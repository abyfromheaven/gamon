import type { Alert } from '../types';
import { SeverityBadge } from './SeverityBadge';

interface AlertDetailPanelProps {
  alert: Alert | null;
  onClose: () => void;
  onMarkResolved: (id: number) => void;
}

export function AlertDetailPanel({ alert, onClose, onMarkResolved }: AlertDetailPanelProps) {
  if (!alert) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/40 z-40 lg:hidden"
        onClick={onClose}
      />

      {/* Panel */}
      <div className="fixed top-0 right-0 h-full w-full max-w-md bg-surface border-l border-border z-50 flex flex-col animate-slide-in-right">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border/50">
          <h3 className="text-sm font-semibold text-text-primary">Alert Detail</h3>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg hover:bg-surface-elevated transition-colors text-text-muted hover:text-text-primary"
          >
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto px-5 py-5 space-y-5">
          {/* Title + Status */}
          <div>
            <div className="flex items-center gap-2 mb-2">
              <SeverityBadge severity={alert.severity} size="md" />
              <span className={`inline-flex items-center gap-1.5 text-xs font-medium ${
                alert.status === 'ongoing' ? 'text-danger' : 'text-success'
              }`}>
                <span className={`w-1.5 h-1.5 rounded-full ${
                  alert.status === 'ongoing' ? 'bg-danger animate-breathe' : 'bg-success'
                }`} />
                {alert.status === 'ongoing' ? 'Ongoing' : 'Resolved'}
              </span>
            </div>
            <h2 className="text-lg font-semibold text-text-primary">{alert.title}</h2>
          </div>

          {/* Info Grid */}
          <div className="space-y-3">
            <InfoRow label="Device Name" value={alert.device} />
            <InfoRow label="Device Type" value={alert.deviceType} />
            <InfoRow label="Monitoring Method" value={alert.monitoringMethod} />
            <InfoRow label="Started" value={alert.startTime} />
            <InfoRow
              label="Resolved"
              value={alert.resolvedTime ?? '—'}
              valueClass={alert.resolvedTime ? 'text-success' : 'text-text-muted'}
            />
          </div>

          {/* Description */}
          <div>
            <p className="text-[11px] text-text-muted uppercase tracking-wider font-semibold mb-2">Description</p>
            <p className="text-sm text-text-secondary leading-relaxed bg-bg/50 rounded-lg p-3 border border-border/30">
              {alert.description}
            </p>
          </div>
        </div>

        {/* Footer Action */}
        {alert.status === 'ongoing' && (
          <div className="px-5 py-4 border-t border-border/50">
            <button
              onClick={() => onMarkResolved(alert.id)}
              className="w-full px-4 py-2.5 bg-success/10 border border-success/30 text-success text-sm font-medium rounded-lg hover:bg-success/20 transition-colors cursor-pointer"
            >
              Mark as Resolved
            </button>
          </div>
        )}
      </div>
    </>
  );
}

function InfoRow({ label, value, valueClass }: { label: string; value: string; valueClass?: string }) {
  return (
    <div className="flex items-center justify-between py-2 border-b border-border/20 last:border-0">
      <span className="text-[11px] text-text-muted uppercase tracking-wider font-semibold">{label}</span>
      <span className={`text-sm font-mono ${valueClass ?? 'text-text-primary'}`}>{value}</span>
    </div>
  );
}
