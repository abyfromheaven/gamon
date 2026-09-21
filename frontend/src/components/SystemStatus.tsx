import type { SystemStatusInfo } from '../types';

interface SystemStatusProps {
  status: SystemStatusInfo;
}

export function SystemStatus({ status }: SystemStatusProps) {
  return (
    <div className="animate-fade-in-up anim-delay-7 bg-surface border border-border rounded-xl p-3">
      <h3 className="text-xs font-semibold text-text-primary mb-2">System Status</h3>

      <div className="grid grid-cols-2 gap-x-3 gap-y-1.5">
        <div className="flex items-center justify-between py-1 px-2 rounded bg-bg/40">
          <span className="text-[11px] text-text-muted">Monitoring</span>
          <div className="flex items-center gap-1.5">
            <span className="font-mono text-[11px] font-medium text-text-primary">{status.monitoring}</span>
            <span className={`w-1.5 h-1.5 rounded-full ${status.monitoring === 'Running' ? 'bg-success animate-breathe' : 'bg-danger'}`} />
          </div>
        </div>
        <div className="flex items-center justify-between py-1 px-2 rounded bg-bg/40">
          <span className="text-[11px] text-text-muted">Interval</span>
          <div className="flex items-center gap-1.5">
            <span className="font-mono text-[11px] font-medium text-text-primary">{status.checkInterval}</span>
            <span className={`w-1.5 h-1.5 rounded-full ${status.monitoring === 'Running' ? 'bg-success animate-breathe' : 'bg-danger'}`} />
          </div>
        </div>
        <div className="flex items-center justify-between py-1 px-2 rounded bg-bg/40">
          <span className="text-[11px] text-text-muted">Last Scan</span>
          <div className="flex items-center gap-1.5">
            <span className="font-mono text-[11px] font-medium text-text-primary">{status.lastScan}</span>
            <span className={`w-1.5 h-1.5 rounded-full ${status.monitoring === 'Running' ? 'bg-success animate-breathe' : 'bg-danger'}`} />
          </div>
        </div>
        <div className="flex items-center justify-between py-1 px-2 rounded bg-bg/40">
          <span className="text-[11px] text-text-muted">Notif</span>
          <div className="flex items-center gap-1.5">
            <span className="font-mono text-[11px] font-medium text-text-primary">{status.notifications}</span>
            <span className={`w-1.5 h-1.5 rounded-full ${status.notifications === 'Active' ? 'bg-success animate-breathe' : 'bg-danger'}`} />
          </div>
        </div>
      </div>
    </div>
  );
}
