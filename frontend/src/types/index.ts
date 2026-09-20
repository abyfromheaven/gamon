// ============================================================
// INTERFACE TYPES - TIPE DATA TERPUSAT
// ============================================================
// File ini berisi semua definisi tipe data yang digunakan di seluruh aplikasi frontend.
// Tipe data ini memastikan data yang dikirim dari server sesuai dengan yang diharapkan.
//
// Catatan: Beberapa tipe diimpor dari api.ts untuk menghindari duplikasi.
// ============================================================

import type {
  AlertSeverity as APIAlertSeverity,
  AlertStatus as APIAlertStatus,
  DeviceType as APIDeviceType,
  MonitorStatus,
} from '../lib/api';

// Re-export tipe dari api.ts agar bisa digunakan di file lain
export type DeviceType = APIDeviceType;
export type AlertSeverity = APIAlertSeverity;
export type AlertStatus = APIAlertStatus;

// ============================================================
// INTERFACE HASIL PING
// ============================================================

/** Hasil pengecekan ping ke satu perangkat */
export interface PingResult {
  ip: string;              // Alamat IP yang dicek
  status: MonitorStatus;   // Status hasil: online/offline/unknown
  latency_ms: number;      // Waktu respons dalam milidetik
  ttl: number;             // Time To Live (TTL) paket ping
  seq: number;             // Nomor urut pengecekan
  timestamp: string;       // Kapan pengecekan dilakukan
}

// ============================================================
// INTERFACE HASIL MONITORING
// ============================================================

/** Hasil monitoring dari WebSocket yang berisi data lengkap perangkat */
export interface MonitorResult extends PingResult {
  device_id: number;                    // ID perangkat di database
  method: string;                       // Metode pengecekan (ICMP Ping, TCP, HTTP)
  details: Record<string, unknown>;     // Detail tambahan (bisa berisi apa saja)
  last_online?: string;                 // Kapan terakhir kali perangkat online
}

// ============================================================
// INTERFACE PERUBAHAN STATUS
// ============================================================

/** Data perubahan status perangkat (contoh: dari online ke offline) */
export interface StatusChange {
  device_id: number;              // ID perangkat
  device_name: string;            // Nama perangkat
  old_status: MonitorStatus | ''; // Status sebelumnya (kosong jika pertama kali)
  new_status: MonitorStatus;      // Status yang baru
  timestamp: string;              // Kapan perubahan terjadi
}

// ============================================================
// INTERFACE ALERT (PERINGATAN)
// ============================================================

/** Data alert/peringatan yang ditampilkan di frontend */
export interface Alert {
  id: number;                     // ID unik alert
  title: string;                  // Judul alert
  device: string;                 // Nama perangkat
  deviceType: DeviceType;         // Jenis perangkat
  status: AlertStatus;            // Status: ongoing/resolved
  severity: AlertSeverity;        // Tingkat keparahan
  startTime: string;              // Kapan alert mulai
  resolvedTime: string | null;    // Kapan alert teratasi
  description: string;            // Deskripsi alert
  monitoringMethod: string;       // Metode pengecekan
  acknowledged: boolean;          // Apakah sudah dikonfirmasi admin
  acknowledgedAt: string | null;  // Kapan dikonfirmasi
  timestamp: number;              // Timestamp untuk pengurutan
}

// ============================================================
// INTERFACE BREAKDOWN PERANGKAT
// ============================================================

/** Ringkasan jumlah perangkat berdasarkan jenis */
export interface DeviceBreakdown {
  type: DeviceType;   // Jenis perangkat
  count: number;      // Jumlah total
  online: number;     // Jumlah yang online
}

// ============================================================
// INTERFACE STATUS SISTEM
// ============================================================

/** Informasi status sistem monitoring */
export interface SystemStatusInfo {
  monitoring: 'Running' | 'Stopped' | 'Error';  // Status monitoring
  checkInterval: string;                          // Interval pengecekan
  lastScan: string;                               // Kapan terakhir kali scan
  notifications: 'Active' | 'Paused';             // Status notifikasi
}

// ============================================================
// INTERFACE DATA DASHBOARD
// ============================================================

/** Data lengkap untuk ditampilkan di dashboard */
export interface DashboardData {
  summary: {
    totalDevices: number;   // Total perangkat
    online: number;         // Perangkat online
    offline: number;        // Perangkat offline
  };
  deviceBreakdown: DeviceBreakdown[];  // Breakdown berdasarkan jenis
  latestAlerts: Alert[];               // Alert terbaru
  systemStatus: SystemStatusInfo;      // Status sistem
}
