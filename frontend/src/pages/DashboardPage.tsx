import { useWebSocket } from '../hooks/useWebSocket';
import { dashboardData } from '../data/dummy';
import { TopBar } from '../components/TopBar';
import { MetricsGrid } from '../components/MetricsGrid';
import { DeviceSummary } from '../components/DeviceSummary';
import { LatestAlerts } from '../components/LatestAlerts';
import { SystemStatus } from '../components/SystemStatus';
import { QuickActions } from '../components/QuickActions';

export function DashboardPage() {
  const { isConnected } = useWebSocket();
  const data = dashboardData;

  return (
    <div className="min-h-screen bg-bg">
      <TopBar isConnected={isConnected} />

      <main className="max-w-[1400px] mx-auto px-4 lg:px-8 py-6 lg:py-8 space-y-6 lg:space-y-8">
        {/* Metrics — 4 questions the dashboard answers */}
        <MetricsGrid summary={data.summary} />

        {/* Middle row: Device Summary + Latest Alerts */}
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 lg:gap-5">
          <div className="lg:col-span-2">
            <DeviceSummary devices={data.deviceBreakdown} />
          </div>
          <div className="lg:col-span-3">
            <LatestAlerts alerts={data.latestAlerts} onViewAll={() => {}} />
          </div>
        </div>

        {/* Bottom row: System Status + Quick Actions */}
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-4 lg:gap-5">
          <div className="lg:col-span-3">
            <SystemStatus status={data.systemStatus} />
          </div>
          <div className="lg:col-span-2">
            <QuickActions
              onAddDevice={() => {}}
              onViewMonitoring={() => {}}
              onViewAlerts={() => {}}
              onRefresh={() => {}}
            />
          </div>
        </div>
      </main>
    </div>
  );
}
