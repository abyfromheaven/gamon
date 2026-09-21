import type { Alert } from '../types';

interface AlertSummaryCardsProps {
  alerts: Alert[];
}

export function AlertSummaryCards({ alerts }: AlertSummaryCardsProps) {
  const total = alerts.length;
  const ongoing = alerts.filter((a) => a.status === 'ongoing').length;
  const resolved = alerts.filter((a) => a.status === 'resolved').length;

  return (
    <section className="grid grid-cols-3 gap-3 lg:gap-4">
      <div className="animate-fade-in-up anim-delay-1 bg-surface border border-border rounded-lg px-3 py-2 transition-colors duration-200 hover:bg-surface-elevated flex items-center justify-between">
        <div>
          <p className="text-[10px] uppercase tracking-wider text-text-secondary font-medium">
            Total Alerts
          </p>
        </div>
        <p className="font-mono text-xl lg:text-2xl font-bold tracking-tight text-text-primary">
          {total}
        </p>
      </div>
      <div className="animate-fade-in-up anim-delay-2 bg-surface border border-danger/20 rounded-lg px-3 py-2 transition-colors duration-200 hover:bg-surface-elevated flex items-center justify-between">
        <div>
          <p className="text-[10px] uppercase tracking-wider text-text-secondary font-medium">
            Ongoing
          </p>
        </div>
        <p className="font-mono text-xl lg:text-2xl font-bold tracking-tight text-danger">
          {ongoing}
        </p>
      </div>
      <div className="animate-fade-in-up anim-delay-3 bg-surface border border-success/20 rounded-lg px-3 py-2 transition-colors duration-200 hover:bg-surface-elevated flex items-center justify-between">
        <div>
          <p className="text-[10px] uppercase tracking-wider text-text-secondary font-medium">
            Resolved
          </p>
        </div>
        <p className="font-mono text-xl lg:text-2xl font-bold tracking-tight text-success">
          {resolved}
        </p>
      </div>
    </section>
  );
}
