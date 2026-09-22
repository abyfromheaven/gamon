package handler

// ============================================================================
// MODUL HANDLER DASHBOARD (handler/dashboard.go)
// ============================================================================
// Modul ini mengelola data statistik utama untuk halaman Dashboard GAMON,
// seperti:
// 1. Total seluruh perangkat yang terdaftar di database.
// 2. Jumlah perangkat berstatus Online (berhasil di-ping).
// 3. Jumlah perangkat berstatus Offline (terputus/gagal di-ping).
// 4. Daftar 5 peringatan terbaru (latest alerts).
// ============================================================================

import (
	"database/sql"
	"log"
	"net/http"
)

// DashboardHandler menyimpan koneksi database SQLite.
type DashboardHandler struct {
	db *sql.DB
}

// NewDashboardHandler membuat instansi baru DashboardHandler.
func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// RingkasanStatistik data jumlah total, online, dan offline perangkat.
type RingkasanStatistik struct {
	TotalDevices   int `json:"total_devices"`
	OnlineDevices  int `json:"online_devices"`
	OfflineDevices int `json:"offline_devices"`
}

// PeringatanTerbaruDashboard format data 5 alert terbaru di dashboard.
type PeringatanTerbaruDashboard struct {
	ID         int    `json:"id"`
	DeviceName string `json:"device_name"`
	Title      string `json:"title"`
	Status     string `json:"status"`
	StartedAt  string `json:"started_at"`
}

// DataDashboardResponse balasan JSON gabungan ringkasan statistik dan alert terbaru.
type DataDashboardResponse struct {
	Summary      RingkasanStatistik           `json:"summary"`
	LatestAlerts []PeringatanTerbaruDashboard `json:"latest_alerts"`
}

// HandleDashboard memproses rute GET /api/dashboard untuk menyajikan ringkasan informasi sistem.
func (h *DashboardHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var ringkasan RingkasanStatistik

	// 1. Hitung total seluruh perangkat
	err := h.db.QueryRow("SELECT COUNT(*) FROM devices").Scan(&ringkasan.TotalDevices)
	if err != nil {
		log.Printf("Gagal menghitung total perangkat: %v", err)
	}

	// 2. Hitung jumlah perangkat berstatus 'online' pada pengecekan terakhir
	err = h.db.QueryRow(`SELECT COUNT(DISTINCT d.id)
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		WHERE ph.status = 'online'`).Scan(&ringkasan.OnlineDevices)
	if err != nil {
		log.Printf("Gagal menghitung perangkat online: %v", err)
	}

	// 3. Hitung jumlah perangkat berstatus 'offline' pada pengecekan terakhir
	err = h.db.QueryRow(`SELECT COUNT(DISTINCT d.id)
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		WHERE ph.status = 'offline'`).Scan(&ringkasan.OfflineDevices)
	if err != nil {
		log.Printf("Gagal menghitung perangkat offline: %v", err)
	}

	// 4. Ambil 5 peringatan terbaru
	barisHasil, err := h.db.Query(`SELECT a.id, d.name, a.title, a.status, a.started_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		ORDER BY a.started_at DESC
		LIMIT 5`)
	if err != nil {
		log.Printf("Gagal membaca peringatan terbaru: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal memuat data dashboard")
		return
	}
	defer barisHasil.Close()

	var daftarPeringatanTerbaru []PeringatanTerbaruDashboard
	for barisHasil.Next() {
		var p PeringatanTerbaruDashboard
		var waktuMulai string
		if err := barisHasil.Scan(&p.ID, &p.DeviceName, &p.Title, &p.Status, &waktuMulai); err != nil {
			log.Printf("Gagal membaca baris peringatan dashboard: %v", err)
			continue
		}
		p.StartedAt = waktuMulai
		daftarPeringatanTerbaru = append(daftarPeringatanTerbaru, p)
	}

	if daftarPeringatanTerbaru == nil {
		daftarPeringatanTerbaru = []PeringatanTerbaruDashboard{}
	}

	dataBalasan := DataDashboardResponse{
		Summary:      ringkasan,
		LatestAlerts: daftarPeringatanTerbaru,
	}

	respondData(w, dataBalasan)
}
