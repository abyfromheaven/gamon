// ============================================================
// MODULE API CLIENT - PENGKOMUNIKASIAN DENGAN BACKEND
// ============================================================
// Module ini berisi semua fungsi untuk berkomunikasi dengan server backend GAMON.
// Setiap fungsi mengirim permintaan (request) ke server dan menerima respons (response).
//
// Contoh penggunaan:
// - fetchDevices()     -> Mengambil daftar semua perangkat
// - fetchMonitoring()  -> Mengambil data status monitoring
// - fetchAlerts()      -> Mengambil daftar alert/peringatan
//
// Semua data dikembalikan dalam format JSON.
// Jika terjadi error, fungsi akan melempar APIError.
// ============================================================

// URL default backend server (bisa diubah via environment variable)
const DEFAULT_API_BASE_URL = 'http://localhost:8080';

const apiBaseURL = (import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL).replace(/\/$/, '');

// ============================================================
// TIPE DATA (TYPES)
// ============================================================

/** Jenis perangkat yang bisa dipantau */
export type DeviceType = 'Server' | 'Router' | 'Switch' | 'Access Point' | 'Website';

/** Metode pengecekan yang digunakan */
export type DeviceMethod = 'ICMP Ping' | 'HTTP Check' | 'TCP Port';

/** Status aktivitas perangkat di sistem */
export type DeviceStatus = 'active' | 'inactive';

/** Status hasil pengecekan (online/offline/tidak diketahui) */
export type MonitorStatus = 'online' | 'offline' | 'unknown';

/** Status alert/peringatan */
export type AlertStatus = 'ongoing' | 'resolved';

/** Tingkat keparahan alert */
export type AlertSeverity = 'low' | 'medium' | 'high' | 'critical' | 'info';

// ============================================================
// INTERFACE PERANGKAT (DEVICE)
// ============================================================

/** Data lengkap satu perangkat dari database */
export interface Device {
  id: number;                    // ID unik perangkat
  name: string;                  // Nama perangkat (contoh: "Router Utama")
  type: DeviceType;              // Jenis perangkat (Server, Router, dll)
  ip: string;                    // Alamat IP (contoh: "192.168.1.1")
  url: string;                   // URL website (jika metode HTTP)
  port: number | null;           // Nomor port (jika metode TCP)
  method: DeviceMethod;          // Metode pengecekan
  location: string;              // Lokasi perangkat (contoh: "Ruang Server")
  check_interval: number;        // Interval pengecekan (dalam detik)
  status: DeviceStatus;          // Status: active/inactive
  description: string;           // Deskripsi tambahan
  created_at?: string;           // Kapan perangkat ditambahkan
  updated_at?: string;           // Kapan terakhir kali diubah
}

/** Data untuk membuat perangkat baru */
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

/** Data untuk memperbarui perangkat (semua field opsional) */
export type DeviceUpdate = Partial<DeviceInput>;

// ============================================================
// INTERFACE ALERT (PERINGATAN)
// ============================================================

/** Data alert/peringatan dari database */
export interface Alert {
  id: number;                    // ID unik alert
  device_id: number;             // ID perangkat yang bermasalah
  device_name: string;           // Nama perangkat
  device_type: DeviceType;       // Jenis perangkat
  device_ip: string;             // IP perangkat
  method: DeviceMethod;          // Metode pengecekan
  title: string;                 // Judul alert
  status: AlertStatus;           // Status: ongoing/resolved
  severity: AlertSeverity;       // Tingkat keparahan
  started_at: string;            // Kapan alert mulai
  resolved_at: string | null;    // Kapan alert teratasi
  description: string;           // Deskripsi alert
  acknowledged: boolean;         // Apakah sudah dikonfirmasi admin
  acknowledged_at: string | null; // Kapan dikonfirmasi
}

/** Filter untuk pencarian alert */
export interface AlertFilters {
  status?: AlertStatus;
  severity?: AlertSeverity;
  device_type?: DeviceType;
}

// ============================================================
// INTERFACE DASHBOARD
// ============================================================

