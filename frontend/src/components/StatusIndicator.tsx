interface StatusIndicatorProps {
  status: 'online' | 'offline' | 'warning' | 'unknown';
  size?: 'sm' | 'md' | 'lg';
}

const sizeMap = {
  sm: 'w-2 h-2',
  md: 'w-2.5 h-2.5',
  lg: 'w-3 h-3',
};

const colorMap = {
  online: 'bg-success',
  offline: 'bg-danger',
  warning: 'bg-warning',
  unknown: 'bg-text-muted',
};

export function StatusIndicator({ status, size = 'md' }: StatusIndicatorProps) {
  return (
    <span className="relative inline-flex">
      {status === 'online' && (
        <span
          className={`absolute inline-flex h-full w-full rounded-full bg-success opacity-40 animate-pulse-ring`}
        />
      )}
      <span
        className={`relative inline-flex rounded-full ${sizeMap[size]} ${colorMap[status]} ${
          status === 'online' ? 'animate-breathe' : ''
        }`}
      />
    </span>
  );
}
