import { useCallback, useEffect, useMemo, useState } from 'react';
import { fetchAlerts, resolveAlert } from '../lib/api';
import { presentAlert } from '../lib/presenters';
import type { Alert, AlertStatus, AlertSeverity, DeviceType, StatusChange } from '../types';
import { PageHeader } from '../components/PageHeader';
import { AlertSummaryCards } from '../components/AlertSummaryCards';
import { AlertFilters } from '../components/AlertFilters';
import { AlertList } from '../components/AlertList';
import { AlertDetailPanel } from '../components/AlertDetailPanel';

export function AlertCenterPage({ lastStatusChange }: { lastStatusChange: StatusChange | null }) {
  const [alerts, setAlerts] = useState<Alert[]>([]); const [selected, setSelected] = useState<Alert | null>(null); const [error, setError] = useState('');
  const [search, setSearch] = useState(''); const [status, setStatus] = useState<AlertStatus | 'all'>('all'); const [severity, setSeverity] = useState<AlertSeverity | 'all'>('all'); const [deviceType, setDeviceType] = useState<DeviceType | 'all'>('all');
  const load = useCallback(async () => { try { setError(''); const data = await fetchAlerts({ ...(status !== 'all' && { status }), ...(severity !== 'all' && { severity }), ...(deviceType !== 'all' && { device_type: deviceType }) }); setAlerts(data.map(presentAlert)); } catch (err) { setError(err instanceof Error ? err.message : 'Gagal memuat alert.'); } }, [status, severity, deviceType]);
  useEffect(() => { void load(); }, [load]);
  useEffect(() => { if (lastStatusChange) void load(); }, [lastStatusChange, load]);
  const filtered = useMemo(() => alerts.filter((alert) => !search || `${alert.title} ${alert.device} ${alert.description}`.toLowerCase().includes(search.toLowerCase())), [alerts, search]);
  const markResolved = async (id: number) => { try { await resolveAlert(id); await load(); setSelected((current) => current?.id === id ? { ...current, status: 'resolved', resolvedTime: new Date().toLocaleString('id-ID') } : current); } catch (err) { setError(err instanceof Error ? err.message : 'Gagal menyelesaikan alert.'); } };
  return <div className="p-4 lg:p-6 space-y-5 lg:space-y-6 max-w-7xl mx-auto"><PageHeader title="Alert Center" subtitle="Monitor and track all system alerts and incidents" />{error && <p className="rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}<AlertSummaryCards alerts={alerts} /><AlertFilters searchQuery={search} onSearchChange={setSearch} statusFilter={status} onStatusFilterChange={setStatus} severityFilter={severity} onSeverityFilterChange={setSeverity} deviceTypeFilter={deviceType} onDeviceTypeFilterChange={setDeviceType} /><div className="text-xs text-text-muted font-mono">Showing {filtered.length} of {alerts.length} alerts</div><AlertList alerts={filtered} selectedAlert={selected} onSelectAlert={setSelected} /><AlertDetailPanel alert={selected} onClose={() => setSelected(null)} onMarkResolved={(id) => void markResolved(id)} /></div>;
}
