import type { AlertSeverity } from '../types';

interface SeverityBadgeProps {
  severity: AlertSeverity;
  size?: 'sm' | 'md';
}

const severityConfig: Record<AlertSeverity, { label: string; text: string; bg: string; dot: string }> = {
  critical: {
    label: 'CRITICAL',
    text: 'text-danger',
    bg: 'bg-danger-muted',
    dot: 'bg-danger',
  },
  high: {
    label: 'HIGH',
    text: 'text-[#EA580C]',
    bg: 'bg-[rgba(234,88,12,0.12)]',
    dot: 'bg-[#EA580C]',
  },
  medium: {
    label: 'MEDIUM',
    text: 'text-warning',
    bg: 'bg-warning-muted',
    dot: 'bg-warning',
  },
  low: {
    label: 'LOW',
    text: 'text-success',
    bg: 'bg-success-muted',
    dot: 'bg-success',
  },
};

export function SeverityBadge({ severity, size = 'sm' }: SeverityBadgeProps) {
  const config = severityConfig[severity];

  if (size === 'md') {
    return (
      <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-xs font-semibold font-mono ${config.bg} ${config.text}`}>
        <span className={`w-1.5 h-1.5 rounded-full ${config.dot}`} />
        {config.label}
      </span>
    );
  }

  return (
    <span className={`inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-bold font-mono ${config.bg} ${config.text}`}>
      <span className={`w-1 h-1 rounded-full ${config.dot}`} />
      {config.label}
    </span>
  );
}
