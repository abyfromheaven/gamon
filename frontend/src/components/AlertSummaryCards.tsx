import type { Alert } from '../types';

interface AlertSummaryCardsProps {
  alerts: Alert[];
}

export function AlertSummaryCards({ alerts }: AlertSummaryCardsProps) {
  const total = alerts.length;
  const ongoing = alerts.filter((a) => a.status === 'ongoing').length;
  const resolved = alerts.filter((a) => a.status === 'resolved').length;
  const critical = alerts.filter((a) => a.severity === 'critical' && a.status === 'ongoing').length;

  const cards = [
    {
      label: 'Total Alerts',
      value: total,
      color: 'text-text-primary',
      border: 'border-border',
    },
    {
      label: 'Ongoing',
      value: ongoing,
      color: 'text-danger',
      border: 'border-danger/30',
    },
    {
      label: 'Resolved',
      value: resolved,
      color: 'text-success',
      border: 'border-success/30',
    },
    {
      label: 'Critical',
      value: critical,
      color: 'text-[#EA580C]',
      border: 'border-[rgba(234,88,12,0.3)]',
    },
  ];

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
      {cards.map((card, i) => (
        <div
          key={card.label}
          className={`animate-fade-in-up anim-delay-${i} bg-surface border ${card.border} rounded-xl p-4 lg:p-5`}
        >
          <p className="text-xs text-text-muted font-medium mb-1">{card.label}</p>
          <p className={`text-2xl lg:text-3xl font-bold font-mono ${card.color}`}>
            {card.value}
          </p>
        </div>
      ))}
    </div>
  );
}
