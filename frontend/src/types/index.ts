export interface PingResult {
  ip: string;
  status: 'online' | 'down';
  latency: number;
  ttl: number;
  seq: number;
  timestamp: string;
}

export interface WSMessage {
  type: 'ping_result' | 'status_change';
  data: PingResult;
}

export interface MonitorRequest {
  ip: string;
}

export interface APIResponse {
  success: boolean;
  message: string;
}

export type DeviceType = 'Server' | 'Router' | 'Switch' | 'Access Point' | 'Website';

export type DeviceMethod = 'ICMP Ping';

export interface Device {
  id: number;
  name: string;
  type: DeviceType;
  ip: string;
  method: DeviceMethod;
  port: number | null;
  location: string;
  status: 'active' | 'inactive';
  lastSeen: string | null;
}

export type AlertSeverity = 'critical' | 'warning' | 'info';

export interface DeviceBreakdown {
  type: DeviceType;
  count: number;
  online: number;
}

export interface Alert {
  id: number;
  device: string;
  message: string;
  severity: AlertSeverity;
  time: string;
  timestamp: number;
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
