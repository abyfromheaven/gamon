// ============================================================================
// MODUL LAYANAN API FRONTEND (src/lib/api.ts)
// ============================================================================
// Modul ini bertugas melakukan komunikasi HTTP (REST API) dari browser React ke
// server backend Golang yang berjalan di port 8080.
// Seluruh fungsi pengiriman data (fetch, add, edit, delete, settings) didefinisikan
// di berkas ini agar rapi dan tidak berantakan (Clean Code & Modular).
// ============================================================================

const DEFAULT_API_BASE_URL = 'http://localhost:8080';

// Alamat URL API Backend Golang
const apiBaseURL = (import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL).replace(/\/$/, '');

// Tipe Data Standar Aplikasi GAMON
export type DeviceType = 'Server' | 'Router' | 'Switch' | 'Access Point' | 'Website';
export type DeviceMethod = 'ICMP Ping' | 'HTTP Check' | 'TCP Port';
export type DeviceStatus = 'active' | 'inactive';
export type MonitorStatus = 'online' | 'offline' | 'unknown';
export type AlertStatus = 'ongoing' | 'resolved';

// Interface Perangkat Jaringan (Device)
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

// Interface Form Tambah Perangkat
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

// Interface Peringatan (Alert)
export interface Alert {
  id: number;
  device_id: number;
  device_name: string;
  device_type: DeviceType;
  device_ip: string;
  method: DeviceMethod;
  title: string;
  status: AlertStatus;
  started_at: string;
  resolved_at: string | null;
  description: string;
  acknowledged: boolean;
  acknowledged_at: string | null;
}

export interface AlertFilters {
  status?: AlertStatus;
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

// Class penanganan error API kustom
export class APIError extends Error {
  readonly status: number;

  constructor(pesan: string, status: number) {
    super(pesan);
    this.name = 'APIError';
    this.status = status;
  }
}

function isSuccessResponse(payload: unknown): payload is { success: boolean } {
  return typeof payload === 'object' && payload !== null && 'success' in payload;
}

// Fungsi utama penanganan HTTP Request (Fetch Wrapper)
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

// --- FUNGSI FUNGSI PERMINTAAN API PERANGKAT (DEVICES) ---

// Mengambil seluruh daftar perangkat
export function fetchDevices(): Promise<Device[]> {
  return request<Device[]>('/api/devices');
}

// Menambah perangkat baru
export function createDevice(device: DeviceInput): Promise<Device> {
  return request<Device>('/api/devices', jsonRequest('POST', device));
}

// Memperbarui data perangkat
export function updateDevice(id: number, device: DeviceUpdate): Promise<Device> {
  return request<Device>(`/api/devices/${id}`, jsonRequest('PUT', device));
}

// Menghapus perangkat
export function deleteDevice(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}`, jsonRequest('DELETE'), false);
}

// Memulai pemantauan perangkat
export function startMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/start`, jsonRequest('POST'), false);
}

// Menghentikan pemantauan perangkat
export function stopMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/stop`, jsonRequest('POST'), false);
}

// Mengubah status aktif/nonaktif perangkat
export function toggleDeviceStatus(id: number, status: DeviceStatus): Promise<{ id: number; status: DeviceStatus; message: string }> {
  return request<{ id: number; status: DeviceStatus; message: string }>(`/api/devices/${id}/status`, jsonRequest('PUT', { status }));
}

// --- FUNGSI PERMINTAAN API PERINGATAN (ALERTS) ---

// Mengambil daftar peringatan
export function fetchAlerts(filters: AlertFilters = {}): Promise<Alert[]> {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined) query.set(key, value);
  }

  const suffix = query.size > 0 ? `?${query}` : '';
  return request<Alert[]>(`/api/alerts${suffix}`);
}

// Memulihkan peringatan
export function resolveAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/resolve`, jsonRequest('PUT'), false);
}

// Mengonfirmasi peringatan oleh admin
export function acknowledgeAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/acknowledge`, jsonRequest('PUT'), false);
}

// Menghitung jumlah peringatan aktif
export function fetchAlertCount(): Promise<{ ongoing: number }> {
  return request<{ ongoing: number }>('/api/alerts/count');
}

// --- FUNGSI PERMINTAAN DASHBOARD & PEMANTAUAN ---

// Mengambil data ringkasan dashboard
export function fetchDashboard(): Promise<Dashboard> {
  return request<Dashboard>('/api/dashboard');
}

// Mengambil status pemantauan seluruh perangkat
export function fetchMonitoring(): Promise<MonitoringRecord[]> {
  return request<MonitoringRecord[]>('/api/monitoring');
}

// Mengambil riwayat latensi ping 50 data terakhir suatu perangkat
export function fetchDeviceHistory(id: number): Promise<PingHistoryRecord[]> {
  return request<PingHistoryRecord[]>(`/api/monitoring/${id}/history`);
}

// --- FUNGSI INTEGRASI TELEGRAM BOT ---

export interface TelegramStatus {
  status: string;
  chat_id: string;
  paired_at: string | null;
}

export interface PairingToken {
  token: string;
  expires_at: string;
}

// Menghasilkan token hubung Telegram baru
export function generatePairingToken(): Promise<PairingToken> {
  return request<PairingToken>('/api/telegram/pair', jsonRequest('POST'));
}

// Mengecek status koneksi Telegram
export function getTelegramStatus(): Promise<TelegramStatus> {
  return request<TelegramStatus>('/api/telegram/status');
}

// Memutus hubungan Telegram
export function disconnectTelegram(): Promise<void> {
  return request<void>('/api/telegram/disconnect', jsonRequest('DELETE'), false);
}

// --- FUNGSI PENGATURAN APLIKASI ---

export interface AppSettings {
  failure_threshold: number;
  notifications_enabled: boolean;
}

// Membaca pengaturan aplikasi
export function getSettings(): Promise<AppSettings> {
  return request<AppSettings>('/api/settings');
}

// Memperbarui pengaturan aplikasi
export function updateSettings(settings: Partial<AppSettings>): Promise<void> {
  return request<void>('/api/settings', jsonRequest('PUT', settings), false);
}
