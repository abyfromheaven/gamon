import { useEffect, useMemo, useRef, useState } from 'react';
import type { Dashboard, MonitoringRecord } from '../lib/api';
import { fetchDashboard, fetchMonitoring } from '../lib/api';
import { presentAlert } from '../lib/presenters';
import type { DashboardData, MonitorResult, StatusChange } from '../types';
import type { Page } from '../components/Sidebar';
import { MetricsGrid } from '../components/MetricsGrid';
import { LatestAlerts } from '../components/LatestAlerts';
import { SystemStatus } from '../components/SystemStatus';
import { EngineLogTerminal, type EngineLogEntry } from '../components/EngineLogTerminal';

interface DashboardPageProps {
  monitorResults: Map<number, MonitorResult>;
  onNavigate: (page: Page) => void;
  isConnected: boolean;
  reconnectKey: number;
  lastStatusChange?: StatusChange | null;
}

export function DashboardPage({ monitorResults, onNavigate, isConnected, reconnectKey, lastStatusChange }: DashboardPageProps) {
  const [dashboard, setDashboard] = useState<Dashboard | null>(null);
  const [monitoring, setMonitoring] = useState<MonitoringRecord[]>([]);
  const [logs, setLogs] = useState<EngineLogEntry[]>(() => {
    try {
      const cached = localStorage.getItem('gamon_engine_logs');
      if (cached) {
        const parsed = JSON.parse(cached) as EngineLogEntry[];
        if (Array.isArray(parsed)) return parsed.slice(-50);
      }
    } catch {
      // ignore
    }
    return [];
  });
  const [error, setError] = useState('');
  const lastProcessedResultRef = useRef<Map<number, string>>(new Map());

  const load = async () => {
    try {
      setError('');
      const [nextDashboard, nextMonitoring] = await Promise.all([
        fetchDashboard(),
        fetchMonitoring(),
      ]);
      setDashboard(nextDashboard);
      setMonitoring(nextMonitoring);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal memuat dashboard.');
    }
  };

  useEffect(() => {
    void load();
  }, [reconnectKey]);

  useEffect(() => {
    if (!lastStatusChange) return;
    void load();
  }, [lastStatusChange]);

  useEffect(() => {
    try {
      localStorage.setItem('gamon_engine_logs', JSON.stringify(logs.slice(-50)));
    } catch {
      // ignore
    }
  }, [logs]);

  const liveMonitoring = useMemo(() => {
    return monitoring.map((record) => {
      const live = monitorResults.get(record.device_id);
      return live
        ? { ...record, status: live.status, latency_ms: live.latency_ms, last_check: live.timestamp }
        : record;
    });
  }, [monitoring, monitorResults]);

  const deviceNameMap = useMemo(() => {
    const map = new Map<number, { name: string; ip: string; method: string }>();
    for (const r of monitoring) {
      map.set(r.device_id, { name: r.device_name, ip: r.ip, method: r.method });
    }
    return map;
  }, [monitoring]);

  useEffect(() => {
    if (monitoring.length > 0 && logs.length === 0) {
      const initialLogs: EngineLogEntry[] = monitoring.map((m) => ({
        id: `init-${m.device_id}`,
        timestamp: m.last_check || new Date().toISOString(),
        deviceName: m.device_name,
        ip: m.ip,
        status: m.status,
        latencyMs: m.latency_ms,
        ttl: 64,
        seq: 1,
        method: m.method,
      }));
      setLogs(initialLogs.slice(-50));
    }
  }, [monitoring]);

  useEffect(() => {
    if (monitorResults.size === 0) return;
    const newEntries: EngineLogEntry[] = [];

    for (const [deviceId, res] of monitorResults.entries()) {
      const lastTs = lastProcessedResultRef.current.get(deviceId);
      if (lastTs !== res.timestamp) {
        lastProcessedResultRef.current.set(deviceId, res.timestamp);
        const devInfo = deviceNameMap.get(deviceId) || { name: `Device #${deviceId}`, ip: res.ip, method: res.method };
        newEntries.push({
          id: `${deviceId}-${res.timestamp}-${res.seq}-${Math.random()}`,
          timestamp: res.timestamp || new Date().toISOString(),
          deviceName: devInfo.name,
          ip: devInfo.ip,
          status: res.status,
          latencyMs: res.latency_ms,
          ttl: res.ttl,
          seq: res.seq,
          method: devInfo.method || res.method || 'ICMP Ping',
        });
      }
    }

    if (newEntries.length > 0) {
      setLogs((prev) => [...prev, ...newEntries].slice(-50));
    }
  }, [monitorResults, deviceNameMap]);

  const data: DashboardData | null = useMemo(() => {
    if (!dashboard) return null;
    return {
      summary: {
        totalDevices: liveMonitoring.length || dashboard.summary.total_devices,
        online: liveMonitoring.filter((item) => item.status === 'online').length,
        offline: liveMonitoring.filter((item) => item.status === 'offline').length,
      },
      deviceBreakdown: [],
      latestAlerts: dashboard.latest_alerts.map(presentAlert),
      systemStatus: {
        monitoring: isConnected ? 'Running' : 'Stopped',
        checkInterval: liveMonitoring.length
          ? `${Math.min(...liveMonitoring.map((item) => item.interval))} seconds`
          : '—',
        lastScan: liveMonitoring.find((item) => item.last_check)?.last_check
          ? new Date(liveMonitoring.find((item) => item.last_check)!.last_check!).toLocaleTimeString('id-ID')
          : '—',
        notifications: isConnected ? 'Active' : 'Paused',
      },
    };
  }, [dashboard, liveMonitoring, isConnected]);

  if (!data) {
    return <div className="p-8 text-text-muted">{error || 'Loading dashboard...'}</div>;
  }

  return (
    <div className="max-w-[1600px] mx-auto px-4 lg:px-8 py-4 lg:py-5 h-full flex flex-col space-y-4 overflow-hidden">
      {error && <p className="rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}
      <MetricsGrid summary={data.summary} />
      
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-4 items-stretch flex-1 min-h-0">
        <div className="lg:col-span-8 h-full flex flex-col min-h-0">
          <EngineLogTerminal logs={logs} onClear={() => { setLogs([]); localStorage.removeItem('gamon_engine_logs'); }} />
        </div>
        <div className="lg:col-span-4 space-y-4 h-full flex flex-col min-h-0">
          <SystemStatus status={data.systemStatus} />
          <div className="flex-1 min-h-0 flex flex-col">
            <LatestAlerts alerts={data.latestAlerts} onViewAll={() => onNavigate('alerts')} />
          </div>
        </div>
      </div>
    </div>
  );
}

