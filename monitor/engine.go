package monitor

// ============================================================
// MODULE ENGINE - JANTUNG SISTEM PEMANTAUAN
// ============================================================
// Module ini adalah komponen utama yang menjalankan semua proses pemantauan.
// Tugas utama Engine:
// 1. Menjalankan pengecekan perangkat secara berkala (misal setiap 3 detik)
// 2. Mengirim perintah ping ke setiap perangkat
// 3. Menyimpan hasil pengecekan ke database
// 4. Mendeteksi jika perangkat mati (offline) dan membuat alert
// 5. Mengirim notifikasi ke Telegram jika ada perubahan status
// 6. Mengirim data real-time ke browser via WebSocket
// ============================================================

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

// ============================================================
// ANTARMUKA (INTERFACE)
// ============================================================

// HubInterface adalah antarmuka untuk mengirim data ke browser.
// Setiap kali ada perubahan status, data akan dikirim ke browser
// agar tampilan di website bisa update secara real-time.
type HubInterface interface {
	Broadcast(msgType string, data interface{})
}

// Notifier adalah antarmuka untuk mengirim notifikasi ke Telegram.
// Ketika perangkat offline atau online, pesan akan dikirim ke Telegram
// agar admin bisa tahu dari ponselnya.
type Notifier interface {
	IsEnabled() bool
	SendAlert(deviceName, deviceIP string)
	SendRecovery(deviceName, deviceIP string)
}

// ============================================================
// STRUKTUR DATA (DATA STRUCTURES)
// ============================================================

// DeviceConfig adalah data konfigurasi untuk satu perangkat yang dipantau.
// Berisi informasi seperti IP address, metode pengecekan, dan interval.
type DeviceConfig struct {
	DeviceID int    // ID unik perangkat di database
	IP       string // Alamat IP perangkat (contoh: 192.168.1.1)
	URL      string // URL website jika metode HTTP
	Port     int    // Nomor port jika metode TCP
	Method   string // Metode pengecekan: "ICMP Ping", "TCP Port", atau "HTTP Check"
	Interval int    // Berapa detik sekali perangkat dicek
}

// DeviceStatus adalah data status terkini dari satu perangkat.
// Data ini dikirim ke browser untuk ditampilkan di halaman monitoring.
type DeviceStatus struct {
	DeviceID   int     `json:"device_id"`   // ID perangkat
	Name       string  `json:"name"`        // Nama perangkat
	Type       string  `json:"type"`        // Jenis (Server, Router, dll)
	IP         string  `json:"ip"`          // Alamat IP
	Status     string  `json:"status"`      // Status: "online" atau "offline"
	LatencyMs  float64 `json:"latency_ms"`  // Waktu respons dalam milidetik
	LastCheck  string  `json:"last_check"`  // Kapan terakhir kali dicek
	LastOnline string  `json:"last_online"` // Kapan terakhir kali perangkat online
}

// StatusChange adalah data yang menandakan perubahan status perangkat.
// Contoh: dari "online" menjadi "offline" atau sebaliknya.
type StatusChange struct {
	DeviceID   int    `json:"device_id"`   // ID perangkat
	DeviceName string `json:"device_name"` // Nama perangkat
	OldStatus  string `json:"old_status"`  // Status sebelumnya
	NewStatus  string `json:"new_status"`  // Status yang baru
	Timestamp  string `json:"timestamp"`   // Kapan perubahan terjadi
}

// ============================================================
// TIPE FUNGSI DAN OPSI KONFIGURASI
// ============================================================

// CheckFunc adalah tipe fungsi untuk mengecek status perangkat.
// Fungsi ini akan dipanggil oleh Engine setiap interval waktu.
type CheckFunc func(DeviceConfig, int) CheckResult

// EngineOption adalah fungsi untuk mengatur konfigurasi Engine.
// Pola ini disebut "Functional Options" agar kode lebih rapi.
type EngineOption func(*Engine)

// WithCheckFunc mengganti fungsi pengecekan bawaan.
// Berguna untuk pengujian (testing) agar tidak perlu melakukan ping beneran.
func WithCheckFunc(check CheckFunc) EngineOption {
	return func(engine *Engine) {
		if check != nil {
			engine.check = check
		}
	}
}

// WithNotifier menambahkan pengirim notifikasi Telegram ke Engine.
func WithNotifier(n Notifier) EngineOption {
	return func(engine *Engine) {
		engine.notifier = n
	}
}

