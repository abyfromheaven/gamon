import { useCallback, useState } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { Sidebar, type Page } from './components/Sidebar';
import { TopBar } from './components/TopBar';
import { DashboardPage } from './pages/DashboardPage';
import { DeviceManagementPage } from './pages/DeviceManagementPage';
import { AlertCenterPage } from './pages/AlertCenterPage';
import { MonitoringPage } from './pages/MonitoringPage';

function App() {
  const [reconnectKey, setReconnectKey] = useState(0);
  const handleReconnect = useCallback(() => setReconnectKey((k) => k + 1), []);
  const { isConnected, monitorResults, lastStatusChange } = useWebSocket(handleReconnect);
  const [currentPage, setCurrentPage] = useState<Page>('dashboard');

  return (
    <div className="min-h-screen bg-bg flex">
      <Sidebar activePage={currentPage} onNavigate={setCurrentPage} />

      <div className="flex-1 flex flex-col min-w-0">
        <div className="lg:hidden">
          <TopBar isConnected={isConnected} />
        </div>
        <div className="hidden lg:block">
          <TopBar isConnected={isConnected} />
        </div>

        <main className="flex-1 overflow-auto">
          {currentPage === 'dashboard' && <DashboardPage monitorResults={monitorResults} onNavigate={setCurrentPage} isConnected={isConnected} reconnectKey={reconnectKey} />}
          {currentPage === 'devices' && <DeviceManagementPage />}
          {currentPage === 'monitoring' && <MonitoringPage monitorResults={monitorResults} onViewAlerts={() => setCurrentPage('alerts')} reconnectKey={reconnectKey} />}
          {currentPage === 'alerts' && <AlertCenterPage lastStatusChange={lastStatusChange} />}
        </main>
      </div>
    </div>
  );
}

export default App;
