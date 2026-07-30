import type { Alert as APIAlert, DashboardAlert, MonitoringRecord } from './api';
import type { Alert, DeviceBreakdown, DeviceType } from '../types';

const deviceTypes: DeviceType[] = ['Server', 'Router', 'Switch', 'Access Point', 'Website'];

export function formatDateTime(value: string | null): string | null {
  if (!value) return null;
  const date = new Date(value);
  return Number.isNaN(date.valueOf())
    ? value
    : date.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' });
}

export function presentAlert(alert: APIAlert | DashboardAlert): Alert {
  const resolvedAt = 'resolved_at' in alert ? alert.resolved_at : null;
  const method = 'method' in alert ? alert.method : 'ICMP Ping';
  const description = 'description' in alert ? alert.description : '';
  const deviceType = 'device_type' in alert ? alert.device_type : 'Server';
  return {
    id: alert.id,
    title: alert.title,
    device: alert.device_name,
    deviceType,
    status: alert.status,
    severity: alert.severity,
    startTime: formatDateTime(alert.started_at) ?? alert.started_at,
    resolvedTime: formatDateTime(resolvedAt),
    description,
    monitoringMethod: method,
    timestamp: new Date(alert.started_at).valueOf(),
  };
}

export function deviceBreakdown(records: MonitoringRecord[]): DeviceBreakdown[] {
  return deviceTypes.map((type) => {
    const devices = records.filter((record) => record.device_type === type);
    return {
      type,
      count: devices.length,
      online: devices.filter((record) => record.status === 'online').length,
    };
  });
}
