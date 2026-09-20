package main

// GAMON (Garda Monitoring) - Aplikasi Pemonitor Jaringan & Perangkat Berbasis Web.
// Berkas main.go merupakan titik masuk utama (entry point) aplikasi backend server Go.

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gamon/database"
	"gamon/handler"
	"gamon/monitor"
	"gamon/notification"
)

// main adalah fungsi utama yang dijalankan saat server backend GAMON dinyalakan.
// Urutan alurnya:
// 1. Membuka & menyiapkan Database SQLite.
// 2. Menjalankan WebSocket Hub (penyebar data real-time ke web frontend).
// 3. Menyiapkan Notifikasi Telegram Bot.
// 4. Menyiapkan Engine Monitoring (Pinger ICMP).
// 5. Mendaftarkan alamat rute API (REST API Endpoints).
// 6. Otomatis mulai memantau perangkat aktif dari database.
// 7. Menjalankan Server Web di port 8080.
func main() {
	// 1. Inisialisasi dan koneksi ke Database SQLite (gamon.db)
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	defer db.Close()

	// 2. Inisialisasi Hub WebSocket untuk pengiriman status real-time ke browser
	hub := handler.NewHub(db)
	go hub.Run() // Dijalankan secara asynchronous (goroutine)

	// 3. Inisialisasi Pengirim Peringatan Telegram Bot
	notifier := notification.NewTelegramNotifier(db)
	telegramHandler := handler.NewTelegramHandler(db)

	// Cek apakah Notifikasi Telegram diaktifkan via variabel lingkungan (TELEGRAM_BOT_TOKEN)
	if notifier.IsEnabled() {
		log.Println("Notifikasi Telegram: AKTIF")
		poller := notification.NewTelegramPoller(os.Getenv("TELEGRAM_BOT_TOKEN"), telegramHandler.CompletePairing, telegramHandler.IsChatIDActive)
		go poller.Start() // Menjalankan poller perintah Telegram di latar belakang
	} else {
		log.Println("Notifikasi Telegram: NONAKTIF (Set TELEGRAM_BOT_TOKEN untuk mengaktifkan)")
	}

	// 4. Inisialisasi Mesin Pemantau (Engine Monitoring ICMP Ping)
	engine := monitor.NewEngine(hub, db, monitor.WithNotifier(notifier))

	// 5. Inisialisasi Pengelola Rute (Handlers) REST API
	deviceHandler := handler.NewDeviceHandler(db, engine, hub)
	alertHandler := handler.NewAlertHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)
	monitoringHandler := handler.NewMonitoringHandler(db)
	settingsHandler := handler.NewSettingsHandler(db)
	legacyAPI := handler.NewAPI(engine, hub)

	// Membuat Mux (Router HTTP bawaan Go)
	mux := http.NewServeMux()

	// Pendaftaran Alamat WebSocket (Real-time update)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(hub, w, r)
	})

	// Pendaftaran Alamat REST API untuk Perangkat & Alert
	mux.HandleFunc("/api/devices", deviceHandler.HandleDevices)
	mux.HandleFunc("/api/devices/", deviceHandler.HandleDevice)
	mux.HandleFunc("/api/alerts", alertHandler.HandleAlerts)
	mux.HandleFunc("/api/alerts/", alertHandler.HandleAlert)
	mux.HandleFunc("/api/dashboard", dashboardHandler.HandleDashboard)
	mux.HandleFunc("/api/monitoring", monitoringHandler.HandleMonitoring)
	mux.HandleFunc("/api/monitoring/", monitoringHandler.HandleMonitoringDevice)
	mux.HandleFunc("/api/telegram/pair", telegramHandler.HandlePair)
	mux.HandleFunc("/api/telegram/status", telegramHandler.HandleStatus)
	mux.HandleFunc("/api/telegram/disconnect", telegramHandler.HandleDisconnect)
	mux.HandleFunc("/api/settings", settingsHandler.HandleSettings)
	mux.HandleFunc("/api/health", legacyAPI.Health)

	// Membungkus router dengan Middleware CORS (agar frontend bisa mengakses backend)
	wrapped := corsMiddleware(mux)

	// 6. Otomatis mulai pemantauan untuk seluruh perangkat berstatus 'active' di database
	go autoStartMonitoring(db, engine)

	fmt.Println("===========================================")
	fmt.Println("   GAMON - Garda Monitoring v0.4")
	fmt.Println("   Web Backend Server (Golang)")
	fmt.Println("   Berjalan di http://localhost:8080")
	fmt.Println("===========================================")

	// 7. Jalankan HTTP Web Server di port 8080
	log.Fatal(http.ListenAndServe(":8080", wrapped))
}

// autoStartMonitoring membaca database dan otomatis memulai pemantauan IP perangkat saat server dinyalakan.
func autoStartMonitoring(db *sql.DB, engine *monitor.Engine) {
	time.Sleep(2 * time.Second) // Tunda 2 detik memastikan seluruh komponen siap

	rows, err := db.Query("SELECT id, ip, method, url, port, check_interval FROM devices WHERE status = 'active'")
	if err != nil {
		log.Printf("Gagal membaca daftar perangkat dari database: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, interval int
		var ip, method, url string
		var port *int

		if err := rows.Scan(&id, &ip, &method, &url, &port, &interval); err != nil {
			log.Printf("Gagal membaca baris perangkat: %v", err)
			continue
		}

		config := monitor.DeviceConfig{
			DeviceID: id,
			IP:       ip,
			URL:      url,
			Method:   method,
			Interval: interval,
		}
		if port != nil {
			config.Port = *port
		}

		engine.Start(config)
		count++
		log.Printf("Otomatis memulai pemantauan perangkat ID %d (%s)", id, ip)
	}

	if count > 0 {
		log.Printf("Berhasil otomatis memulai pemantauan untuk %d perangkat aktif", count)
	}
}

// corsMiddleware memberikan izin Cross-Origin Resource Sharing (CORS) agar aplikasi web React (frontend)
// dapat berkomunikasi dengan backend server tanpa diblokir oleh browser.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
