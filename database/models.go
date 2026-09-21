package database

// Module Models berisi definisi struktur data (struct) yang merepresentasikan tabel-tabel
// di Database SQLite serta format data JSON yang dikirimkan ke frontend (React).

import "time"

// Device merepresentasikan tabel 'devices' (Perangkat Jaringan yang dipantau).
type Device struct {
	ID            int       `json:"id"`             // ID unik perangkat (Primary Key)
	Name          string    `json:"name"`           // Nama perangkat (contoh: Router Utama, Server SIAKAD)
	Type          string    `json:"type"`           // Jenis perangkat (Server, Router, Switch, Access Point, Website)
	IP            string    `json:"ip"`             // Alamat IP perangkat (contoh: 192.168.1.1)
	URL           string    `json:"url"`            // URL (opsional, jika memantau website HTTP)
	Port          *int      `json:"port"`           // Port spesifik (opsional)
	Method        string    `json:"method"`         // Metode pengecekan (ICMP Ping, TCP Port, HTTP GET)
	Location      string    `json:"location"`       // Lokasi fisik/ruangan perangkat
	CheckInterval int       `json:"check_interval"` // Interval waktu pengecekan (dalam detik)
	Status        string    `json:"status"`         // Status pemantauan (active / inactive)
	Description   string    `json:"description"`    // Catatan tambahan mengenai perangkat
	CreatedAt     time.Time `json:"created_at"`     // Waktu pertama kali ditambahkan
	UpdatedAt     time.Time `json:"updated_at"`     // Waktu terakhir kali diubah
}

// PingHistory merepresentasikan tabel 'ping_history' (Riwayat Pengecekan Latensi & Status).
type PingHistory struct {
	ID        int       `json:"id"`         // ID unik riwayat
	DeviceID  int       `json:"device_id"`  // ID perangkat yang di-ping (Foreign Key ke devices)
	Status    string    `json:"status"`     // Hasil status (online / offline / warning)
	LatencyMs float64   `json:"latency_ms"` // Waktu respon dalam milidetik (ms)
	TTL       int       `json:"ttl"`        // Time To Live dari paket ICMP
	Seq       int       `json:"seq"`        // Nomor urut paket ping
	Details   string    `json:"details"`    // Detail tambahan (misal pesan error)
	Timestamp time.Time `json:"timestamp"`  // Waktu saat pengecekan dilakukan
}

// Alert merepresentasikan tabel 'alerts' (Catatan Peringatan saat Perangkat Bermasalah/Mati).
type Alert struct {
	ID             int        `json:"id"`              // ID unik peringatan
	DeviceID       int        `json:"device_id"`       // ID perangkat bermasalah
	Title          string     `json:"title"`           // Judul peringatan (contoh: Perangkat Tidak Merespon)
	Status         string     `json:"status"`          // Status peringatan (ongoing / resolved)
	AlertType      string     `json:"alert_type"`      // Jenis peringatan
	StartedAt      time.Time  `json:"started_at"`      // Waktu mulai terjadi masalah
	ResolvedAt     *time.Time `json:"resolved_at"`     // Waktu masalah berhasil teratasi (opsional)
	Description    string     `json:"description"`     // Rincian penyebab masalah
	Acknowledged   bool       `json:"acknowledged"`    // Apakah sudah dikonfirmasi oleh admin
	AcknowledgedAt *time.Time `json:"acknowledged_at"` // Waktu konfirmasi admin
}

// DeviceWithType merepresentasikan data perangkat beserta jenis tipe pengelompokannya.
type DeviceWithType struct {
	Device
	TypeName string `json:"type_name"`
}

// MonitoringStatus merepresentasikan status pemantauan real-time yang dikirimkan via WebSocket.
type MonitoringStatus struct {
	DeviceID   int       `json:"device_id"`   // ID perangkat
	DeviceName string    `json:"device_name"` // Nama perangkat
	DeviceType string    `json:"device_type"` // Tipe perangkat
	IP         string    `json:"ip"`          // Alamat IP
	Method     string    `json:"method"`      // Metode pemantauan (ICMP Ping)
	Status     string    `json:"status"`      // Status saat ini (online/offline)
	LatencyMs  float64   `json:"latency_ms"`  // Latensi terkini dalam ms
	LastCheck  time.Time `json:"last_check"`  // Waktu pengecekan terakhir
	Interval   int       `json:"interval"`    // Frekuensi pengecekan (detik)
}

// DashboardSummary merepresentasikan data ringkasan untuk halaman Dashboard utama.
type DashboardSummary struct {
	TotalDevices   int     `json:"total_devices"`   // Jumlah total seluruh perangkat
	OnlineDevices  int     `json:"online_devices"`  // Jumlah perangkat yang aktif/normal
	OfflineDevices int     `json:"offline_devices"` // Jumlah perangkat yang mati/terputus
	WarningDevices int     `json:"warning_devices"` // Jumlah perangkat bermasalah/lambat
	LatestAlerts   []Alert `json:"latest_alerts"`   // Daftar peringatan terbaru
}
