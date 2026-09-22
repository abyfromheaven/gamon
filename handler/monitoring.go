package handler

// ============================================================================
// MODUL HANDLER PEMANTAUAN / MONITORING (handler/monitoring.go)
// ============================================================================
// Modul ini mengelola API endpoint untuk:
// 1. Menampilkan status pemantauan real-time seluruh perangkat.
// 2. Menampilkan riwayat (history) 50 data ping terakhir suatu perangkat untuk
//    ditampilkan dalam bentuk grafik latensi (chart) pada frontend web.
// ============================================================================

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// MonitoringHandler menyimpan koneksi database SQLite.
type MonitoringHandler struct {
	db *sql.DB
}

// NewMonitoringHandler membuat instansi baru MonitoringHandler.
func NewMonitoringHandler(db *sql.DB) *MonitoringHandler {
	return &MonitoringHandler{db: db}
}

// DataStatusPemantauan format data status pemantauan 1 perangkat.
type DataStatusPemantauan struct {
	DeviceID   int     `json:"device_id"`
	DeviceName string  `json:"device_name"`
	DeviceType string  `json:"device_type"`
	IP         string  `json:"ip"`
	Method     string  `json:"method"`
	Status     string  `json:"status"`
	LatencyMs  float64 `json:"latency_ms"`
	LastCheck  *string `json:"last_check"`
	Interval   int     `json:"interval"`
}

// CatatanRiwayatPing format data 1 baris riwayat ping.
type CatatanRiwayatPing struct {
	ID        int     `json:"id"`
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	TTL       int     `json:"ttl"`
	Seq       int     `json:"seq"`
	Details   string  `json:"details"`
	Timestamp string  `json:"timestamp"`
}

// HandleMonitoring memproses rute GET /api/monitoring (Status pemantauan seluruh perangkat).
func (h *MonitoringHandler) HandleMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	h.listMonitoringStatus(w, r)
}

// HandleMonitoringDevice memproses rute /api/monitoring/{id}/history (Riwayat ping perangkat).
func (h *MonitoringHandler) HandleMonitoringDevice(w http.ResponseWriter, r *http.Request) {
	jalur := strings.TrimPrefix(r.URL.Path, "/api/monitoring/")
	bagian := strings.Split(jalur, "/")

	idTeks := bagian[0]
	idPerangkat, err := strconv.Atoi(idTeks)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID Perangkat tidak valid")
		return
	}

	if len(bagian) > 1 && bagian[1] == "history" {
		h.getDeviceHistory(w, r, idPerangkat)
		return
	}

	respondError(w, http.StatusNotFound, "Endpoint tidak ditemukan")
}

// listMonitoringStatus mengambil status pemantauan terbaru seluruh perangkat dari SQLite.
func (h *MonitoringHandler) listMonitoringStatus(w http.ResponseWriter, _ *http.Request) {
	barisHasil, err := h.db.Query(`SELECT d.id, d.name, d.type, d.ip, d.method, d.check_interval,
		COALESCE(ph.status, 'unknown') as status,
		COALESCE(ph.latency_ms, 0) as latency_ms,
		COALESCE(ph.timestamp, '') as last_check
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		ORDER BY d.name`)
	if err != nil {
		log.Printf("Gagal membaca status pemantauan: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil status pemantauan")
		return
	}
	defer barisHasil.Close()

	var daftarStatus []DataStatusPemantauan
	for barisHasil.Next() {
		var s DataStatusPemantauan
		var waktuPengecekan string
		if err := barisHasil.Scan(&s.DeviceID, &s.DeviceName, &s.DeviceType, &s.IP, &s.Method, &s.Interval, &s.Status, &s.LatencyMs, &waktuPengecekan); err != nil {
			log.Printf("Gagal membaca baris status pemantauan: %v", err)
			continue
		}
		if waktuPengecekan != "" {
			s.LastCheck = &waktuPengecekan
		}
		daftarStatus = append(daftarStatus, s)
	}

	if daftarStatus == nil {
		daftarStatus = []DataStatusPemantauan{}
	}
	respondData(w, daftarStatus)
}

// getDeviceHistory mengambil 50 baris riwayat latensi ping terakhir dari tabel ping_history.
func (h *MonitoringHandler) getDeviceHistory(w http.ResponseWriter, _ *http.Request, idPerangkat int) {
	barisHasil, err := h.db.Query(`SELECT id, status, latency_ms, ttl, seq, details, timestamp
		FROM ping_history
		WHERE device_id = ?
		ORDER BY timestamp DESC
		LIMIT 50`, idPerangkat)
	if err != nil {
		log.Printf("Gagal membaca riwayat ping perangkat ID %d: %v", idPerangkat, err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil riwayat pemantauan")
		return
	}
	defer barisHasil.Close()

	var daftarRiwayat []CatatanRiwayatPing
	for barisHasil.Next() {
		var p CatatanRiwayatPing
		var waktu string
		if err := barisHasil.Scan(&p.ID, &p.Status, &p.LatencyMs, &p.TTL, &p.Seq, &p.Details, &waktu); err != nil {
			log.Printf("Gagal membaca baris riwayat ping: %v", err)
			continue
		}
		p.Timestamp = waktu
		daftarRiwayat = append(daftarRiwayat, p)
	}

	if daftarRiwayat == nil {
		daftarRiwayat = []CatatanRiwayatPing{}
	}
	respondData(w, daftarRiwayat)
}
