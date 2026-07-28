import type { Alert } from '../types';
import { SeverityBadge } from './SeverityBadge';

interface AlertListProps {
  alerts: Alert[];
  selectedAlert: Alert | null;
  onSelectAlert: (alert: Alert) => void;
}

export function AlertList({ alerts, selectedAlert, onSelectAlert }: AlertListProps) {
  if (alerts.length === 0) {
    return (
      <div className="animate-fade-in bg-surface border border-border rounded-xl p-12 text-center">
        <svg className="w-12 h-12 text-text-muted mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M14.857 17.082a23.848 23.848 0 005.454-1.31A8.967 8.967 0 0118 9.75v-.7V9A6 6 0 006 9v.75a8.967 8.967 0 01-2.312 6.022c1.733.64 3.56 1.085 5.455 1.31m5.714 0a24.255 24.255 0 01-5.714 0m5.714 0a3 3 0 11-5.714 0" />
        </svg>
        <p className="text-sm text-text-muted">No alerts match your filters</p>
      </div>
    );
  }

  return (
    <div className="animate-fade-in-up anim-delay-5 bg-surface border border-border rounded-xl overflow-hidden">
      {/* Desktop Table */}
      <div className="hidden lg:block overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/50">
              <th className="text-left text-[11px] font-semibold text-text-muted uppercase tracking-wider px-5 py-3">Alert</th>
              <th className="text-left text-[11px] font-semibold text-text-muted uppercase tracking-wider px-5 py-3">Device</th>
              <th className="text-left text-[11px] font-semibold text-text-muted uppercase tracking-wider px-5 py-3">Status</th>
              <th className="text-left text-[11px] font-semibold text-text-muted uppercase tracking-wider px-5 py-3">Severity</th>
              <th className="text-left text-[11px] font-semibold text-text-muted uppercase tracking-wider px-5 py-3">Started</th>
            </tr>
          </thead>
          <tbody>
            {alerts.map((alert) => {
              const isSelected = selectedAlert?.id === alert.id;
              return (
                <tr
                  key={alert.id}
                  onClick={() => onSelectAlert(alert)}
                  className={`border-b border-border/30 last:border-0 cursor-pointer transition-colors ${
                    isSelected
                      ? 'bg-accent-muted'
                      : 'hover:bg-surface-elevated/50'
                  }`}
                >
                  <td className="px-5 py-3.5">
                    <span className="text-sm text-text-primary font-medium">{alert.title}</span>
                  </td>
                  <td className="px-5 py-3.5">
                    <div className="flex flex-col">
                      <span className="text-sm text-text-primary font-mono">{alert.device}</span>
                      <span className="text-[11px] text-text-muted">{alert.deviceType}</span>
                    </div>
                  </td>
                  <td className="px-5 py-3.5">
                    <span className={`inline-flex items-center gap-1.5 text-xs font-medium ${
                      alert.status === 'ongoing' ? 'text-danger' : 'text-success'
                    }`}>
                      <span className={`w-1.5 h-1.5 rounded-full ${
                        alert.status === 'ongoing' ? 'bg-danger animate-breathe' : 'bg-success'
                      }`} />
                      {alert.status === 'ongoing' ? 'Ongoing' : 'Resolved'}
                    </span>
                  </td>
                  <td className="px-5 py-3.5">
                    <SeverityBadge severity={alert.severity} />
                  </td>
                  <td className="px-5 py-3.5">
                    <span className="text-xs text-text-muted font-mono">{alert.startTime}</span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Mobile Cards */}
      <div className="lg:hidden divide-y divide-border/30">
        {alerts.map((alert) => {
          const isSelected = selectedAlert?.id === alert.id;
          return (
            <div
              key={alert.id}
              onClick={() => onSelectAlert(alert)}
              className={`p-4 cursor-pointer transition-colors ${
                isSelected ? 'bg-accent-muted' : 'active:bg-surface-elevated/50'
              }`}
            >
              <div className="flex items-start justify-between gap-2 mb-2">
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-text-primary font-medium truncate">{alert.title}</p>
                  <p className="text-xs text-text-muted font-mono mt-0.5">{alert.device}</p>
                </div>
                <SeverityBadge severity={alert.severity} />
              </div>
              <div className="flex items-center justify-between">
                <span className={`inline-flex items-center gap-1.5 text-[11px] font-medium ${
                  alert.status === 'ongoing' ? 'text-danger' : 'text-success'
                }`}>
                  <span className={`w-1 h-1 rounded-full ${
                    alert.status === 'ongoing' ? 'bg-danger' : 'bg-success'
                  }`} />
                  {alert.status === 'ongoing' ? 'Ongoing' : 'Resolved'}
                </span>
                <span className="text-[10px] text-text-muted font-mono">{alert.startTime}</span>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
