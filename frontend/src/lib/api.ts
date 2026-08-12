const DEFAULT_API_BASE_URL = 'http://localhost:8080';

const apiBaseURL = (import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL).replace(/\/$/, '');

export type DeviceType = 'Server' | 'Router' | 'Switch' | 'Access Point' | 'Website';
export type DeviceMethod = 'ICMP Ping' | 'HTTP Check' | 'TCP Port';
export type DeviceStatus = 'active' | 'inactive';
export type MonitorStatus = 'online' | 'offline' | 'unknown';
export type AlertStatus = 'ongoing' | 'resolved';
export type AlertSeverity = 'low' | 'medium' | 'high' | 'critical' | 'info';

export interface Device {
  id: number;
  name: string;
  type: DeviceType;
  ip: string;
  url: string;
  port: number | null;
  method: DeviceMethod;
  location: string;
  check_interval: number;
  status: DeviceStatus;
  description: string;
  created_at?: string;
  updated_at?: string;
}

export interface DeviceInput {
  name: string;
  type: DeviceType;
  ip: string;
  url?: string;
  port?: number | null;
  method?: DeviceMethod;
  location?: string;
  check_interval?: number;
  status?: DeviceStatus;
  description?: string;
}

export type DeviceUpdate = Partial<DeviceInput>;

export interface Alert {
  id: number;
  device_id: number;
  device_name: string;
  device_type: DeviceType;
  device_ip: string;
  method: DeviceMethod;
  title: string;
  status: AlertStatus;
  severity: AlertSeverity;
  started_at: string;
  resolved_at: string | null;
  description: string;
  acknowledged: boolean;
  acknowledged_at: string | null;
}

export interface AlertFilters {
  status?: AlertStatus;
  severity?: AlertSeverity;
  device_type?: DeviceType;
}

export interface DashboardSummary {
  total_devices: number;
  online_devices: number;
  offline_devices: number;
}

export interface DashboardAlert {
  id: number;
  device_name: string;
  title: string;
  severity: AlertSeverity;
  status: AlertStatus;
  started_at: string;
}

export interface Dashboard {
  summary: DashboardSummary;
  latest_alerts: DashboardAlert[];
}

export interface MonitoringRecord {
  device_id: number;
  device_name: string;
  device_type: DeviceType;
  ip: string;
  method: DeviceMethod;
  status: MonitorStatus;
  latency_ms: number;
  last_check: string | null;
  interval: number;
}

export interface PingHistoryRecord {
  id: number;
  status: MonitorStatus;
  latency_ms: number;
  ttl: number;
  seq: number;
  details: string;
  timestamp: string;
}

interface DataResponse<T> {
  success: boolean;
  data: T;
}

interface MessageResponse {
  success: boolean;
  message: string;
}

export class APIError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'APIError';
    this.status = status;
  }
}

function isSuccessResponse(payload: unknown): payload is { success: boolean } {
  return typeof payload === 'object' && payload !== null && 'success' in payload;
}

async function request<T>(path: string, init: RequestInit = {}, expectsData = true): Promise<T> {
  let response: Response;
  try {
    response = await fetch(`${apiBaseURL}${path}`, {
      ...init,
      headers: {
        Accept: 'application/json',
        ...init.headers,
      },
    });
  } catch {
    throw new APIError('Tidak dapat terhubung ke backend Gamon.', 0);
  }

  let payload: unknown;
  try {
    payload = await response.json();
  } catch {
    throw new APIError('Backend Gamon mengirim respons yang bukan JSON.', response.status);
  }

  if (!response.ok || !isSuccessResponse(payload) || !payload.success) {
    const message = 'message' in (payload as object) && typeof (payload as MessageResponse).message === 'string' ? (payload as MessageResponse).message : `Request gagal (${response.status}).`;
    throw new APIError(message, response.status);
  }

  if (expectsData && !('data' in payload)) {
    throw new APIError('Respons backend tidak memiliki data yang diharapkan.', response.status);
  }

  return ('data' in payload ? (payload as DataResponse<T>).data : undefined) as T;
}

function jsonRequest(method: 'POST' | 'PUT' | 'DELETE', body?: unknown): RequestInit {
  return {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  };
}

export function fetchDevices(): Promise<Device[]> {
  return request<Device[]>('/api/devices');
}

export function createDevice(device: DeviceInput): Promise<Device> {
  return request<Device>('/api/devices', jsonRequest('POST', device));
}

export function updateDevice(id: number, device: DeviceUpdate): Promise<Device> {
  return request<Device>(`/api/devices/${id}`, jsonRequest('PUT', device));
}

export function deleteDevice(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}`, jsonRequest('DELETE'), false);
}

export function startMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/start`, jsonRequest('POST'), false);
}

export function stopMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/stop`, jsonRequest('POST'), false);
}

export function toggleDeviceStatus(id: number, status: DeviceStatus): Promise<{ id: number; status: DeviceStatus; message: string }> {
  return request<{ id: number; status: DeviceStatus; message: string }>(`/api/devices/${id}/status`, jsonRequest('PUT', { status }));
}

export function fetchAlerts(filters: AlertFilters = {}): Promise<Alert[]> {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined) query.set(key, value);
  }

  const suffix = query.size > 0 ? `?${query}` : '';
  return request<Alert[]>(`/api/alerts${suffix}`);
}

export function resolveAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/resolve`, jsonRequest('PUT'), false);
}

export function acknowledgeAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/acknowledge`, jsonRequest('PUT'), false);
}

export function fetchAlertCount(): Promise<{ ongoing: number }> {
  return request<{ ongoing: number }>('/api/alerts/count');
}

export function fetchDashboard(): Promise<Dashboard> {
  return request<Dashboard>('/api/dashboard');
}

export function fetchMonitoring(): Promise<MonitoringRecord[]> {
  return request<MonitoringRecord[]>('/api/monitoring');
}

export function fetchDeviceHistory(id: number): Promise<PingHistoryRecord[]> {
  return request<PingHistoryRecord[]>(`/api/monitoring/${id}/history`);
}

export interface TelegramStatus {
  status: string;
  chat_id: string;
  paired_at: string | null;
}

export interface PairingToken {
  token: string;
  expires_at: string;
}

export function generatePairingToken(): Promise<PairingToken> {
  return request<PairingToken>('/api/telegram/pair', jsonRequest('POST'));
}

export function getTelegramStatus(): Promise<TelegramStatus> {
  return request<TelegramStatus>('/api/telegram/status');
}

export function disconnectTelegram(): Promise<void> {
  return request<void>('/api/telegram/disconnect', jsonRequest('DELETE'), false);
}

export interface AppSettings {
  failure_threshold: number;
  notifications_enabled: boolean;
}

export function getSettings(): Promise<AppSettings> {
  return request<AppSettings>('/api/settings');
}

export function updateSettings(settings: Partial<AppSettings>): Promise<void> {
  return request<void>('/api/settings', jsonRequest('PUT', settings), false);
}
