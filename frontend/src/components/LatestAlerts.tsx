import type { Alert, AlertSeverity } from '../types';

interface LatestAlertsProps {
  alerts: Alert[];
  onViewAll?: () => void;
}

const severityStyles: Record<AlertSeverity, { dot: string; text: string; bg: string }> = {
  critical: {
    dot: 'bg-danger',
    text: 'text-danger',
    bg: 'bg-danger-muted',
  },
  high: {
    dot: 'bg-[#EA580C]',
    text: 'text-[#EA580C]',
    bg: 'bg-[rgba(234,88,12,0.12)]',
  },
  medium: {
    dot: 'bg-warning',
    text: 'text-warning',
    bg: 'bg-warning-muted',
  },
  low: {
    dot: 'bg-success',
    text: 'text-success',
    bg: 'bg-success-muted',
  },
  info: {
    dot: 'bg-accent',
    text: 'text-accent',
    bg: 'bg-accent-muted',
  },
};

const severityLabels: Record<AlertSeverity, string> = {
  critical: 'CRIT',
  high: 'HIGH',
  medium: 'MED',
  low: 'LOW',
  info: 'INFO',
};

export function LatestAlerts({ alerts, onViewAll }: LatestAlertsProps) {
  return (
    <div className="animate-fade-in-up anim-delay-6 bg-surface border border-border rounded-xl p-5 lg:p-6 h-full flex flex-col">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-semibold text-text-primary">Latest Alerts</h3>
        <span className="text-xs text-text-muted font-mono">{alerts.length} recent</span>
      </div>

      <div className="space-y-2 flex-1 overflow-y-auto">
        {alerts.map((alert, i) => {
          const style = severityStyles[alert.severity];
          return (
            <div
              key={alert.id}
              className={`animate-slide-in-right anim-delay-${Math.min(i + 5, 7)} flex items-start gap-3 p-3 rounded-lg bg-bg/50 hover:bg-surface-elevated transition-colors cursor-pointer group`}
            >
              <span className={`mt-1 w-2 h-2 rounded-full shrink-0 ${style.dot}`} />
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-0.5">
                  <span className="font-mono text-xs font-semibold text-text-primary truncate">
                    {alert.device}
                  </span>
                  <span className={`text-[10px] font-mono font-bold px-1.5 py-0.5 rounded ${style.bg} ${style.text}`}>
                    {severityLabels[alert.severity]}
                  </span>
                </div>
                <p className="text-xs text-text-muted truncate">{alert.title}</p>
              </div>
              <span className="text-[10px] text-text-muted font-mono shrink-0 mt-0.5">
                {alert.startTime}
              </span>
            </div>
          );
        })}
      </div>

      {onViewAll && (
        <button
          onClick={onViewAll}
          className="mt-4 text-xs text-accent hover:text-accent/80 font-medium transition-colors text-center"
        >
          View all alerts →
        </button>
      )}
    </div>
  );
}
