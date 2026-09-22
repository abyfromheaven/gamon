package monitor

// ============================================================================
// MODUL ENGINE MONITORING (monitor/engine.go)
// ============================================================================
// Module Engine adalah jantung utama (core) pemantauan GAMON.
// Engine bertanggung jawab:
// 1. Menjalankan timer interval pemantauan tiap perangkat (goroutine & ticker).
// 2. Memanggil fungsi PingOnce untuk mengecek status ICMP Ping.
// 3. Menyimpan hasil ping ke tabel database 'ping_history'.
// 4. Melacak kegagalan berturut-turut (misal 3 kali berturut-turut terputus -> memicu status Offline & Peringatan Alert).
// 5. Mengirimkan notifikasi ke Telegram Bot dan memancarkan data (broadcast) ke WebSocket.
// ============================================================================

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"gamon/database"
)

// HubInterface adalah antarmuka untuk menyebarkan pesan ke seluruh koneksi web browser (WebSocket).
type HubInterface interface {
	Broadcast(msgType string, data interface{})
}

// Notifier adalah antarmuka untuk pengiriman Notifikasi (misalnya ke Telegram Bot).
type Notifier interface {
	IsEnabled() bool
	SendAlert(deviceName, deviceIP string)
	SendRecovery(deviceName, deviceIP string)
}

// DeviceConfig adalah data konfigurasi pemantauan perangkat.
type DeviceConfig struct {
	DeviceID int    // ID Perangkat
	IP       string // Alamat IP
	URL      string // URL
	Port     int    // Port
	Method   string // Metode pemantauan (ICMP Ping)
	Interval int    // Interval pemantauan (detik)
}

