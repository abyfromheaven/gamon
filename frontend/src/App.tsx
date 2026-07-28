import { useState } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { Sidebar, type Page } from './components/Sidebar';
import { TopBar } from './components/TopBar';
import { DashboardPage } from './pages/DashboardPage';
import { DeviceManagementPage } from './pages/DeviceManagementPage';

function App() {
  const { isConnected } = useWebSocket();
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
          {currentPage === 'dashboard' && <DashboardPage />}
          {currentPage === 'devices' && <DeviceManagementPage />}
        </main>
      </div>
    </div>
  );
}

export default App;
