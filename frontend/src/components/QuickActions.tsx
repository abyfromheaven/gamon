interface QuickActionsProps {
  onAddDevice?: () => void;
  onViewMonitoring?: () => void;
  onViewAlerts?: () => void;
  onRefresh?: () => void;
}

interface ActionButtonProps {
  label: string;
  icon: string;
  onClick?: () => void;
}

function ActionButton({ label, icon, onClick }: ActionButtonProps) {
  return (
    <button
      onClick={onClick}
      className="flex items-center gap-3 w-full px-4 py-3 rounded-lg bg-bg/50 hover:bg-surface-elevated border border-border/50 hover:border-border text-sm text-text-secondary hover:text-text-primary transition-all duration-200 cursor-pointer group"
    >
      <span className="text-text-muted group-hover:text-accent transition-colors font-mono text-xs">
        {icon}
      </span>
      <span className="font-medium">{label}</span>
    </button>
  );
}

export function QuickActions({ onAddDevice, onViewMonitoring, onViewAlerts, onRefresh }: QuickActionsProps) {
  return (
    <div className="animate-fade-in-up anim-delay-8 bg-surface border border-border rounded-xl p-5 lg:p-6 h-full">
      <h3 className="text-sm font-semibold text-text-primary mb-4">Quick Actions</h3>

      <div className="grid grid-cols-2 gap-2">
        <ActionButton label="Add Device" icon="+" onClick={onAddDevice} />
        <ActionButton label="Monitoring" icon="◎" onClick={onViewMonitoring} />
        <ActionButton label="Alerts" icon="!" onClick={onViewAlerts} />
        <ActionButton label="Refresh" icon="↻" onClick={onRefresh} />
      </div>
    </div>
  );
}
