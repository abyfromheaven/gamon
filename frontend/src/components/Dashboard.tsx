import type { PingResult } from '../types';
import { StatusCard } from './StatusCard';

interface Props {
  result: PingResult | null;
  isConnected: boolean;
}

export function Dashboard({ result, isConnected }: Props) {
  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold text-white">Status Monitoring</h2>
        <span className={`text-xs px-2 py-1 rounded ${
          isConnected
            ? 'bg-green-900/50 text-green-400'
            : 'bg-red-900/50 text-red-400'
        }`}>
          {isConnected ? 'WS Connected' : 'WS Disconnected'}
        </span>
      </div>

      {result ? (
        <StatusCard result={result} />
      ) : (
        <div className="text-center py-12 text-gray-500">
          <p className="text-lg">Belum ada target monitoring</p>
          <p className="text-sm mt-1">Masukkan IP di atas untuk memulai</p>
        </div>
      )}
    </div>
  );
}