// DeviceStatus adalah format data status perangkat terkini.
type DeviceStatus struct {
	DeviceID  int     `json:"device_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	IP        string  `json:"ip"`
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	LastCheck string  `json:"last_check"`
}

// StatusChange merepresentasikan perubahan status perangkat (misal dari online ke offline atau sebaliknya).
type StatusChange struct {
	DeviceID   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	OldStatus  string `json:"old_status"`
	NewStatus  string `json:"new_status"`
	Timestamp  string `json:"timestamp"`
}

// CheckFunc tipe fungsi pelaksana pengecekan ping.
type CheckFunc func(DeviceConfig, int) CheckResult

// EngineOption adalah opsi konfigurasi untuk pembuatan Engine baru.
type EngineOption func(*Engine)

// WithCheckFunc mengganti fungsi ping bawaan (berguna untuk pengujian / testing mock).
func WithCheckFunc(fungsiCek CheckFunc) EngineOption {
	return func(mesin *Engine) {
		if fungsiCek != nil {
			mesin.check = fungsiCek
		}
	}
}

// WithNotifier menambahkan pengirim notifikasi (Telegram).
func WithNotifier(pengirimNotif Notifier) EngineOption {
	return func(mesin *Engine) {
		mesin.notifier = pengirimNotif
	}
}

// Engine adalah struktur utama mesin pemonitoring.
type Engine struct {
	hub      HubInterface // WebSocket hub
	db       *sql.DB      // Koneksi database SQLite
	notifier Notifier     // Notifier Telegram

	mu         sync.Mutex                 // Mutex untuk keamanan akses data bersama (concurrency thread-safe)
	targets    map[int]context.CancelFunc // Map goroutine pemantau per ID perangkat
	lastStatus map[int]string             // Map status terakhir perangkat
	failures   map[int]int                // Map jumlah kegagalan berturut-turut
	check      CheckFunc                  // Fungsi eksekusi ping
}

// NewEngine membuat objek Engine pemonitoring baru.
func NewEngine(hub HubInterface, db *sql.DB, opsi ...EngineOption) *Engine {
	mesin := &Engine{
		hub:        hub,
		db:         db,
		targets:    make(map[int]context.CancelFunc),
		lastStatus: make(map[int]string),
		failures:   make(map[int]int),
		check: func(konfig DeviceConfig, urutan int) CheckResult {
			hasil := PingOnce(konfig.IP, urutan)
			hasil.DeviceID = konfig.DeviceID
			hasil.Method = konfig.Method
			return hasil
		},
	}
	for _, pilihan := range opsi {
		pilihan(mesin)
	}
	return mesin
}

// Start memulai proses pemantauan periodik untuk 1 perangkat secara asynchronous (goroutine).
func (e *Engine) Start(konfig DeviceConfig) {
	intervalDetik := konfig.Interval
	if intervalDetik <= 0 {
		intervalDetik = 3 // Default 3 detik sekali
	}

	e.mu.Lock()
	if _, sudahAktif := e.targets[konfig.DeviceID]; sudahAktif {
		e.mu.Unlock()
		return // Jika perangkat sudah dipantau, abaikan
	}
	konteks, batalkan := context.WithCancel(context.Background())
	e.targets[konfig.DeviceID] = batalkan
	e.mu.Unlock()

	// Menjalankan perulangan pengecekan di goroutine terpisah
	go e.checkLoop(konteks, konfig, time.Duration(intervalDetik)*time.Second)
	log.Printf("Memulai pemantauan perangkat ID %d (%s)", konfig.DeviceID, konfig.IP)
}

// Stop menghentikan pemantauan untuk 1 perangkat.
func (e *Engine) Stop(idPerangkat int) {
	e.mu.Lock()
	if batalkan, ada := e.targets[idPerangkat]; ada {
		batalkan() // Hentikan goroutine pemantau
		delete(e.targets, idPerangkat)
	}
	delete(e.lastStatus, idPerangkat)
	delete(e.failures, idPerangkat)
	e.mu.Unlock()
}

// IsMonitoring mengecek apakah perangkat sedang aktif dipantau.
func (e *Engine) IsMonitoring(idPerangkat int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ada := e.targets[idPerangkat]
	return ada
}

// checkLoop adalah perulangan periodik (ticker) yang mengecek perangkat setiap interval waktu.
func (e *Engine) checkLoop(konteks context.Context, konfig DeviceConfig, interval time.Duration) {
	urutan := 1
	e.runCheck(konfig, urutan) // Pengecekan pertama

	penghitungWaktu := time.NewTicker(interval)
	defer penghitungWaktu.Stop()
	for {
		select {
		case <-konteks.Done(): // Berhenti jika perintah Stop() dipanggil
			return
		case <-penghitungWaktu.C: // Berjalan tiap detik interval tercapai
			urutan++
			e.runCheck(konfig, urutan)
		}
	}
}

// runCheck mengeksekusi ping, menyimpan riwayat ke DB, melacak status, serta mengirim sinyal ke WS & Telegram.
func (e *Engine) runCheck(konfig DeviceConfig, urutan int) {
	hasil := e.check(konfig, urutan)
	hasil.DeviceID = konfig.DeviceID
	hasil.IP = konfig.IP
	hasil.Method = konfig.Method
	if hasil.Timestamp == "" {
		hasil.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if hasil.Details == nil {
		hasil.Details = map[string]any{}
	}

	// 1. Simpan riwayat ping ke database
	e.saveCheckResult(hasil)

	// 2. Lacak perubahan status dan kebutuhan alert
	perubahan, butuhPulih, butuhPeringatan := e.trackStatus(hasil)

	// 3. Jika perangkat kembali pulih (online)
	if butuhPulih {
		e.resolveAlerts(hasil.DeviceID)
		if e.notifier != nil && e.notifier.IsEnabled() {
			namaPerangkat := e.getDeviceName(hasil.DeviceID)
			e.notifier.SendRecovery(namaPerangkat, hasil.IP)
		}
	}

	// 4. Jika perangkat mati (offline) melampaui batas ambang (threshold)
	if butuhPeringatan {
		e.createAlert(hasil.DeviceID)
		if e.notifier != nil && e.notifier.IsEnabled() {
			namaPerangkat := e.getDeviceName(hasil.DeviceID)
			e.notifier.SendAlert(namaPerangkat, hasil.IP)
		}
	}

	// 5. Broadcast perubahan status ke frontend via WebSocket
	if perubahan != nil {
		e.hub.Broadcast("status_change", *perubahan)
	}
	e.hub.Broadcast("check_result", hasil)
}

// saveCheckResult menyimpan log hasil ping ke tabel ping_history di SQLite.
func (e *Engine) saveCheckResult(hasil CheckResult) {
	rincianTeks, err := json.Marshal(hasil.Details)
	if err != nil {
		rincianTeks = []byte("{}")
	}
	_, err = e.db.Exec(`INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`, hasil.DeviceID, hasil.Status, hasil.LatencyMs, hasil.TTL, hasil.Seq, string(rincianTeks))
	if err != nil {
		log.Printf("Gagal menyimpan riwayat ping untuk perangkat %d: %v", hasil.DeviceID, err)
	}
}

// trackStatus mengevaluasi kegagalan berturut-turut untuk menentukan apakah status berubah.
func (e *Engine) trackStatus(hasil CheckResult) (*StatusChange, bool, bool) {
	e.mu.Lock()
	statusLama := e.lastStatus[hasil.DeviceID]
	statusBaru := statusLama
	pulihkanAlert := false
	buatAlert := false
	ambangBatas := e.getFailureThreshold()

	switch hasil.Status {
	case StatusOffline:
		e.failures[hasil.DeviceID]++
		// Hanya picu alert jika gagal berturut-turut telah mencapai ambang batas (misal 3 kali)
		if e.failures[hasil.DeviceID] >= ambangBatas && statusLama != StatusOffline {
			statusBaru = StatusOffline
			buatAlert = true
		}
	case StatusOnline:
		e.failures[hasil.DeviceID] = 0
		statusBaru = StatusOnline
		if statusLama == StatusOffline {
			pulihkanAlert = true
		}
	default:
		e.mu.Unlock()
		return nil, false, false
	}

	if statusBaru == statusLama {
		e.mu.Unlock()
		return nil, false, false
	}

	perubahan := &StatusChange{
		DeviceID:   hasil.DeviceID,
		DeviceName: e.getDeviceName(hasil.DeviceID),
		OldStatus:  statusLama,
		NewStatus:  statusBaru,
		Timestamp:  hasil.Timestamp,
	}
	e.lastStatus[hasil.DeviceID] = statusBaru
	e.mu.Unlock()

	if statusLama == "" {
		return nil, pulihkanAlert, buatAlert
	}

	return perubahan, pulihkanAlert, buatAlert
}

// createAlert mencatat peristiwa perangkat mati ke tabel alerts.
func (e *Engine) createAlert(idPerangkat int) {
	_, err := e.db.Exec(`INSERT INTO alerts (device_id, title, status, description)
		VALUES (?, 'Perangkat Tidak Merespon (Offline)', 'ongoing', 'Perangkat tidak membalas ICMP Ping')`, idPerangkat)
	if err != nil {
		log.Printf("Gagal membuat catatan alert untuk perangkat %d: %v", idPerangkat, err)
	}
}

// resolveAlerts mengubah status peringatan bermasalah menjadi pulih (resolved).
func (e *Engine) resolveAlerts(idPerangkat int) {
	_, err := e.db.Exec(`UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
		WHERE device_id = ? AND status = 'ongoing'`, idPerangkat)
	if err != nil {
		log.Printf("Gagal memperbarui status alert teratasi untuk perangkat %d: %v", idPerangkat, err)
	}
}

// getDeviceName mengambil nama perangkat dari database.
func (e *Engine) getDeviceName(idPerangkat int) string {
	var nama string
	if err := e.db.QueryRow("SELECT name FROM devices WHERE id = ?", idPerangkat).Scan(&nama); err != nil {
		return "Unknown"
	}
	return nama
}

// getFailureThreshold mengambil batas ambang batas kegagalan ping dari pengaturan database (default: 3 kali).
func (e *Engine) getFailureThreshold() int {
	nilaiTeks := database.GetSetting(e.db, "failure_threshold", "3")
	ambang, err := strconv.Atoi(nilaiTeks)
	if err != nil || ambang < 1 {
		return 3
	}
	return ambang
}
