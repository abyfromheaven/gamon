import { useState, useEffect } from 'react';
import { ConnectionBadge } from './ConnectionBadge';

interface TopBarProps {
  isConnected: boolean;
}

export function TopBar({ isConnected }: TopBarProps) {
  const [time, setTime] = useState(new Date());

  useEffect(() => {
    const interval = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(interval);
  }, []);

  const formattedTime = time.toLocaleTimeString('id-ID', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });

  const formattedDate = time.toLocaleDateString('id-ID', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });

  return (
    <header className="animate-fade-in sticky top-0 z-40 bg-bg/80 backdrop-blur-sm border-b border-border">
      <div className="max-w-[1400px] mx-auto px-4 lg:px-8 h-14 flex items-center justify-between">
        {/* Left: Brand */}
        <div className="flex items-center gap-3">
          <img src="/gamon-logo.svg" alt="Gamon" className="w-8 h-8 text-accent" />
          <div className="flex items-baseline gap-2">
            <span className="font-semibold text-text-primary tracking-tight">GAMON</span>
            <span className="text-text-muted text-xs hidden sm:inline">Garda Monitoring</span>
          </div>
        </div>

        {/* Center: Clock */}
        <div className="hidden md:flex flex-col items-center">
          <span className="font-mono text-lg font-semibold text-text-primary tracking-wide">
            {formattedTime}
          </span>
          <span className="text-[10px] text-text-muted uppercase tracking-wider">
            {formattedDate}
          </span>
        </div>

        {/* Right: Connection Status */}
        <ConnectionBadge isConnected={isConnected} />
      </div>
    </header>
  );
}
