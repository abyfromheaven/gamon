import type {
  AlertSeverity as APIAlertSeverity,
  AlertStatus as APIAlertStatus,
  DeviceType as APIDeviceType,
  MonitorStatus,
} from '../lib/api';

export type DeviceType = APIDeviceType;
export type AlertSeverity = APIAlertSeverity;
export type AlertStatus = APIAlertStatus;

export interface PingResult {
  ip: string;
  status: MonitorStatus;
  latency_ms: number;
  ttl: number;
  seq: number;
  timestamp: string;
}

export interface MonitorResult extends PingResult {
  device_id: number;
  method: string;
  details: Record<string, unknown>;
}

export interface StatusChange {
  device_id: number;
  device_name: string;
  old_status: MonitorStatus | '';
  new_status: MonitorStatus;
  timestamp: string;
}

export interface Alert {
  id: number;
  title: string;
  device: string;
  deviceType: DeviceType;
  status: AlertStatus;
  severity: AlertSeverity;
  startTime: string;
  resolvedTime: string | null;
  description: string;
  monitoringMethod: string;
  timestamp: number;
}

export interface DeviceBreakdown {
  type: DeviceType;
  count: number;
  online: number;
}

export interface SystemStatusInfo {
  monitoring: 'Running' | 'Stopped' | 'Error';
  checkInterval: string;
  lastScan: string;
  notifications: 'Active' | 'Paused';
}

export interface DashboardData {
  summary: {
    totalDevices: number;
    online: number;
    offline: number;
    warnings: number;
  };
  deviceBreakdown: DeviceBreakdown[];
  latestAlerts: Alert[];
  systemStatus: SystemStatusInfo;
}
