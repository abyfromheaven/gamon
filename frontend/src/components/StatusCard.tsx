import type { PingResult } from '../types';

interface StatusCardProps {
  result: PingResult;
}

export function StatusCard({ result }: StatusCardProps) {
  const isOnline = result.status === 'online';

  return (
    <div className="rounded-xl border border-border bg-surface p-5">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-text-primary">{result.ip}</p>
          <p className="mt-1 text-xs text-text-muted">Last check: {result.timestamp}</p>
        </div>
        <span className={isOnline ? 'text-success' : 'text-danger'}>
          {isOnline ? 'Online' : 'Offline'}
        </span>
      </div>

      <dl className="mt-5 grid grid-cols-3 gap-3 text-sm">
        <div>
          <dt className="text-text-muted">Latency</dt>
          <dd className="mt-1 font-mono text-text-primary">{result.latency} ms</dd>
        </div>
        <div>
          <dt className="text-text-muted">TTL</dt>
          <dd className="mt-1 font-mono text-text-primary">{result.ttl}</dd>
        </div>
        <div>
          <dt className="text-text-muted">Sequence</dt>
          <dd className="mt-1 font-mono text-text-primary">#{result.seq}</dd>
        </div>
      </dl>
    </div>
  );
}
