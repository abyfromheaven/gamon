import { useEffect, useState } from 'react';

export interface AlertBannerData {
  device_id: number;
  device_name: string;
  device_ip: string;
  old_status: string;
  new_status: string;
  timestamp: string;
}

interface AlertBannerProps {
  data: AlertBannerData;
  onDetail: (deviceId: number) => void;
  onDismiss: () => void;
}

export function AlertBanner({ data, onDetail, onDismiss }: AlertBannerProps) {
  const [visible, setVisible] = useState(false);
  const [flashing, setFlashing] = useState(true);

  const isOffline = data.new_status === 'offline';
  const bgColor = isOffline ? 'bg-danger' : 'bg-warning';
  const textColor = isOffline ? 'text-white' : 'text-bg';

  useEffect(() => {
    requestAnimationFrame(() => setVisible(true));
    const flashTimer = setTimeout(() => setFlashing(false), 1200);
    const dismissTimer = setTimeout(() => {
      setVisible(false);
      setTimeout(onDismiss, 400);
    }, 10000);
    return () => {
      clearTimeout(flashTimer);
      clearTimeout(dismissTimer);
    };
  }, [onDismiss]);

  return (
    <div
      className={`fixed top-14 left-0 right-0 z-40 ${bgColor} ${textColor} shadow-2xl transition-all duration-400 ${
        visible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-full'
      } ${flashing ? 'animate-banner-flash' : ''}`}
    >
      <div className="max-w-[1400px] mx-auto px-4 lg:px-8 py-3 flex items-center gap-4">
        <div className="relative shrink-0">
          <svg className="w-6 h-6 animate-pulse-ring" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
          </svg>
        </div>
        <div className="flex-1 min-w-0">
          <p className="font-semibold text-sm">
            PERINGATAN: {data.device_name} ({data.device_ip}) berubah dari{' '}
            <span className="uppercase font-bold">{data.old_status}</span> ke{' '}
            <span className="uppercase font-bold">{data.new_status}</span>
          </p>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <button
            onClick={() => onDetail(data.device_id)}
            className="px-3 py-1.5 rounded-md bg-white/20 hover:bg-white/30 text-xs font-semibold transition-colors cursor-pointer"
          >
            Lihat Detail
          </button>
          <button
            onClick={() => {
              setVisible(false);
              setTimeout(onDismiss, 400);
            }}
            className="p-1 rounded hover:bg-white/20 transition-colors cursor-pointer"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}
