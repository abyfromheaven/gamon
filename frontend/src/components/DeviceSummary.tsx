import { StatusIndicator } from './StatusIndicator';
import type { DeviceBreakdown } from '../types';

interface DeviceSummaryProps {
  devices: DeviceBreakdown[];
}

function getDeviceIcon(type: string): string {
  switch (type) {
    case 'Server': return '◻';
    case 'Router': return '◎';
    case 'Switch': return '⬡';
    case 'Access Point': return '◈';
    case 'Website': return '◇';
    default: return '○';
  }
}

export function DeviceSummary({ devices }: DeviceSummaryProps) {
  const total = devices.reduce((sum, d) => sum + d.count, 0);

  return (
    <div className="animate-fade-in-up anim-delay-5 bg-surface border border-border rounded-xl p-5 lg:p-6 h-full">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-semibold text-text-primary">Device Summary</h3>
        <span className="text-xs text-text-muted font-mono">{total} total</span>
      </div>

      <div className="space-y-3">
        {devices.map((device) => {
          const offlineCount = device.count - device.online;
          const allOnline = offlineCount === 0;

          return (
            <div key={device.type} className="flex items-center justify-between group">
              <div className="flex items-center gap-3">
                <span className="text-text-muted text-sm w-4 text-center font-mono">
                  {getDeviceIcon(device.type)}
                </span>
                <span className="text-sm text-text-secondary group-hover:text-text-primary transition-colors">
                  {device.type}
                </span>
              </div>

              <div className="flex items-center gap-3">
                <span className="font-mono text-sm text-text-primary font-medium">
                  {device.online}
                </span>
                <span className="text-text-muted text-xs">/</span>
                <span className="font-mono text-sm text-text-muted">
                  {device.count}
                </span>
                <StatusIndicator status={allOnline ? 'online' : 'warning'} size="sm" />
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
