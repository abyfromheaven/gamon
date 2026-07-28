import { useState, useMemo } from 'react';
import type { Alert, AlertStatus, AlertSeverity, DeviceType } from '../types';
import { alertData } from '../data/alerts';
import { PageHeader } from '../components/PageHeader';
import { AlertSummaryCards } from '../components/AlertSummaryCards';
import { AlertFilters } from '../components/AlertFilters';
import { AlertList } from '../components/AlertList';
import { AlertDetailPanel } from '../components/AlertDetailPanel';

export function AlertCenterPage() {
  const [alerts, setAlerts] = useState<Alert[]>(alertData);
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<AlertStatus | 'all'>('all');
  const [severityFilter, setSeverityFilter] = useState<AlertSeverity | 'all'>('all');
  const [deviceTypeFilter, setDeviceTypeFilter] = useState<DeviceType | 'all'>('all');
  const [selectedAlert, setSelectedAlert] = useState<Alert | null>(null);

  const filteredAlerts = useMemo(() => {
    return alerts.filter((alert) => {
      const matchesSearch =
        searchQuery === '' ||
        alert.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        alert.device.toLowerCase().includes(searchQuery.toLowerCase()) ||
        alert.description.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStatus = statusFilter === 'all' || alert.status === statusFilter;
      const matchesSeverity = severityFilter === 'all' || alert.severity === severityFilter;
      const matchesDeviceType = deviceTypeFilter === 'all' || alert.deviceType === deviceTypeFilter;

      return matchesSearch && matchesStatus && matchesSeverity && matchesDeviceType;
    });
  }, [alerts, searchQuery, statusFilter, severityFilter, deviceTypeFilter]);

  const handleMarkResolved = (id: number) => {
    setAlerts((prev) =>
      prev.map((a) =>
        a.id === id
          ? { ...a, status: 'resolved' as const, resolvedTime: new Date().toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' }) }
          : a
      )
    );
    setSelectedAlert((prev) =>
      prev && prev.id === id
        ? { ...prev, status: 'resolved' as const, resolvedTime: new Date().toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' }) }
        : prev
    );
  };

  return (
    <div className="p-4 lg:p-6 space-y-5 lg:space-y-6 max-w-7xl mx-auto">
      <PageHeader
        title="Alert Center"
        subtitle="Monitor and track all system alerts and incidents"
      />

      <AlertSummaryCards alerts={alerts} />

      <AlertFilters
        searchQuery={searchQuery}
        onSearchChange={setSearchQuery}
        statusFilter={statusFilter}
        onStatusFilterChange={setStatusFilter}
        severityFilter={severityFilter}
        onSeverityFilterChange={setSeverityFilter}
        deviceTypeFilter={deviceTypeFilter}
        onDeviceTypeFilterChange={setDeviceTypeFilter}
      />

      <div className="text-xs text-text-muted font-mono">
        Showing {filteredAlerts.length} of {alerts.length} alerts
      </div>

      <AlertList
        alerts={filteredAlerts}
        selectedAlert={selectedAlert}
        onSelectAlert={setSelectedAlert}
      />

      <AlertDetailPanel
        alert={selectedAlert}
        onClose={() => setSelectedAlert(null)}
        onMarkResolved={handleMarkResolved}
      />
    </div>
  );
}
