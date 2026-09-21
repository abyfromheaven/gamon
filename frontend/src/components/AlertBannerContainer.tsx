import { useCallback, useEffect, useRef, useState } from 'react';
import { AlertBanner, type AlertBannerData } from './AlertBanner';
import { EmergencyAlertModal } from './EmergencyAlertModal';
import type { MonitorResult, StatusChange } from '../types';
import { audioAlert } from '../lib/audioAlert';
import { sendDesktopNotification, startTabTitleFlash, stopTabTitleFlash } from '../lib/desktopNotify';
import { getNotificationClientSettings } from '../lib/notificationSettings';

interface AlertBannerContainerProps {
  statusChange: StatusChange | null;
  monitorResults?: Map<number, MonitorResult>;
  onNavigateToMonitoring: (deviceId: number) => void;
}

interface AcknowledgedTrack {
  deviceName: string;
  deviceIp: string;
  acknowledgedAt: number;
  lastAlertedAt: number;
}

export function AlertBannerContainer({ statusChange, monitorResults, onNavigateToMonitoring }: AlertBannerContainerProps) {
  const [queue, setQueue] = useState<AlertBannerData[]>([]);
  const [current, setCurrent] = useState<AlertBannerData | null>(null);
  const acknowledgedMap = useRef<Map<number, AcknowledgedTrack>>(new Map());
  const maxQueue = 3;

  // Handle incoming status change events (WebSocket)
  useEffect(() => {
    if (!statusChange) return;
    if (!statusChange.old_status) return; // skip new device insertion

    const isOffline = statusChange.new_status === 'offline';
    const isRecovery = statusChange.new_status === 'online' && statusChange.old_status === 'offline';

    // If device recovered, remove it from acknowledged offline tracker
    if (isRecovery) {
      acknowledgedMap.current.delete(statusChange.device_id);
    }

    const banner: AlertBannerData = {
      device_id: statusChange.device_id,
      device_name: statusChange.device_name,
      device_ip: '',
      old_status: statusChange.old_status,
      new_status: statusChange.new_status,
      timestamp: statusChange.timestamp,
    };

    // 1. Mandatory Alert: Native Desktop Push Notification & Tab Title Flash
    if (isOffline) {
      sendDesktopNotification(
        `🚨 PERINGATAN CRITICAL: ${statusChange.device_name} OFFLINE!`,
        `Perangkat ${statusChange.device_name} berubah status dari ${statusChange.old_status} menjadi ${statusChange.new_status}.`
      );
      startTabTitleFlash(`${statusChange.device_name} OFFLINE`);
    }

    // 2. Read user settings for Sound Alarm
    const settings = getNotificationClientSettings();

    if (isOffline && settings.soundAlarmEnabled) {
      audioAlert.setVolume(settings.volume);
      audioAlert.startSiren();
    }

    setCurrent((prev) => {
      if (!prev) return banner;
      setQueue((q) => {
        const next = [...q, banner];
        return next.length > maxQueue ? next.slice(-maxQueue) : next;
      });
      return prev;
    });
  }, [statusChange]);

  // Periodic Re-notification Checker (Runs every 10 seconds)
  useEffect(() => {
    const interval = setInterval(() => {
      const settings = getNotificationClientSettings();
      const reNotifMinutes = settings.reNotificationIntervalMinutes;
      if (!reNotifMinutes || reNotifMinutes <= 0) return;

      const intervalMs = reNotifMinutes * 60 * 1000;
      const now = Date.now();

      acknowledgedMap.current.forEach((track, deviceId) => {
        // Check live monitor status if available
        if (monitorResults) {
          const liveResult = monitorResults.get(deviceId);
          if (liveResult && liveResult.status === 'online') {
            // Device has recovered!
            acknowledgedMap.current.delete(deviceId);
            return;
          }
        }

        // Check if re-notification interval has elapsed
        if (now - track.lastAlertedAt >= intervalMs) {
          track.lastAlertedAt = now;

          const reAlertBanner: AlertBannerData = {
            device_id: deviceId,
            device_name: track.deviceName,
            device_ip: track.deviceIp,
            old_status: 'offline',
            new_status: 'offline',
            timestamp: new Date().toISOString(),
            is_re_notification: true,
          };

          // Trigger sound, title flash, desktop notification
          sendDesktopNotification(
            `⚠️ PENGINGAT ULANG: ${track.deviceName} MASIH OFFLINE!`,
            `Perangkat ${track.deviceName} masih belum pulih setelah di-acknowledge.`
          );
          startTabTitleFlash(`${track.deviceName} MASIH OFFLINE`);

          if (settings.soundAlarmEnabled) {
            audioAlert.setVolume(settings.volume);
            audioAlert.startSiren();
          }

          // Enqueue or display
          setCurrent((prev) => {
            if (!prev) return reAlertBanner;
            setQueue((q) => [...q, reAlertBanner]);
            return prev;
          });
        }
      });
    }, 10000);

    return () => clearInterval(interval);
  }, [monitorResults]);

  const dismiss = useCallback(() => {
    // When current alert is acknowledged or dismissed, track it if offline
    if (current && current.new_status === 'offline') {
      const now = Date.now();
      acknowledgedMap.current.set(current.device_id, {
        deviceName: current.device_name,
        deviceIp: current.device_ip,
        acknowledgedAt: now,
        lastAlertedAt: now,
      });
    }

    // Stop siren sound & tab title flashing when alert is acknowledged/dismissed
    audioAlert.stopSiren();
    stopTabTitleFlash();

    setCurrent(null);
    setQueue((prev) => {
      const next = [...prev];
      next.shift();
      return next;
    });
  }, [current]);

  useEffect(() => {
    if (!current && queue.length > 0) {
      setCurrent(queue[0]);
      setQueue((prev) => prev.slice(1));
    }
  }, [current, queue]);

  if (!current) return null;

  const clientSettings = getNotificationClientSettings();
  const isOffline = current.new_status === 'offline';
  const showBlockScreenModal = isOffline && clientSettings.blockScreenEnabled;

  return (
    <>
      {showBlockScreenModal ? (
        <EmergencyAlertModal
          data={current}
          onAcknowledge={() => {
            onNavigateToMonitoring(current.device_id);
            dismiss();
          }}
        />
      ) : (
        <AlertBanner
          data={current}
          onDetail={onNavigateToMonitoring}
          onDismiss={dismiss}
        />
      )}
    </>
  );
}
