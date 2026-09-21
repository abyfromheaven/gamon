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
      className={`animate-fade-in-up anim-delay-${delay} bg-surface border ${borderColors[accent]} rounded-lg px-3 py-2 transition-colors duration-200 hover:bg-surface-elevated flex items-center justify-between`}
    >
      <div>
        <p className="text-[10px] uppercase tracking-wider text-text-secondary font-medium">
          {label}
        </p>
        {sublabel && (
          <p className="text-[9px] text-text-muted">{sublabel}</p>
        )}
      </div>
      <p className={`font-mono text-xl lg:text-2xl font-bold tracking-tight ${accentColors[accent]}`}>
        {value}
      </p>
    </div>
  );
}
