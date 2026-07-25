import { useState } from 'react';
import { useWebSocket } from './hooks/useWebSocket';
import { InputForm } from './components/InputForm';
import { Dashboard } from './components/Dashboard';

function App() {
  const { latestResult, isConnected } = useWebSocket();
  const [monitoringIp, setMonitoringIp] = useState<string | null>(null);

  const handleMonitor = (ip: string) => {
    setMonitoringIp(ip);
  };

  return (
    <div className="min-h-screen bg-gray-900 text-white">
      <div className="max-w-2xl mx-auto px-4 py-8">
        {/* Header */}
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-green-400 mb-1">GAMON</h1>
          <p className="text-gray-500 text-sm">Garda Monitoring — Realtime Network Monitor</p>
        </div>

        {/* Input */}
        <div className="mb-8">
          <InputForm onMonitor={handleMonitor} />
        </div>

        {/* Monitoring Target Label */}
        {monitoringIp && (
          <div className="mb-4 text-center">
            <span className="text-gray-400 text-sm">
              Monitoring target:{' '}
              <span className="text-white font-mono">{monitoringIp}</span>
            </span>
          </div>
        )}

        {/* Dashboard */}
        <Dashboard result={latestResult} isConnected={isConnected} />
      </div>
    </div>
  );
}

export default App;
