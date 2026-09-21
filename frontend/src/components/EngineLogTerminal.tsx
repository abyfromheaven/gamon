import { useEffect, useRef, useState } from 'react';

export interface EngineLogEntry {
  id: string;
  timestamp: string;
  deviceName: string;
  ip: string;
  status: 'online' | 'offline' | 'unknown';
  latencyMs: number;
  ttl: number;
  seq: number;
  method: string;
}

interface EngineLogTerminalProps {
  logs: EngineLogEntry[];
  onClear?: () => void;
}

export function EngineLogTerminal({ logs, onClear }: EngineLogTerminalProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);

  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  const handleScroll = () => {
    if (!scrollRef.current) return;
    const { scrollTop, scrollHeight, clientHeight } = scrollRef.current;
    const isAtBottom = scrollHeight - scrollTop - clientHeight < 40;
    setAutoScroll(isAtBottom);
  };

  return (
    <div className="bg-surface border border-border rounded-xl p-4 lg:p-5 h-full flex flex-col shadow-lg">
      <div className="flex items-center justify-between mb-3 pb-2.5 border-b border-border/60">
        <div className="flex items-center gap-2.5">
          <span className="text-xs font-mono font-semibold text-text-primary tracking-wide">
            MONITORING ENGINE LOGS
          </span>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-xs text-text-muted font-mono">{logs.length} events</span>
          {onClear && (
            <button
              onClick={onClear}
              className="px-2 py-0.5 text-[11px] font-mono rounded bg-bg hover:bg-surface-elevated text-text-muted hover:text-text-primary transition-colors border border-border"
              title="Clear terminal logs"
            >
              Clear
            </button>
          )}
        </div>
      </div>

      <div
        ref={scrollRef}
        onScroll={handleScroll}
        className="font-mono text-xs bg-bg/80 border border-border/80 rounded-lg p-3 h-[280px] lg:h-[360px] overflow-y-auto space-y-1.5 select-text"
      >
        {logs.length === 0 ? (
          <div className="h-full flex items-center justify-center text-text-muted">
            Menunggu aktivitas engine monitoring...
          </div>
        ) : (
          logs.map((log) => {
            const isOnline = log.status === 'online';
            const timeStr = (() => {
              try {
                const d = new Date(log.timestamp);
                return isNaN(d.getTime()) ? log.timestamp : d.toLocaleTimeString('id-ID');
              } catch {
                return log.timestamp;
              }
            })();

            return (
              <div
                key={log.id}
                className="flex flex-col sm:flex-row sm:items-center justify-between gap-1 py-1 px-2 rounded hover:bg-surface/60 transition-colors border-b border-border/20 last:border-0"
              >
                <div className="flex items-center gap-2 min-w-0">
                  <span className="text-text-muted shrink-0 text-[11px]">[{timeStr}]</span>
                  <span
                    className={`font-bold px-1.5 py-0.5 rounded text-[10px] uppercase tracking-wider shrink-0 ${
                      isOnline
                        ? 'bg-success-muted text-success'
                        : 'bg-danger-muted text-danger animate-pulse'
                    }`}
                  >
                    {log.status}
                  </span>
                  <span className="font-semibold text-text-primary truncate">{log.deviceName}</span>
                  <span className="text-text-muted text-[11px] truncate">({log.ip})</span>
                </div>

                <div className="flex items-center gap-3 text-[11px] pl-5 sm:pl-0 shrink-0 text-text-muted">
                  <span>method: <strong className="text-text-primary">{log.method || 'ICMP Ping'}</strong></span>
                  <span>seq: <strong className="text-text-primary">#{log.seq}</strong></span>
                  {isOnline ? (
                    <>
                      <span>ttl: <strong className="text-text-primary">{log.ttl}</strong></span>
                      <span className="text-success font-medium">+{log.latencyMs.toFixed(1)}ms</span>
                    </>
                  ) : (
                    <span className="text-danger font-medium">Request Timeout</span>
                  )}
                </div>
              </div>
            );
          })
        )}
      </div>
    </div>
  );
}
