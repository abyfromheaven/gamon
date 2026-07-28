interface MetricCardProps {
  value: number;
  label: string;
  sublabel?: string;
  accent?: 'default' | 'success' | 'danger' | 'warning';
  delay?: number;
}

const accentColors = {
  default: 'text-text-primary',
  success: 'text-success',
  danger: 'text-danger',
  warning: 'text-warning',
};

const borderColors = {
  default: 'border-border',
  success: 'border-success/20',
  danger: 'border-danger/20',
  warning: 'border-warning/20',
};

export function MetricCard({ value, label, sublabel, accent = 'default', delay = 0 }: MetricCardProps) {
  return (
    <div
      className={`animate-fade-in-up anim-delay-${delay} bg-surface border ${borderColors[accent]} rounded-xl p-5 lg:p-6 transition-colors duration-200 hover:bg-surface-elevated`}
    >
      <p className={`font-mono text-4xl lg:text-5xl font-bold tracking-tight ${accentColors[accent]}`}>
        {value}
      </p>
      <p className="mt-2 text-xs uppercase tracking-widest text-text-secondary font-medium">
        {label}
      </p>
      {sublabel && (
        <p className="mt-1 text-xs text-text-muted">{sublabel}</p>
      )}
    </div>
  );
}
