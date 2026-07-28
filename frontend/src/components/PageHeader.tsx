interface PageHeaderProps {
  title: string;
  subtitle: string;
  actionLabel: string;
  onAction: () => void;
}

export function PageHeader({ title, subtitle, actionLabel, onAction }: PageHeaderProps) {
  return (
    <div className="animate-fade-in flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
      <div>
        <h1 className="text-2xl font-bold text-text-primary tracking-tight">{title}</h1>
        <p className="mt-1 text-sm text-text-muted">{subtitle}</p>
      </div>
      <button
        onClick={onAction}
        className="inline-flex items-center gap-2 px-4 py-2.5 bg-accent hover:bg-accent/90 text-white text-sm font-semibold rounded-lg transition-colors duration-200 cursor-pointer shrink-0"
      >
        <span className="text-lg leading-none">+</span>
        {actionLabel}
      </button>
    </div>
  );
}
