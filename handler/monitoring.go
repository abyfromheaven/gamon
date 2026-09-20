package handler

// ============================================================
// MODULE MONITORING HANDLER
// ============================================================
// Module ini menangani permintaan data monitoring dari browser.
// Tugas utama:
// 1. Menampilkan daftar semua perangkat beserta statusnya
// 2. Menampilkan detail riwayat pengecekan satu perangkat
// 3. Menghitung berapa lama perangkat dalam status tertentu (uptime)
//
// Endpoint yang dilayani:
// - GET /api/monitoring          -> Daftar semua perangkat + status
// - GET /api/monitoring/:id/history -> Riwayat pengecekan satu perangkat
// ============================================================

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MonitoringHandler menangani semua permintaan terkait data monitoring.
// Struct ini memiliki akses ke database untuk mengambil data.
type MonitoringHandler struct {
	db *sql.DB
}

// NewMonitoringHandler membuat handler baru yang siap digunakan.
func NewMonitoringHandler(db *sql.DB) *MonitoringHandler {
	return &MonitoringHandler{db: db}
}

// HandleMonitoring menangani permintaan GET /api/monitoring.
// Mengembalikan daftar semua perangkat beserta status terkininya.
func (h *MonitoringHandler) HandleMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.listMonitoringStatus(w, r)
}

// HandleMonitoringDevice menangani permintaan untuk satu perangkat tertentu.
// Contoh: GET /api/monitoring/1/history (riwayat pengecekan perangkat ID 1)
func (h *MonitoringHandler) HandleMonitoringDevice(w http.ResponseWriter, r *http.Request) {
	// Ambil ID perangkat dari URL (contoh: /api/monitoring/1 -> "1")
	path := strings.TrimPrefix(r.URL.Path, "/api/monitoring/")
	parts := strings.Split(path, "/")

	idStr := parts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	// Jika URL-nya /api/monitoring/:id/history, tampilkan riwayat
	if len(parts) > 1 && parts[1] == "history" {
		h.getDeviceHistory(w, r, id)
		return
	}

	respondError(w, http.StatusNotFound, "Endpoint not found")
}

// ============================================================
// STRUKTUR DATA YANG DIKIRIM KE BROWSER
// ============================================================

// MonitoringStatus adalah format data yang dikirim ke browser.
// Berisi informasi lengkap tentang status satu perangkat.
type MonitoringStatus struct {
	DeviceID              int     `json:"device_id"`               // ID perangkat
	DeviceName            string  `json:"device_name"`             // Nama perangkat
	DeviceType            string  `json:"device_type"`             // Jenis perangkat
	IP                    string  `json:"ip"`                      // Alamat IP
	Method                string  `json:"method"`                  // Metode pengecekan
	Status                string  `json:"status"`                  // Status: online/offline/unknown
	LatencyMs             float64 `json:"latency_ms"`              // Waktu respons (milidetik)
	LastCheck             *string `json:"last_check"`              // Kapan terakhir kali dicek
	LastOnline            *string `json:"last_online"`             // Kapan terakhir kali online
	CurrentStatusSince    string  `json:"current_status_since"`    // Sejak kapan status saat ini aktif
	CurrentStatusDuration string  `json:"current_status_duration"` // Berapa lama status saat ini berlangsung
	Interval              int     `json:"interval"`                // Interval pengecekan (detik)
}

// ============================================================
// FUNGSI UTAMA: MENGAMBIL DAFTAR STATUS PERANGKAT
// ============================================================

// listMonitoringStatus mengambil data status semua perangkat dari database.
// Fungsi ini menggabungkan data dari tabel devices dan ping_history
// untuk menampilkan status terkini setiap perangkat.
func (h *MonitoringHandler) listMonitoringStatus(w http.ResponseWriter, _ *http.Request) {

	// Query SQL: ambil data perangkat dan gabungkan dengan hasil ping terakhir
	rows, err := h.db.Query(`SELECT d.id, d.name, d.type, d.ip, d.method, d.check_interval,
		COALESCE(ph.status, 'unknown') as status,
		COALESCE(ph.latency_ms, 0) as latency_ms,
		COALESCE(ph.timestamp, '') as last_check,
		d.last_online
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		ORDER BY d.name`)
	if err != nil {
		log.Printf("Error listing monitoring status: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list monitoring status")
		return
	}
	defer rows.Close()

	// Baca hasil query satu per satu
	var statuses []MonitoringStatus
	for rows.Next() {
		var s MonitoringStatus
		var lastCheck string
		var lastOnline sql.NullString

		// Baca data dari database
		if err := rows.Scan(&s.DeviceID, &s.DeviceName, &s.DeviceType, &s.IP, &s.Method,
			&s.Interval, &s.Status, &s.LatencyMs, &lastCheck, &lastOnline); err != nil {
			log.Printf("Error scanning monitoring status: %v", err)
			continue
		}

		// Isi data last_check jika ada
		if lastCheck != "" {
			s.LastCheck = &lastCheck
		}

		// Isi data last_online jika ada
		if lastOnline.Valid {
			s.LastOnline = &lastOnline.String
		}

		// Hitung sejak kapan status saat ini aktif dan berapa lamanya
		since, duration := h.calculateStatusDuration(s.DeviceID, s.Status)
		s.CurrentStatusSince = since
		s.CurrentStatusDuration = duration

		statuses = append(statuses, s)
	}

	// Jika tidak ada data, kirim array kosong (bukan null)
	if statuses == nil {
		statuses = []MonitoringStatus{}
	}
	respondData(w, statuses)
}

