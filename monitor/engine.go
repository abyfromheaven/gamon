package monitor

// Module Engine adalah jantung utama (core) pemantauan GAMON.
// Engine bertanggung jawab:
// 1. Menjalankan timer interval pemantauan tiap perangkat (menggunakan goroutine & ticker).
// 2. Memanggil fungsi PingOnce untuk mengecek status ICMP Ping.
// 3. Menyimpan hasil ping ke tabel database 'ping_history'.
// 4. Melacak kegagalan berturut-turut (misal 3 kali berturut-turut terputus -> memicu status Offline & Peringatan Alert).
// 5. Mengirimkan notifikasi ke Telegram Bot dan memancarkan data (broadcast) ke WebSocket.

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

// WithCheckFunc mengganti fungsi ping bawaan (berguna untuk testing / mock).
func WithCheckFunc(check CheckFunc) EngineOption {
	return func(engine *Engine) {
		if check != nil {
			engine.check = check
		}
	}
}

// WithNotifier menambahkan pengirim notifikasi (Telegram).
func WithNotifier(n Notifier) EngineOption {
	return func(engine *Engine) {
		engine.notifier = n
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
func NewEngine(hub HubInterface, db *sql.DB, options ...EngineOption) *Engine {
	engine := &Engine{
		hub:        hub,
		db:         db,
		targets:    make(map[int]context.CancelFunc),
		lastStatus: make(map[int]string),
		failures:   make(map[int]int),
		check: func(config DeviceConfig, seq int) CheckResult {
			result := PingOnce(config.IP, seq)
			result.DeviceID = config.DeviceID
			result.Method = config.Method
			return result
		},
	}
	for _, option := range options {
		option(engine)
	}
	return engine
}

// Start memulai proses pemantauan periodik untuk 1 perangkat secara asynchronous (goroutine).
func (e *Engine) Start(config DeviceConfig) {
	interval := config.Interval
	if interval <= 0 {
		interval = 3 // Default 3 detik sekali
	}

	e.mu.Lock()
	if _, exists := e.targets[config.DeviceID]; exists {
		e.mu.Unlock()
		return // Jika perangkat sudah dipantau, abaikan
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.targets[config.DeviceID] = cancel
	e.mu.Unlock()

	// Menjalankan perulangan pengecekan di goroutine terpisah
	go e.checkLoop(ctx, config, time.Duration(interval)*time.Second)
	log.Printf("Memulai pemantauan perangkat ID %d (%s)", config.DeviceID, config.IP)
}

// Stop menghentikan pemantauan untuk 1 perangkat.
func (e *Engine) Stop(deviceID int) {
	e.mu.Lock()
	if cancel, exists := e.targets[deviceID]; exists {
		cancel() // Hentikan goroutine pemantau
		delete(e.targets, deviceID)
	}
	delete(e.lastStatus, deviceID)
	delete(e.failures, deviceID)
	e.mu.Unlock()
}

// IsMonitoring mengecek apakah perangkat sedang aktif dipantau.
func (e *Engine) IsMonitoring(deviceID int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.targets[deviceID]
	return exists
}

// checkLoop adalah perulangan periodik (ticker) yang mengecek perangkat setiap interval waktu.
func (e *Engine) checkLoop(ctx context.Context, config DeviceConfig, interval time.Duration) {
	seq := 1
	e.runCheck(config, seq) // Pengecekan pertama

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done(): // Berhenti jika perintah Stop() dipanggil
			return
		case <-ticker.C: // Berjalan tiap detik interval tercapai
			seq++
			e.runCheck(config, seq)
		}
	}
}

// runCheck mengeksekusi ping, menyimpan riwayat ke DB, melacak status, serta mengirim sinyal ke WS & Telegram.
func (e *Engine) runCheck(config DeviceConfig, seq int) {
	result := e.check(config, seq)
	result.DeviceID = config.DeviceID
	result.IP = config.IP
	result.Method = config.Method
	if result.Timestamp == "" {
		result.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if result.Details == nil {
		result.Details = map[string]any{}
	}

	// 1. Simpan riwayat ping ke database
	e.saveCheckResult(result)

	// 2. Lacak perubahan status dan kebutuhan alert
	change, shouldResolve, shouldAlert := e.trackStatus(result)

	// 3. Jika perangkat kembali pulih (online)
	if shouldResolve {
		e.resolveAlerts(result.DeviceID)
		if e.notifier != nil && e.notifier.IsEnabled() {
			deviceName := e.getDeviceName(result.DeviceID)
			e.notifier.SendRecovery(deviceName, result.IP)
		}
	}

	// 4. Jika perangkat mati (offline) melampaui batas ambang (threshold)
	if shouldAlert {
		e.createAlert(result.DeviceID)
		if e.notifier != nil && e.notifier.IsEnabled() {
			deviceName := e.getDeviceName(result.DeviceID)
			e.notifier.SendAlert(deviceName, result.IP)
		}
	}

	// 5. Broadcast perubahan status ke frontend via WebSocket
	if change != nil {
		e.hub.Broadcast("status_change", *change)
	}
	e.hub.Broadcast("check_result", result)
}

// saveCheckResult menyimpan log hasil ping ke tabel ping_history di SQLite.
func (e *Engine) saveCheckResult(result CheckResult) {
	details, err := json.Marshal(result.Details)
	if err != nil {
		details = []byte("{}")
	}
	_, err = e.db.Exec(`INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`, result.DeviceID, result.Status, result.LatencyMs, result.TTL, result.Seq, string(details))
	if err != nil {
		log.Printf("Gagal menyimpan riwayat ping untuk perangkat %d: %v", result.DeviceID, err)
	}
}

// trackStatus mengevaluasi kegagalan berturut-turut untuk menentukan apakah status berubah.
func (e *Engine) trackStatus(result CheckResult) (*StatusChange, bool, bool) {
	e.mu.Lock()
	oldStatus := e.lastStatus[result.DeviceID]
	newStatus := oldStatus
	resolveAlert := false
	createAlert := false
	threshold := e.getFailureThreshold()

	switch result.Status {
	case StatusOffline:
		e.failures[result.DeviceID]++
		// Hanya picu alert jika gagal berturut-turut telah mencapai ambang batas (misal 3 kali)
		if e.failures[result.DeviceID] >= threshold && oldStatus != StatusOffline {
			newStatus = StatusOffline
			createAlert = true
		}
	case StatusOnline:
		e.failures[result.DeviceID] = 0
		newStatus = StatusOnline
		if oldStatus == StatusOffline {
			resolveAlert = true
		}
	default:
		e.mu.Unlock()
		return nil, false, false
	}

	if newStatus == oldStatus {
		e.mu.Unlock()
		return nil, false, false
	}

	change := &StatusChange{
		DeviceID:   result.DeviceID,
		DeviceName: e.getDeviceName(result.DeviceID),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Timestamp:  result.Timestamp,
	}
	e.lastStatus[result.DeviceID] = newStatus
	e.mu.Unlock()

	if oldStatus == "" {
		return nil, resolveAlert, createAlert
	}

	return change, resolveAlert, createAlert
}

// createAlert mencatat peristiwa perangkat mati ke tabel alerts.
func (e *Engine) createAlert(deviceID int) {
	_, err := e.db.Exec(`INSERT INTO alerts (device_id, title, status, severity, description)
		VALUES (?, 'Perangkat Tidak Merespon (Offline)', 'ongoing', 'critical', 'Perangkat tidak membalas ICMP Ping')`, deviceID)
	if err != nil {
		log.Printf("Gagal membuat catatan alert untuk perangkat %d: %v", deviceID, err)
	}
}

// resolveAlerts mengubah status peringatan bermasalah menjadi pulih (resolved).
func (e *Engine) resolveAlerts(deviceID int) {
	_, err := e.db.Exec(`UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
		WHERE device_id = ? AND status = 'ongoing'`, deviceID)
	if err != nil {
		log.Printf("Gagal memperbarui status alert teratasi untuk perangkat %d: %v", deviceID, err)
	}
}

// getDeviceName mengambil nama perangkat dari database.
func (e *Engine) getDeviceName(deviceID int) string {
	var name string
	if err := e.db.QueryRow("SELECT name FROM devices WHERE id = ?", deviceID).Scan(&name); err != nil {
		return "Unknown"
	}
	return name
}

// getFailureThreshold mengambil batas ambang batas kegagalan ping dari pengaturan database (default: 3 kali).
func (e *Engine) getFailureThreshold() int {
	val := database.GetSetting(e.db, "failure_threshold", "3")
	threshold, err := strconv.Atoi(val)
	if err != nil || threshold < 1 {
		return 3
	}
	return threshold
}
