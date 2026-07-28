import type { SystemStatusInfo } from '../types';

interface SystemStatusProps {
  status: SystemStatusInfo;
}

interface StatusRowProps {
  label: string;
  value: string;
  isRunning: boolean;
}

function StatusRow({ label, value, isRunning }: StatusRowProps) {
  return (
    <div className="flex items-center justify-between py-2">
      <span className="text-xs text-text-muted">{label}</span>
      <div className="flex items-center gap-2">
        <span className="font-mono text-xs text-text-primary font-medium">{value}</span>
        <span
          className={`w-1.5 h-1.5 rounded-full ${
            isRunning ? 'bg-success animate-breathe' : 'bg-danger'
          }`}
        />
      </div>
    </div>
  );
}

export function SystemStatus({ status }: SystemStatusProps) {
  return (
    <div className="animate-fade-in-up anim-delay-7 bg-surface border border-border rounded-xl p-5 lg:p-6 h-full">
      <h3 className="text-sm font-semibold text-text-primary mb-4">System Status</h3>

      <div className="divide-y divide-border/50">
        <StatusRow
          label="Monitoring"
          value={status.monitoring}
          isRunning={status.monitoring === 'Running'}
        />
        <StatusRow
          label="Check Interval"
          value={status.checkInterval}
          isRunning={status.monitoring === 'Running'}
        />
        <StatusRow
          label="Last Scan"
          value={status.lastScan}
          isRunning={status.monitoring === 'Running'}
        />
        <StatusRow
          label="Notifications"
          value={status.notifications}
          isRunning={status.notifications === 'Active'}
        />
      </div>
    </div>
  );
}
