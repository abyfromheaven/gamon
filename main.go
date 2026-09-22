package main

// ============================================================================
// MODUL UTAMA APKLIKASI (main.go)
// ============================================================================
// GAMON (Garda Monitoring) - Sistem Pemonitor Jaringan & Perangkat Berbasis Web.
// Berkas main.go merupakan titik masuk utama (entry point) aplikasi backend Golang.
//
// Urutan alur inisialisasi aplikasi saat dinyalakan:
// 1. Membuka & menyiapkan Database SQLite (gamon.db) dengan WAL Mode.
// 2. Menjalankan Pusat WebSocket (Hub) untuk pengiriman data real-time ke web React.
// 3. Menyiapkan Pengirim Notifikasi & Poller Telegram Bot.
// 4. Menyiapkan Mesin Pemantau Jaringan (ICMP Ping Engine).
// 5. Mendaftarkan seluruh alamat rute API (REST Endpoints) & Middleware CORS.
// 6. Otomatis mulai memantau perangkat yang berstatus 'active' di database.
// 7. Menjalankan HTTP Web Server di port 8080.
// ============================================================================

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

func main() {
	// 1. Inisialisasi dan koneksi ke Database SQLite (gamon.db)
	koneksiDb, err := database.NewDB()
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	defer koneksiDb.Close()

	// 2. Inisialisasi Pusat WebSocket (Hub) untuk penyiaran data real-time ke browser
	pusatWebsocket := handler.NewHub(koneksiDb)
	go pusatWebsocket.Run() // Dijalankan secara asynchronous di goroutine terpisah

	// 3. Inisialisasi Pengirim Notifikasi & Penangan Telegram Bot
	pengirimTelegram := notification.NewTelegramNotifier(koneksiDb)
	pengelolaTelegram := handler.NewTelegramHandler(koneksiDb)

	// Cek apakah token bot Telegram diatur di variabel lingkungan (TELEGRAM_BOT_TOKEN)
	if pengirimTelegram.IsEnabled() {
		log.Println("Notifikasi Telegram: AKTIF")
		pengecekPesanTelegram := notification.NewTelegramPoller(os.Getenv("TELEGRAM_BOT_TOKEN"), pengelolaTelegram.CompletePairing, pengelolaTelegram.IsChatIDActive)
		go pengecekPesanTelegram.Start() // Menjalankan poller perintah Telegram di latar belakang
	} else {
		log.Println("Notifikasi Telegram: NONAKTIF (Set TELEGRAM_BOT_TOKEN untuk mengaktifkan)")
	}

	// 4. Inisialisasi Mesin Pemantau (Engine Monitoring ICMP Ping)
	mesinPemantau := monitor.NewEngine(pusatWebsocket, koneksiDb, monitor.WithNotifier(pengirimTelegram))

	// 5. Inisialisasi Pengelola Rute (Handlers) REST API
	pengelolaPerangkat := handler.NewDeviceHandler(koneksiDb, mesinPemantau, pusatWebsocket)
	pengelolaPeringatan := handler.NewAlertHandler(koneksiDb)
	pengelolaDashboard := handler.NewDashboardHandler(koneksiDb)
	pengelolaMonitoring := handler.NewMonitoringHandler(koneksiDb)
	pengelolaPengaturan := handler.NewSettingsHandler(koneksiDb)
	apiKesehatan := handler.NewAPI(mesinPemantau, pusatWebsocket)

	// Membuat Router HTTP bawaan Golang (ServeMux)
	routerHTTP := http.NewServeMux()

	// Pendaftaran Alamat WebSocket (Real-time update)
	routerHTTP.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(pusatWebsocket, w, r)
	})

	// Pendaftaran Alamat REST API
	routerHTTP.HandleFunc("/api/devices", pengelolaPerangkat.HandleDevices)
	routerHTTP.HandleFunc("/api/devices/", pengelolaPerangkat.HandleDevice)
	routerHTTP.HandleFunc("/api/alerts", pengelolaPeringatan.HandleAlerts)
	routerHTTP.HandleFunc("/api/alerts/", pengelolaPeringatan.HandleAlert)
	routerHTTP.HandleFunc("/api/dashboard", pengelolaDashboard.HandleDashboard)
	routerHTTP.HandleFunc("/api/monitoring", pengelolaMonitoring.HandleMonitoring)
	routerHTTP.HandleFunc("/api/monitoring/", pengelolaMonitoring.HandleMonitoringDevice)
	routerHTTP.HandleFunc("/api/telegram/pair", pengelolaTelegram.HandlePair)
	routerHTTP.HandleFunc("/api/telegram/status", pengelolaTelegram.HandleStatus)
	routerHTTP.HandleFunc("/api/telegram/disconnect", pengelolaTelegram.HandleDisconnect)
	routerHTTP.HandleFunc("/api/settings", pengelolaPengaturan.HandleSettings)
	routerHTTP.HandleFunc("/api/health", apiKesehatan.Health)

	// Membungkus router dengan Middleware CORS (izin akses lintas domain dari frontend React)
	routerDenganCORS := corsMiddleware(routerHTTP)

	// 6. Otomatis mulai pemantauan untuk seluruh perangkat berstatus 'active'
	go autoStartMonitoring(koneksiDb, mesinPemantau)

	fmt.Println("===========================================")
	fmt.Println("   GAMON - Garda Monitoring v0.4")
	fmt.Println("   Web Backend Server (Golang)")
	fmt.Println("   Berjalan di http://localhost:8080")
	fmt.Println("===========================================")

	// 7. Jalankan HTTP Web Server di port 8080
	log.Fatal(http.ListenAndServe(":8080", routerDenganCORS))
}

// autoStartMonitoring membaca database dan otomatis memulai pemantauan IP perangkat saat server dinyalakan.
func autoStartMonitoring(koneksiDb *sql.DB, mesinPemantau *monitor.Engine) {
	time.Sleep(2 * time.Second) // Tunda 2 detik memastikan seluruh komponen server siap

	barisData, err := koneksiDb.Query("SELECT id, ip, method, url, port, check_interval FROM devices WHERE status = 'active'")
	if err != nil {
		log.Printf("Gagal membaca daftar perangkat dari database: %v", err)
		return
	}
	defer barisData.Close()

	jumlahDipantau := 0
	for barisData.Next() {
		var idPerangkat, interval int
		var ip, metode, url string
		var port *int

		if err := barisData.Scan(&idPerangkat, &ip, &metode, &url, &port, &interval); err != nil {
			log.Printf("Gagal membaca baris data perangkat: %v", err)
			continue
		}

		konfigurasi := monitor.DeviceConfig{
			DeviceID: idPerangkat,
			IP:       ip,
			URL:      url,
			Method:   metode,
			Interval: interval,
		}
		if port != nil {
			konfigurasi.Port = *port
		}

		mesinPemantau.Start(konfigurasi)
		jumlahDipantau++
		log.Printf("Otomatis memulai pemantauan perangkat ID %d (%s)", idPerangkat, ip)
	}

	if jumlahDipantau > 0 {
		log.Printf("Berhasil otomatis memulai pemantauan untuk %d perangkat aktif", jumlahDipantau)
	}
}

// corsMiddleware memberikan izin Cross-Origin Resource Sharing (CORS) agar aplikasi web React (frontend)
// dapat berkomunikasi dengan backend Golang tanpa diblokir oleh kebijakan keamanan browser.
func corsMiddleware(selanjutnya http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		selanjutnya.ServeHTTP(w, r)
	})
}