// ============================================================
// KALKULASI UPTIME / STATUS DURATION
// ============================================================

// calculateStatusDuration menghitung:
// 1. Sejak kapan perangkat dalam status saat ini (misal: "Online sejak 20:13")
// 2. Berapa lama status saat ini berlangsung (misal: "2j 15m 30d")
//
// Cara kerja:
// - Cari waktu paling baru di mana status BERBEDA dari status saat ini
// - Waktu itu menandakan awal dari periode status saat ini
// - Hitung selisih waktu dari saat itu hingga sekarang
func (h *MonitoringHandler) calculateStatusDuration(deviceID int, currentStatus string) (string, string) {
	// Jika status tidak diketahui, kembalikan kosong
	if currentStatus == "unknown" {
		return "", ""
	}

	// Cari timestamp terakhir di mana status BERBEDA dari status saat ini
	// Contoh: jika status saat ini "online", cari timestamp terakhir "offline"
	var sinceTime sql.NullString
	err := h.db.QueryRow(`
		SELECT timestamp FROM ping_history 
		WHERE device_id = ? AND status != ? 
		ORDER BY timestamp DESC LIMIT 1
	`, deviceID, currentStatus).Scan(&sinceTime)

	if err != nil || !sinceTime.Valid {
		// Jika tidak ada riwayat perubahan, ambil data pertama kali
		err2 := h.db.QueryRow(`
			SELECT timestamp FROM ping_history 
			WHERE device_id = ? 
			ORDER BY timestamp ASC LIMIT 1
		`, deviceID).Scan(&sinceTime)
		if err2 != nil || !sinceTime.Valid {
			return "", ""
		}
	}

	// Ubah string timestamp menjadi objek waktu
	parsedTime, err := time.Parse("2006-01-02 15:04:05", sinceTime.String)
	if err != nil {
		// Coba format lain (RFC3339) jika format pertama gagal
		parsedTime, err = time.Parse(time.RFC3339, sinceTime.String)
		if err != nil {
			return sinceTime.String, ""
		}
	}

	// Hitung berapa lama status saat ini berlangsung
	now := time.Now()
	duration := now.Sub(parsedTime)

	// Format waktu mulai (contoh: "20:13:45")
	since := parsedTime.Format("15:04:05")

	// Format durasi (contoh: "2j 15m 30d")
	durationStr := formatDuration(duration)

	return since, durationStr
}

// formatDuration mengubah durasi waktu menjadi format yang mudah dibaca.
//
// Contoh hasil:
// - "5d"              (5 detik)
// - "3m 15d"          (3 menit 15 detik)
// - "1j 30m 45d"      (1 jam 30 menit 45 detik)
//
// Keterangan singkatan:
// - j = jam (hours)
// - m = menit (minutes)
// - d = detik (seconds)
func formatDuration(d time.Duration) string {
	// Jika durasi kurang dari 1 detik
	if d < time.Second {
		return "0d"
	}

	// Hitung jam, menit, dan detik
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	// Format sesuai kebutuhan
	if hours > 0 {
		return strconv.Itoa(hours) + "j " + strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "d"
	}
	if minutes > 0 {
		return strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "d"
	}
	return strconv.Itoa(seconds) + "d"
}

// ============================================================
// RIWAYAT PENGECEKAN PERANGKAT
// ============================================================

// getDeviceHistory mengambil riwayat pengecekan untuk satu perangkat.
// Mengembalikan 50 data terakhir yang disimpan di database.
// Data ini digunakan untuk menampilkan grafik latency di detail panel.
func (h *MonitoringHandler) getDeviceHistory(w http.ResponseWriter, _ *http.Request, deviceID int) {
	// Struktur data riwayat ping yang dikirim ke browser
	type PingRecord struct {
		ID        int     `json:"id"`         // ID unik record
		Status    string  `json:"status"`     // Status saat itu (online/offline)
		LatencyMs float64 `json:"latency_ms"` // Waktu respons saat itu
		TTL       int     `json:"ttl"`        // TTL paket ping
		Seq       int     `json:"seq"`        // Nomor urut pengecekan
		Details   string  `json:"details"`    // Detail tambahan (JSON)
		Timestamp string  `json:"timestamp"`  // Kapan pengecekan dilakukan
	}

	// Ambil 50 riwayat terakhir dari database
	rows, err := h.db.Query(`SELECT id, status, latency_ms, ttl, seq, details, timestamp
		FROM ping_history
		WHERE device_id = ?
		ORDER BY timestamp DESC
		LIMIT 50`, deviceID)
	if err != nil {
		log.Printf("Error getting device history: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get device history")
		return
	}
	defer rows.Close()

	// Baca hasil query
	var history []PingRecord
	for rows.Next() {
		var p PingRecord
		var timestamp string
		if err := rows.Scan(&p.ID, &p.Status, &p.LatencyMs, &p.TTL, &p.Seq, &p.Details, &timestamp); err != nil {
			log.Printf("Error scanning ping record: %v", err)
			continue
		}
		p.Timestamp = timestamp
		history = append(history, p)
	}

	// Kirim array kosong jika tidak ada data
	if history == nil {
		history = []PingRecord{}
	}
	respondData(w, history)
}
