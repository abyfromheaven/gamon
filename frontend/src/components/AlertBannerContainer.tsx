import { useCallback, useEffect, useState } from 'react';
import { AlertBanner, type AlertBannerData } from './AlertBanner';
import { EmergencyAlertModal } from './EmergencyAlertModal';
import type { StatusChange } from '../types';
import { audioAlert } from '../lib/audioAlert';
import { sendDesktopNotification, startTabTitleFlash, stopTabTitleFlash } from '../lib/desktopNotify';
import { getNotificationClientSettings } from '../lib/notificationSettings';

interface AlertBannerContainerProps {
  statusChange: StatusChange | null;
  onNavigateToMonitoring: (deviceId: number) => void;
}

export function AlertBannerContainer({ statusChange, onNavigateToMonitoring }: AlertBannerContainerProps) {
  const [queue, setQueue] = useState<AlertBannerData[]>([]);
  const [current, setCurrent] = useState<AlertBannerData | null>(null);
  const maxQueue = 3;

  useEffect(() => {
    if (!statusChange) return;
    if (!statusChange.old_status) return; // skip device baru

    const banner: AlertBannerData = {
      device_id: statusChange.device_id,
      device_name: statusChange.device_name,
      device_ip: '',
      old_status: statusChange.old_status,
      new_status: statusChange.new_status,
      timestamp: statusChange.timestamp,
    };

    // 1. Mandatory Alert: Native Desktop Push Notification & Tab Title Flash
    const isOffline = statusChange.new_status === 'offline';
    if (isOffline) {
      sendDesktopNotification(
        `🚨 PERINGATAN CRITICAL: ${statusChange.device_name} OFFLINE!`,
        `Perangkat ${statusChange.device_name} berubah status dari ${statusChange.old_status} menjadi ${statusChange.new_status}.`
      );
      startTabTitleFlash(`${statusChange.device_name} OFFLINE`);
    }

    // 2. Read user settings for Sound Alarm & Block Screen Modal
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

  const dismiss = useCallback(() => {
    // Stop siren sound & tab title flashing when alert is acknowledged/dismissed
    audioAlert.stopSiren();
    stopTabTitleFlash();

    setCurrent(null);
    setQueue((prev) => {
      const next = [...prev];
      next.shift();
      return next;
    });
  }, []);

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
