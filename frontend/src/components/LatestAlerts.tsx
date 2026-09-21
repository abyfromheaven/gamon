import type { Alert } from '../types';

interface LatestAlertsProps {
  alerts: Alert[];
  onViewAll?: () => void;
}

export function LatestAlerts({ alerts, onViewAll }: LatestAlertsProps) {
  return (
    <div className="animate-fade-in-up anim-delay-6 bg-surface border border-border rounded-xl p-5 lg:p-6 h-full flex flex-col">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-sm font-semibold text-text-primary">Latest Alerts</h3>
        <span className="text-xs text-text-muted font-mono">{alerts.length} recent</span>
      </div>

      <div className="space-y-2 flex-1 overflow-y-auto">
        {alerts.map((alert, i) => (
          <div
            key={alert.id}
            className={`animate-slide-in-right anim-delay-${Math.min(i + 5, 7)} flex items-start gap-3 p-3 rounded-lg bg-bg/50 hover:bg-surface-elevated transition-colors cursor-pointer group`}
          >
            <span className="mt-1 w-2 h-2 rounded-full shrink-0 bg-danger" />
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-0.5">
                <span className="font-mono text-xs font-semibold text-text-primary truncate">
                  {alert.device}
                </span>
              </div>
              <p className="text-xs text-text-muted truncate">{alert.title}</p>
            </div>
            <span className="text-[10px] text-text-muted font-mono shrink-0 mt-0.5">
              {alert.startTime}
            </span>
          </div>
        ))}
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
