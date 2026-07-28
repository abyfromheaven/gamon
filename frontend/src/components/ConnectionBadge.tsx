interface ConnectionBadgeProps {
  isConnected: boolean;
}

export function ConnectionBadge({ isConnected }: ConnectionBadgeProps) {
  return (
    <div
      className={`inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-medium transition-colors duration-300 ${
        isConnected
          ? 'bg-success-muted text-success'
          : 'bg-danger-muted text-danger'
      }`}
    >
      <span
        className={`w-1.5 h-1.5 rounded-full ${
          isConnected ? 'bg-success animate-breathe' : 'bg-danger'
        }`}
      />
      <span className="font-mono">
        {isConnected ? 'Golang Connected' : 'Golang Disconnected'}
      </span>
    </div>
  );
}
