import { useEffect, useMemo, useRef, useState } from 'react';
import type { Device, MonitoringRecord, PingHistoryRecord } from '../lib/api';
import { fetchDeviceHistory, fetchDevices, fetchMonitoring } from '../lib/api';
import type { MonitorResult } from '../types';
import { PageHeader } from '../components/PageHeader';
import { StatusIndicator } from '../components/StatusIndicator';
import { LatencyChart } from '../components/LatencyChart';

export function MonitoringPage({ monitorResults, reconnectKey }: { monitorResults: Map<number, MonitorResult>; onViewAlerts: () => void; reconnectKey: number }) {
  const [devices, setDevices] = useState<Device[]>([]);
  const [records, setRecords] = useState<MonitoringRecord[]>([]);
  const [selected, setSelected] = useState<MonitoringRecord | null>(null);
  const [history, setHistory] = useState<PingHistoryRecord[]>([]);
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState<'all' | MonitoringRecord['status']>('all');
  const [typeFilter, setTypeFilter] = useState<'all' | Device['type']>('all');
  const [error, setError] = useState('');
  const historyCache = useRef<Map<number, PingHistoryRecord[]>>(new Map());

  useEffect(() => {
    void Promise.all([fetchDevices(), fetchMonitoring()])
      .then(([nextDevices, nextRecords]) => { setDevices(nextDevices); setRecords(nextRecords); })
      .catch((err: unknown) => setError(err instanceof Error ? err.message : 'Gagal memuat monitoring.'));
  }, [reconnectKey]);

  const live = useMemo(
    () => records.map((record) => {
      const result = monitorResults.get(record.device_id);
      return result ? { ...record, status: result.status, latency_ms: result.latency_ms, last_check: result.timestamp } : record;
    }),
    [records, monitorResults]
  );

  const filtered = live.filter(
    (record) =>
      (statusFilter === 'all' || record.status === statusFilter) &&
      (typeFilter === 'all' || record.device_type === typeFilter) &&
      (!search || `${record.device_name} ${record.ip}`.toLowerCase().includes(search.toLowerCase()))
  );

  const choose = async (record: MonitoringRecord) => {
    setSelected(record);
    const cached = historyCache.current.get(record.device_id);
    if (cached) {
      setHistory(cached);
      return;
    }
    try {
      const data = await fetchDeviceHistory(record.device_id);
      const sliced = data.slice(0, 50);
      historyCache.current.set(record.device_id, sliced);
      setHistory(sliced);
    } catch {
      setHistory([]);
    }
  };

  useEffect(() => {
    if (!selected) return;
    const deviceId = selected.device_id;
    const liveResult = monitorResults.get(deviceId);
    if (!liveResult) return;

    const newRecord: PingHistoryRecord = {
      id: Date.now(),
      status: liveResult.status,
      latency_ms: liveResult.latency_ms,
      ttl: liveResult.ttl,
      seq: liveResult.seq,
      details: JSON.stringify(liveResult.details),
      timestamp: liveResult.timestamp,
    };

    const cached = historyCache.current.get(deviceId) || [];
    const updated = [...cached, newRecord].slice(-50);
    historyCache.current.set(deviceId, updated);
    setHistory(updated);
  }, [monitorResults, selected?.device_id]);

  const chartData = useMemo(
    () => history
      .slice(-50)
      .reverse()
      .map((item) => {
        const date = new Date(item.timestamp);
        const time = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
        return { time, latency: item.latency_ms };
      }),
    [history]
  );

  const counts = {
    total: live.length,
    online: live.filter((item) => item.status === 'online').length,
    offline: live.filter((item) => item.status === 'offline').length,
  };

  const detail = selected && devices.find((device) => device.id === selected.device_id);

  return (
    <div className="max-w-[1200px] mx-auto px-4 lg:px-8 py-6 lg:py-8 space-y-5">
      <PageHeader
        title="Monitoring"
        subtitle={`${counts.total} devices · ${counts.online} online · ${counts.offline} offline`}
      />

      {error && <p className="rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}

      <div className="flex flex-wrap gap-3">
        <input
          className="min-w-56 flex-1 rounded-lg border border-border bg-surface p-3 text-sm"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Search device or IP..."
        />
        <select
          className="rounded-lg border border-border bg-surface px-3 text-sm"
          value={statusFilter}
          onChange={(event) => setStatusFilter(event.target.value as typeof statusFilter)}
        >
          <option value="all">All status</option>
          <option value="online">Online</option>
          <option value="offline">Offline</option>
          <option value="unknown">Unknown</option>
        </select>
        <select
          className="rounded-lg border border-border bg-surface px-3 text-sm"
          value={typeFilter}
          onChange={(event) => setTypeFilter(event.target.value as typeof typeFilter)}
        >
          <option value="all">All device types</option>
          {['Server', 'Router', 'Switch', 'Access Point', 'Website'].map((type) => (
            <option key={type}>{type}</option>
          ))}
        </select>
      </div>

      <div className="grid lg:grid-cols-3 gap-5">
        {/* Table */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-surface overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-text-muted">
                <th className="p-4">Device</th>
                <th>Status</th>
                <th>Latency</th>
                <th>Last Check</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((record) => (
                <tr
                  key={record.device_id}
                  onClick={() => void choose(record)}
                  className={`cursor-pointer border-t border-border/50 hover:bg-surface-elevated transition-colors ${
                    selected?.device_id === record.device_id ? 'bg-surface-elevated' : ''
                  }`}
                >
                  <td className="p-4">
                    <p className="font-medium">{record.device_name}</p>
                    <p className="font-mono text-xs text-text-muted">{record.ip} · {record.device_type}</p>
                  </td>
                  <td>
                    <StatusIndicator status={record.status} />
                  </td>
                  <td className="font-mono text-xs">{record.latency_ms} ms</td>
                  <td className="text-xs text-text-muted">
                    {record.last_check ? new Date(record.last_check).toLocaleTimeString('id-ID') : '—'}
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={4} className="p-8 text-center text-text-muted">Tidak ada device</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* Detail Panel */}
        <div className="rounded-xl border border-border bg-surface p-5 space-y-4">
          {detail && selected ? (
            <>
              <div>
                <h3 className="text-lg font-semibold">{detail.name}</h3>
                <p className="text-xs text-text-muted">{detail.type} · {detail.method}</p>
              </div>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-text-muted">IP</span>
                  <span className="font-mono">{detail.ip}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Status</span>
                  <StatusIndicator status={selected.status} />
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Latency</span>
                  <span className="font-mono">{selected.latency_ms} ms</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Interval</span>
                  <span>{selected.interval}s</span>
                </div>
                {detail.location && (
                  <div className="flex justify-between">
                    <span className="text-text-muted">Lokasi</span>
                    <span>{detail.location}</span>
                  </div>
                )}
                {detail.description && (
                  <div className="pt-2 border-t border-border/50">
                    <p className="text-text-muted text-xs mb-1">Deskripsi</p>
                    <p className="text-sm">{detail.description}</p>
                  </div>
                )}
              </div>

              {/* Latency Chart */}
              <div className="pt-3 border-t border-border/50">
                <p className="text-xs text-text-muted mb-3">Latency History (50 data terakhir)</p>
                <LatencyChart data={chartData} />
              </div>
            </>
          ) : (
            <div className="flex items-center justify-center h-full min-h-[300px] text-text-muted text-sm">
              Pilih device untuk melihat detail
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