// ============================================================
// STRUKTUR ENGINE (MESIN PEMANTAUAN)
// ============================================================

// Engine adalah struktur utama yang menjalankan semua proses pemantauan.
// Engine menyimpan semua data tentang perangkat yang sedang dipantau
// dan mengelola proses pengecekan secara berkala.
type Engine struct {
	hub      HubInterface // Untuk mengirim data ke browser
	db       *sql.DB      // Koneksi ke database
	notifier Notifier     // Untuk mengirim notifikasi Telegram

	mu         sync.Mutex                 // Pengaman agar data tidak rusak saat diakses bersamaan
	targets    map[int]context.CancelFunc // Menyimpan "tombol berhenti" untuk setiap perangkat
	lastStatus map[int]string             // Menyimpan status terakhir setiap perangkat
	failures   map[int]int                // Menyimpan jumlah kegagalan berturut-turut
	check      CheckFunc                  // Fungsi untuk mengecek perangkat
}

// ============================================================
// FUNGSI-FUNGSI UTAMA
// ============================================================

// NewEngine membuat Engine baru dan siap digunakan.
// Fungsi ini seperti "menyalakan mesin" pemantauan.
func NewEngine(hub HubInterface, db *sql.DB, options ...EngineOption) *Engine {
	engine := &Engine{
		hub:        hub,
		db:         db,
		targets:    make(map[int]context.CancelFunc),
		lastStatus: make(map[int]string),
		failures:   make(map[int]int),
		check: func(config DeviceConfig, seq int) CheckResult {
			// Fungsi bawaan: panggil PingOnce untuk mengecek perangkat
			result := PingOnce(config.IP, seq)
			result.DeviceID = config.DeviceID
			result.Method = config.Method
			return result
		},
	}
	// Terapkan semua opsi konfigurasi yang diberikan
	for _, option := range options {
		option(engine)
	}
	return engine
}

// Start memulai pemantauan untuk satu perangkat.
// Pemantauan akan berjalan terus menerus di latar belakang
// sampai fungsi Stop() dipanggil.
func (e *Engine) Start(config DeviceConfig) {
	// Tentukan interval pengecekan (default: 3 detik)
	interval := config.Interval
	if interval <= 0 {
		interval = 3
	}

	// Cek apakah perangkat sudah sedang dipantau
	e.mu.Lock()
	if _, exists := e.targets[config.DeviceID]; exists {
		e.mu.Unlock()
		return // Jika sudah dipantau, jangan mulai lagi
	}

	// Buat "tombol berhenti" untuk perangkat ini
	ctx, cancel := context.WithCancel(context.Background())
	e.targets[config.DeviceID] = cancel
	e.mu.Unlock()

	// Mulai pengecekan berkala di latar belakang (goroutine)
	go e.checkLoop(ctx, config, time.Duration(interval)*time.Second)
	log.Printf("Memulai pemantauan perangkat ID %d (%s)", config.DeviceID, config.IP)
}

// Stop menghentikan pemantauan untuk satu perangkat.
// Setelah dipanggil, perangkat tidak akan dicek lagi.
func (e *Engine) Stop(deviceID int) {
	e.mu.Lock()
	if cancel, exists := e.targets[deviceID]; exists {
		cancel() // Beritahu goroutine untuk berhenti
		delete(e.targets, deviceID)
	}
	delete(e.lastStatus, deviceID)
	delete(e.failures, deviceID)
	e.mu.Unlock()
}

// IsMonitoring mengecek apakah perangkat sedang dipantau.
// Mengembalikan true jika perangkat sedang aktif dipantau.
func (e *Engine) IsMonitoring(deviceID int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.targets[deviceID]
	return exists
}

// ============================================================
// PROSES PENGECEKAN BERKALA
// ============================================================

// checkLoop adalah perulangan yang mengecek perangkat secara berkala.
// Fungsi ini berjalan terus menerus sampai context dibatalkan (Stop dipanggil).
func (e *Engine) checkLoop(ctx context.Context, config DeviceConfig, interval time.Duration) {
	// Lakukan pengecekan pertama kali langsung
	seq := 1
	e.runCheck(config, seq)

	// Buat timer untuk pengecekan selanjutnya
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Perintah berhenti diterima, hentikan perulangan
			return
		case <-ticker.C:
			// Waktu pengecekan tiba, lakukan pengecekan
			seq++
			e.runCheck(config, seq)
		}
	}
}

