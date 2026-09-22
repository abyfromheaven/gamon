package handler

// ============================================================================
// MODUL HANDLER PERINGATAN / ALERTS (handler/alert.go)
// ============================================================================
// Modul ini mengelola rute API untuk data Peringatan (Alerts) ketika perangkat
// jaringan mengalami masalah (Offline/Down).
// Fitur utama:
// 1. Menampilkan seluruh daftar peringatan dengan filter status / jenis perangkat.
// 2. Mengonfirmasi peringatan (Acknowledge) oleh admin.
// 3. Memulihkan peringatan (Resolve) setelah masalah teratasi.
// 4. Menghitung jumlah peringatan yang sedang berlangsung (ongoing).
// ============================================================================

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// AlertHandler menyimpan dependency ke koneksi database SQLite.
type AlertHandler struct {
	db *sql.DB
}

// NewAlertHandler membuat instansi baru AlertHandler.
func NewAlertHandler(db *sql.DB) *AlertHandler {
	return &AlertHandler{db: db}
}

// AlertResponse format standar balasan JSON untuk data peringatan ke frontend web.
type AlertResponse struct {
	ID             int     `json:"id"`
	DeviceID       int     `json:"device_id"`
	DeviceName     string  `json:"device_name"`
	DeviceType     string  `json:"device_type"`
	DeviceIP       string  `json:"device_ip"`
	Method         string  `json:"method,omitempty"`
	Title          string  `json:"title"`
	Status         string  `json:"status"`
	StartedAt      string  `json:"started_at"`
	ResolvedAt     *string `json:"resolved_at"`
	Description    string  `json:"description"`
	Acknowledged   bool    `json:"acknowledged"`
	AcknowledgedAt *string `json:"acknowledged_at"`
}