/** Ringkasan data untuk dashboard utama */
export interface DashboardSummary {
  total_devices: number;    // Total semua perangkat
  online_devices: number;   // Perangkat yang online
  offline_devices: number;  // Perangkat yang offline
}

/** Data alert terbaru untuk ditampilkan di dashboard */
export interface DashboardAlert {
  id: number;
  device_name: string;
  title: string;
  severity: AlertSeverity;
  status: AlertStatus;
  started_at: string;
}

/** Data lengkap dashboard */
export interface Dashboard {
  summary: DashboardSummary;
  latest_alerts: DashboardAlert[];
}

// ============================================================
// INTERFACE MONITORING
// ============================================================

/** Data status monitoring satu perangkat */
export interface MonitoringRecord {
  device_id: number;                  // ID perangkat
  device_name: string;                // Nama perangkat
  device_type: DeviceType;            // Jenis perangkat
  ip: string;                         // Alamat IP
  method: DeviceMethod;               // Metode pengecekan
  status: MonitorStatus;              // Status: online/offline/unknown
  latency_ms: number;                 // Waktu respons (milidetik)
  last_check: string | null;          // Kapan terakhir kali dicek
  last_online: string | null;         // Kapan terakhir kali online
  current_status_since: string;       // Sejak kapan status saat ini aktif
  current_status_duration: string;    // Berapa lama status saat ini berlangsung
  interval: number;                   // Interval pengecekan (detik)
}

/** Data riwayat pengecekan (ping history) */
export interface PingHistoryRecord {
  id: number;         // ID unik record
  status: string;     // Status saat itu (online/offline)
  latency_ms: number; // Waktu respons saat itu
  ttl: number;        // TTL paket ping
  seq: number;        // Nomor urut pengecekan
  details: string;    // Detail tambahan (JSON)
  timestamp: string;  // Kapan pengecekan dilakukan
}

// ============================================================
// FUNGSI-FUNGSI API
// ============================================================

// Struktur respons dari server
interface DataResponse<T> {
  success: boolean;
  data: T;
}

interface MessageResponse {
  success: boolean;
  message: string;
}

// Error khusus untuk masalah API
export class APIError extends Error {
  readonly status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = 'APIError';
    this.status = status;
  }
}

// Fungsi untuk mengecek apakah respons berhasil
function isSuccessResponse(payload: unknown): payload is { success: boolean } {
  return typeof payload === 'object' && payload !== null && 'success' in payload;
}

// Fungsi utama untuk mengirim permintaan ke server
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
    const message = 'message' in (payload as object) && typeof (payload as MessageResponse).message === 'string'
      ? (payload as MessageResponse).message
      : `Request gagal (${response.status}).`;
    throw new APIError(message, response.status);
  }

  if (expectsData && !('data' in payload)) {
    throw new APIError('Respons backend tidak memiliki data yang diharapkan.', response.status);
  }

  return ('data' in payload ? (payload as DataResponse<T>).data : undefined) as T;
}

// Fungsi pembuat request dengan body JSON
function jsonRequest(method: 'POST' | 'PUT' | 'DELETE', body?: unknown): RequestInit {
  return {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  };
}

// ============================================================
// FUNGSI-FUNGSI UNTUK PERANGKAT (DEVICES)
// ============================================================

/** Mengambil daftar semua perangkat dari server */
export function fetchDevices(): Promise<Device[]> {
  return request<Device[]>('/api/devices');
}

/** Membuat perangkat baru */
export function createDevice(device: DeviceInput): Promise<Device> {
  return request<Device>('/api/devices', jsonRequest('POST', device));
}

/** Memperbarui data perangkat */
export function updateDevice(id: number, device: DeviceUpdate): Promise<Device> {
  return request<Device>(`/api/devices/${id}`, jsonRequest('PUT', device));
}

/** Menghapus perangkat */
export function deleteDevice(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}`, jsonRequest('DELETE'), false);
}

/** Memulai pemantauan untuk satu perangkat */
export function startMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/start`, jsonRequest('POST'), false);
}

