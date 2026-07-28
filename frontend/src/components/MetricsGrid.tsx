import { MetricCard } from './MetricCard';
import type { DashboardData } from '../types';

interface MetricsGridProps {
  summary: DashboardData['summary'];
}

export function MetricsGrid({ summary }: MetricsGridProps) {
  return (
    <section className="grid grid-cols-2 lg:grid-cols-4 gap-3 lg:gap-4">
      <MetricCard
        value={summary.totalDevices}
        label="Total Device"
        sublabel="sedang dimonitor"
        accent="default"
        delay={1}
      />
      <MetricCard
        value={summary.online}
        label="Online"
        sublabel="perangkat aktif"
        accent="success"
        delay={2}
      />
      <MetricCard
        value={summary.offline}
        label="Offline"
        sublabel="tidak merespon"
        accent="danger"
        delay={3}
      />
      <MetricCard
        value={summary.warnings}
        label="Warning"
        sublabel="perlu perhatian"
        accent="warning"
        delay={4}
      />
    </section>
  );
}