// runCheck adalah fungsi utama yang menjalankan satu kali pengecekan.
// Fungsi ini melakukan:
// 1. Mengirim ping ke perangkat
// 2. Menyimpan hasil ke database
// 3. Memeriksa apakah status berubah
// 4. Mengirim notifikasi jika diperlukan
// 5. Mengirim data ke browser via WebSocket
func (e *Engine) runCheck(config DeviceConfig, seq int) {
	// 1. Lakukan pengecekan ping
	result := e.check(config, seq)
	result.DeviceID = config.DeviceID
	result.IP = config.IP
	result.Method = config.Method

	// Isi timestamp jika belum ada
	if result.Timestamp == "" {
		result.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if result.Details == nil {
		result.Details = map[string]any{}
	}

	// Tambahkan informasi last_online ke data yang dikirim ke browser
	lastOnline := e.getLastOnline(config.DeviceID)
	if lastOnline != "" {
		result.Details["last_online"] = lastOnline
	}

	// 2. Simpan hasil pengecekan ke database
	e.saveCheckResult(result)

	// 3. Periksa apakah status perangkat berubah
	change, shouldResolve, shouldAlert := e.trackStatus(result)

	// 4. Jika perangkat kembali online setelah sebelumnya offline
	if shouldResolve {
		// Tandai alert sebagai "teratasi" di database
		e.resolveAlerts(result.DeviceID)
		// Kirim notifikasi pemulihan ke Telegram
		if e.notifier != nil && e.notifier.IsEnabled() {
			deviceName := e.getDeviceName(result.DeviceID)
			e.notifier.SendRecovery(deviceName, result.IP)
		}
	}

	// 5. Jika perangkat offline melebihi batas toleransi
	if shouldAlert {
		// Buat alert baru di database
		e.createAlert(result.DeviceID)
		// Kirim notifikasi gangguan ke Telegram
		if e.notifier != nil && e.notifier.IsEnabled() {
			deviceName := e.getDeviceName(result.DeviceID)
			e.notifier.SendAlert(deviceName, result.IP)
		}
	}

	// 6. Kirim data perubahan status ke browser via WebSocket
	if change != nil {
		e.hub.Broadcast("status_change", *change)
	}
	// Kirim hasil pengecekan ke browser (untuk update tampilan)
	e.hub.Broadcast("check_result", result)
}

// ============================================================
// PENYIMPANAN DATA
// ============================================================

// saveCheckResult menyimpan hasil pengecekan ke database.
// Data ini berguna untuk melihat riwayat kapan perangkat online/offline.
func (e *Engine) saveCheckResult(result CheckResult) {
	details, err := json.Marshal(result.Details)
	if err != nil {
		details = []byte("{}")
	}
	_, err = e.db.Exec(`INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		result.DeviceID, result.Status, result.LatencyMs, result.TTL, result.Seq, string(details))
	if err != nil {
		log.Printf("Gagal menyimpan riwayat ping untuk perangkat %d: %v", result.DeviceID, err)
	}
}

// ============================================================
// PELACAKAN STATUS PERANGKAT
// ============================================================

// trackStatus memeriksa apakah status perangkat berubah.
// Fungsi ini menghitung kegagalan berturut-turut dan menentukan:
// - Kapan perangkat harus dianggap "offline" (setelah gagal beberapa kali)
// - Kapan perangkat dianggap "online" kembali
//
// Mengembalikan 3 nilai:
// - change: data perubahan status (jika ada)
// - shouldResolve: true jika alert harus ditandai "teratasi"
// - shouldAlert: true jika harus membuat alert baru
func (e *Engine) trackStatus(result CheckResult) (*StatusChange, bool, bool) {
	e.mu.Lock()
	oldStatus := e.lastStatus[result.DeviceID]
	newStatus := oldStatus
	resolveAlert := false
	createAlert := false

	// Ambil batas toleransi kegagalan dari pengaturan
	threshold := e.getFailureThreshold()

	switch result.Status {
	case StatusOffline:
		// Perangkat tidak merespon, tambahkan hitungan kegagalan
		e.failures[result.DeviceID]++

		// Jika sudah melebihi batas toleransi dan sebelumnya online,
		// maka status berubah menjadi offline
		if e.failures[result.DeviceID] >= threshold && oldStatus != StatusOffline {
			newStatus = StatusOffline
			createAlert = true
		}

	case StatusOnline:
		// Perangkat merespon, reset hitungan kegagalan
		e.failures[result.DeviceID] = 0
		newStatus = StatusOnline

		// Simpan waktu terakhir kali perangkat online di database
		e.updateLastOnline(result.DeviceID)

		// Jika sebelumnya offline, berarti perangkat sudah pulih
		if oldStatus == StatusOffline {
			resolveAlert = true
		}

	default:
		// Status tidak dikenali, abaikan
		e.mu.Unlock()
		return nil, false, false
	}

	// Jika status tidak berubah, tidak perlu melakukan apa-apa
	if newStatus == oldStatus {
		e.mu.Unlock()
		return nil, false, false
	}

	// Buat data perubahan status untuk dikirim ke browser
	change := &StatusChange{
		DeviceID:   result.DeviceID,
		DeviceName: e.getDeviceName(result.DeviceID),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Timestamp:  result.Timestamp,
	}
	e.lastStatus[result.DeviceID] = newStatus
	e.mu.Unlock()

	// Jika ini pengecekan pertama kali (belum ada status sebelumnya),
	// jangan kirim notifikasi perubahan
	if oldStatus == "" {
		return nil, resolveAlert, createAlert
	}

	return change, resolveAlert, createAlert
}

// ============================================================
// PENGELOLAAN ALERT (PERINGATAN)
// ============================================================

// createAlert membuat peringatan baru di database ketika perangkat offline.
// Alert ini akan muncul di dashboard untuk memberitahu admin.
func (e *Engine) createAlert(deviceID int) {
	_, err := e.db.Exec(`INSERT INTO alerts (device_id, title, status, severity, description)
		VALUES (?, 'Perangkat Tidak Merespon (Offline)', 'ongoing', 'critical', 'Perangkat tidak membalas ICMP Ping')`, deviceID)
	if err != nil {
		log.Printf("Gagal membuat catatan alert untuk perangkat %d: %v", deviceID, err)
	}
}

// resolveAlerts menandai alert sebagai "teratasi" di database.
// Dipanggil ketika perangkat kembali online setelah sebelumnya offline.
func (e *Engine) resolveAlerts(deviceID int) {
	_, err := e.db.Exec(`UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
		WHERE device_id = ? AND status = 'ongoing'`, deviceID)
	if err != nil {
		log.Printf("Gagal memperbarui status alert teratasi untuk perangkat %d: %v", deviceID, err)
	}
}

// ============================================================
// FUNGSI PEMBANTU (HELPER FUNCTIONS)
// ============================================================

// getDeviceName mengambil nama perangkat dari database.
// Jika perangkat tidak ditemukan, mengembalikan "Unknown".
func (e *Engine) getDeviceName(deviceID int) string {
	var name string
	if err := e.db.QueryRow("SELECT name FROM devices WHERE id = ?", deviceID).Scan(&name); err != nil {
		return "Unknown"
	}
	return name
}

// getFailureThreshold mengambil batas toleransi kegagalan dari pengaturan.
// Default: 3 kali gagal berturut-turut baru dianggap offline.
// Admin bisa mengubah nilai ini di halaman pengaturan.
func (e *Engine) getFailureThreshold() int {
	val := database.GetSetting(e.db, "failure_threshold", "3")
	threshold, err := strconv.Atoi(val)
	if err != nil || threshold < 1 {
		return 3
	}
	return threshold
}

// ============================================================
// PENGELOLAAN LAST ONLINE (WAKTU TERAKHIR ONLINE)
// ============================================================

// updateLastOnline mencatat waktu terakhir kali perangkat online di database.
// Fitur ini berguna untuk mengetahui kapan terakhir kali perangkat berfungsi normal.
func (e *Engine) updateLastOnline(deviceID int) {
	_, err := e.db.Exec("UPDATE devices SET last_online = CURRENT_TIMESTAMP WHERE id = ?", deviceID)
	if err != nil {
		log.Printf("Gagal memperbarui last_online untuk perangkat %d: %v", deviceID, err)
	}
}

// getLastOnline mengambil waktu terakhir kali perangkat online dari database.
// Mengembalikan string kosong jika belum ada data.
func (e *Engine) getLastOnline(deviceID int) string {
	var lastOnline sql.NullString
	err := e.db.QueryRow("SELECT last_online FROM devices WHERE id = ?", deviceID).Scan(&lastOnline)
	if err != nil || !lastOnline.Valid {
		return ""
	}
	return lastOnline.String
}