/** Menghentikan pemantauan untuk satu perangkat */
export function stopMonitor(id: number): Promise<void> {
  return request<void>(`/api/devices/${id}/stop`, jsonRequest('POST'), false);
}

/** Mengubah status perangkat (active/inactive) */
export function toggleDeviceStatus(id: number, status: DeviceStatus): Promise<{ id: number; status: DeviceStatus; message: string }> {
  return request<{ id: number; status: DeviceStatus; message: string }>(`/api/devices/${id}/status`, jsonRequest('PUT', { status }));
}

// ============================================================
// FUNGSI-FUNGSI UNTUK ALERT (PERINGATAN)
// ============================================================

/** Mengambil daftar alert dengan filter */
export function fetchAlerts(filters: AlertFilters = {}): Promise<Alert[]> {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined) query.set(key, value);
  }

  const suffix = query.size > 0 ? `?${query}` : '';
  return request<Alert[]>(`/api/alerts${suffix}`);
}

/** Menandai alert sebagai "teratasi" */
export function resolveAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/resolve`, jsonRequest('PUT'), false);
}

/** Mengonfirmasi alert (acknowledge) */
export function acknowledgeAlert(id: number): Promise<void> {
  return request<void>(`/api/alerts/${id}/acknowledge`, jsonRequest('PUT'), false);
}

/** Mengambil jumlah alert yang masih aktif */
export function fetchAlertCount(): Promise<{ ongoing: number }> {
  return request<{ ongoing: number }>('/api/alerts/count');
}

// ============================================================
// FUNGSI-FUNGSI UNTUK DASHBOARD
// ============================================================

/** Mengambil data ringkasan untuk dashboard */
export function fetchDashboard(): Promise<Dashboard> {
  return request<Dashboard>('/api/dashboard');
}

// ============================================================
// FUNGSI-FUNGSI UNTUK MONITORING
// ============================================================

/** Mengambil data status monitoring semua perangkat */
export function fetchMonitoring(): Promise<MonitoringRecord[]> {
  return request<MonitoringRecord[]>('/api/monitoring');
}

/** Mengambil riwayat pengecekan untuk satu perangkat */
export function fetchDeviceHistory(id: number): Promise<PingHistoryRecord[]> {
  return request<PingHistoryRecord[]>(`/api/monitoring/${id}/history`);
}

// ============================================================
// FUNGSI-FUNGSI UNTUK TELEGRAM BOT
// ============================================================

/** Status koneksi Telegram Bot */
export interface TelegramStatus {
  status: string;           // Status koneksi
  chat_id: string;          // ID chat Telegram
  paired_at: string | null; // Kapan terpasang
}

/** Token untuk pairing Telegram Bot */
export interface PairingToken {
  token: string;         // Token unik untuk pairing
  expires_at: string;    // Kapan token kadaluarsa
}

/** Membuat token baru untuk pairing Telegram */
export function generatePairingToken(): Promise<PairingToken> {
  return request<PairingToken>('/api/telegram/pair', jsonRequest('POST'));
}

/** Mengambil status koneksi Telegram */
export function getTelegramStatus(): Promise<TelegramStatus> {
  return request<TelegramStatus>('/api/telegram/status');
}

/** Memutus koneksi Telegram Bot */
export function disconnectTelegram(): Promise<void> {
  return request<void>('/api/telegram/disconnect', jsonRequest('DELETE'), false);
}

// ============================================================
// FUNGSI-FUNGSI UNTUK PENGATURAN (SETTINGS)
// ============================================================

/** Pengaturan aplikasi */
export interface AppSettings {
  failure_threshold: number;       // Batas toleransi kegagalan ping
  notifications_enabled: boolean;  // Apakah notifikasi aktif
}

/** Mengambil pengaturan dari server */
export function getSettings(): Promise<AppSettings> {
  return request<AppSettings>('/api/settings');
}

/** Memperbarui pengaturan di server */
export function updateSettings(settings: Partial<AppSettings>): Promise<void> {
  return request<void>('/api/settings', jsonRequest('PUT', settings), false);
}
