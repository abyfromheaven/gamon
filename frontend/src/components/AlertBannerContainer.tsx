import { useCallback, useEffect, useState } from 'react';
import { AlertBanner, type AlertBannerData } from './AlertBanner';
import type { StatusChange } from '../types';

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

    const banner: AlertBannerData = {
      device_id: statusChange.device_id,
      device_name: statusChange.device_name,
      device_ip: '',
      old_status: statusChange.old_status,
      new_status: statusChange.new_status,
      timestamp: statusChange.timestamp,
    };

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

  return (
    <AlertBanner
      data={current}
      onDetail={onNavigateToMonitoring}
      onDismiss={dismiss}
    />
  );
}
