import type { PingResult } from '../types';

interface Props {
  result: PingResult;
}

export function StatusCard({ result }: Props) {
  const isOnline = result.status === 'online';

  return (
    <div className={`p-6 rounded-xl border-2 transition-all duration-300 ${
      isOnline
        ? 'bg-green-900/30 border-green-500/50'
        : 'bg-red-900/30 border-red-500/50'
    }`}>
      <div className="flex items-center justify-between mb-4">
        <span className="text-gray-400 text-sm font-mono">{result.ip}</span>
        <span className={`px-3 py-1 rounded-full text-xs font-bold ${
          isOnline
            ? 'bg-green-500/20 text-green-400'
            : 'bg-red-500/20 text-red-400'
        }`}>
          {isOnline ? 'ONLINE' : 'DOWN'}
        </span>
      </div>

      <div className="grid grid-cols-3 gap-4 text-center">
        <div>
          <p className="text-gray-400 text-xs mb-1">Latency</p>
          <p className={`text-xl font-bold ${
            isOnline ? 'text-green-400' : 'text-red-400'
          }`}>
            {isOnline ? `${result.latency.toFixed(1)}` : '—'}
          </p>
          <p className="text-gray-500 text-xs">{isOnline ? 'ms' : ''}</p>
        </div>
        <div>
          <p className="text-gray-400 text-xs mb-1">TTL</p>
          <p className="text-xl font-bold text-blue-400">
            {isOnline ? result.ttl : '—'}
          </p>
        </div>
        <div>
          <p className="text-gray-400 text-xs mb-1">Seq</p>
          <p className="text-xl font-bold text-purple-400">{result.seq}</p>
        </div>
      </div>

      <div className="mt-4 pt-3 border-t border-gray-700/50">
        <p className="text-gray-500 text-xs text-right">
          {new Date(result.timestamp).toLocaleTimeString('id-ID')}
        </p>
      </div>
    </div>
  );
}