// HandleAlerts memproses rute GET /api/alerts (Mendapatkan daftar peringatan).
func (h *AlertHandler) HandleAlerts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAlerts(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// HandleAlert memproses rute /api/alerts/{id} serta aksi /resolve, /acknowledge, dan /count.
func (h *AlertHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	idTeks := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
	idTeks = strings.Split(idTeks, "/")[0]

	// Jika memanggil /api/alerts/count -> Kembalikan jumlah peringatan aktif
	if idTeks == "count" {
		h.alertCount(w, r)
		return
	}

	idPeringatan, err := strconv.Atoi(idTeks)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID Peringatan tidak valid")
		return
	}

	bagianSubURL := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/alerts/"+idTeks), "/")

	if len(bagianSubURL) > 1 {
		switch bagianSubURL[1] {
		case "resolve":
			h.resolveAlert(w, r, idPeringatan)
			return
		case "acknowledge":
			h.acknowledgeAlert(w, r, idPeringatan)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.getAlert(w, r, idPeringatan)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// listAlerts mengambil daftar seluruh peringatan dari tabel alerts digabung (JOIN) dengan tabel devices.
func (h *AlertHandler) listAlerts(w http.ResponseWriter, r *http.Request) {
	perintahSQL := `SELECT a.id, a.device_id, d.name, d.type, d.ip, d.method,
		a.title, a.status,
		a.started_at, a.resolved_at, a.description,
		a.acknowledged, a.acknowledged_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE 1=1`
	parameterQuery := []interface{}{}

	// Filter berdasarkan status jika ada di URL Query Parameters (?status=ongoing)
	if statusFilter := r.URL.Query().Get("status"); statusFilter != "" {
		perintahSQL += " AND a.status = ?"
		parameterQuery = append(parameterQuery, statusFilter)
	}

	// Filter berdasarkan tipe perangkat (?device_type=Server)
	if tipePerangkat := r.URL.Query().Get("device_type"); tipePerangkat != "" {
		perintahSQL += " AND d.type = ?"
		parameterQuery = append(parameterQuery, tipePerangkat)
	}

	perintahSQL += " ORDER BY a.started_at DESC"

	barisHasil, err := h.db.Query(perintahSQL, parameterQuery...)
	if err != nil {
		log.Printf("Gagal membaca daftar peringatan: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil daftar peringatan")
		return
	}
	defer barisHasil.Close()

	var daftarPeringatan []AlertResponse
	for barisHasil.Next() {
		var p AlertResponse
		var waktuMulai string
		var waktuSelesai, waktuKonfirmasi *string

		if err := barisHasil.Scan(&p.ID, &p.DeviceID, &p.DeviceName, &p.DeviceType, &p.DeviceIP, &p.Method,
			&p.Title, &p.Status,
			&waktuMulai, &waktuSelesai, &p.Description,
			&p.Acknowledged, &waktuKonfirmasi); err != nil {
			log.Printf("Gagal membaca baris peringatan: %v", err)
			continue
		}
		p.StartedAt = waktuMulai
		p.ResolvedAt = waktuSelesai
		p.AcknowledgedAt = waktuKonfirmasi
		daftarPeringatan = append(daftarPeringatan, p)
	}

	if daftarPeringatan == nil {
		daftarPeringatan = []AlertResponse{}
	}
	respondData(w, daftarPeringatan)
}

// getAlert mengambil 1 rincian data peringatan berdasarkan ID.
func (h *AlertHandler) getAlert(w http.ResponseWriter, _ *http.Request, idPeringatan int) {
	var p AlertResponse
	var waktuMulai string
	var waktuSelesai, waktuKonfirmasi *string

	err := h.db.QueryRow(`SELECT a.id, a.device_id, d.name, d.type, d.ip, a.title, a.status,
		a.started_at, a.resolved_at, a.description, a.acknowledged, a.acknowledged_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE a.id = ?`, idPeringatan).
		Scan(&p.ID, &p.DeviceID, &p.DeviceName, &p.DeviceType, &p.DeviceIP, &p.Title, &p.Status,
			&waktuMulai, &waktuSelesai, &p.Description, &p.Acknowledged, &waktuKonfirmasi)

	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Peringatan tidak ditemukan")
		return
	}
	if err != nil {
		log.Printf("Gagal mengambil data peringatan ID %d: %v", idPeringatan, err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data peringatan")
		return
	}

	p.StartedAt = waktuMulai
	p.ResolvedAt = waktuSelesai
	p.AcknowledgedAt = waktuKonfirmasi
	respondData(w, p)
}

// resolveAlert menandai bahwa gangguan pada perangkat telah teratasi (status -> 'resolved').
func (h *AlertHandler) resolveAlert(w http.ResponseWriter, _ *http.Request, idPeringatan int) {
	hasilExec, err := h.db.Exec("UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP WHERE id = ? AND status = 'ongoing'", idPeringatan)
	if err != nil {
		log.Printf("Gagal menyelesaikan peringatan ID %d: %v", idPeringatan, err)
		respondError(w, http.StatusInternalServerError, "Gagal memperbarui peringatan")
		return
	}

	jumlahBaris, _ := hasilExec.RowsAffected()
	if jumlahBaris == 0 {
		respondError(w, http.StatusNotFound, "Peringatan tidak ditemukan atau sudah teratasi")
		return
	}

	log.Printf("Peringatan ID %d berhasil dipulihkan (resolved)", idPeringatan)
	respondSuccess(w, "Peringatan berhasil dipulihkan")
}

// acknowledgeAlert menandai bahwa admin sudah menyadari/membaca peringatan tersebut (acknowledged = true).
func (h *AlertHandler) acknowledgeAlert(w http.ResponseWriter, _ *http.Request, idPeringatan int) {
	hasilExec, err := h.db.Exec("UPDATE alerts SET acknowledged = TRUE, acknowledged_at = CURRENT_TIMESTAMP WHERE id = ? AND acknowledged = FALSE", idPeringatan)
	if err != nil {
		log.Printf("Gagal mengonfirmasi peringatan ID %d: %v", idPeringatan, err)
		respondError(w, http.StatusInternalServerError, "Gagal mengonfirmasi peringatan")
		return
	}

	jumlahBaris, _ := hasilExec.RowsAffected()
	if jumlahBaris == 0 {
		respondError(w, http.StatusNotFound, "Peringatan tidak ditemukan atau sudah dikonfirmasi")
		return
	}

	log.Printf("Peringatan ID %d dikonfirmasi admin (acknowledged)", idPeringatan)
	respondSuccess(w, "Peringatan berhasil dikonfirmasi")
}

// alertCount menghitung berapa jumlah total peringatan yang masih aktif (ongoing).
func (h *AlertHandler) alertCount(w http.ResponseWriter, _ *http.Request) {
	var jumlahAktif int
	err := h.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE status = 'ongoing'").Scan(&jumlahAktif)
	if err != nil {
		log.Printf("Gagal menghitung jumlah peringatan aktif: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal menghitung peringatan")
		return
	}

	respondData(w, map[string]int{"ongoing": jumlahAktif})
}
