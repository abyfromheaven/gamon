package main

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
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	hub := handler.NewHub(db)
	go hub.Run()

	notifier := notification.NewTelegramNotifier(db)
	telegramHandler := handler.NewTelegramHandler(db)

	if notifier.IsEnabled() {
		log.Println("Telegram notification: enabled")
		poller := notification.NewTelegramPoller(os.Getenv("TELEGRAM_BOT_TOKEN"), telegramHandler.CompletePairing, telegramHandler.IsChatIDActive)
		go poller.Start()
	} else {
		log.Println("Telegram notification: disabled (set TELEGRAM_BOT_TOKEN to enable)")
	}

	engine := monitor.NewEngine(hub, db, monitor.WithNotifier(notifier))

	deviceHandler := handler.NewDeviceHandler(db, engine, hub)
	alertHandler := handler.NewAlertHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)
	monitoringHandler := handler.NewMonitoringHandler(db)
	settingsHandler := handler.NewSettingsHandler(db)
	legacyAPI := handler.NewAPI(engine, hub)

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(hub, w, r)
	})

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

	wrapped := corsMiddleware(mux)

	go autoStartMonitoring(db, engine)

	fmt.Println("===========================================")
	fmt.Println("   GAMON - Garda Monitoring v0.4")
	fmt.Println("   Web Backend Server")
	fmt.Println("   Running on http://localhost:8080")
	fmt.Println("===========================================")

	log.Fatal(http.ListenAndServe(":8080", wrapped))
}

func autoStartMonitoring(db *sql.DB, engine *monitor.Engine) {
	time.Sleep(2 * time.Second)

	rows, err := db.Query("SELECT id, ip, method, url, port, check_interval FROM devices WHERE status = 'active'")
	if err != nil {
		log.Printf("Failed to query devices for auto-start: %v", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id, interval int
		var ip, method, url string
		var port *int

		if err := rows.Scan(&id, &ip, &method, &url, &port, &interval); err != nil {
			log.Printf("Error scanning device for auto-start: %v", err)
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
		log.Printf("Auto-started monitoring for device %d (%s)", id, ip)
	}

	if count > 0 {
		log.Printf("Auto-started monitoring for %d active devices", count)
	}
}

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
