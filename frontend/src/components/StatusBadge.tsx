interface StatusBadgeProps {
  status: 'active' | 'inactive' | 'warning';
}

const statusConfig = {
  active: {
    label: 'Aktif',
    dot: 'bg-success animate-breathe',
    text: 'text-success',
    bg: 'bg-success-muted',
  },
  inactive: {
    label: 'Nonaktif',
    dot: 'bg-text-muted',
    text: 'text-text-muted',
    bg: 'bg-bg/50',
  },
  warning: {
    label: 'Warning',
    dot: 'bg-warning',
    text: 'text-warning',
    bg: 'bg-warning-muted',
  },
};

export function StatusBadge({ status }: StatusBadgeProps) {
  const config = statusConfig[status];

  return (
    <span className={`inline-flex items-center gap-1.5 px-2 py-1 rounded-full text-xs font-medium ${config.bg} ${config.text}`}>
      <span className={`w-1.5 h-1.5 rounded-full ${config.dot}`} />
      {config.label}
    </span>
  );
}
