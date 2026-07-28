package main

import (
	"fmt"
	"log"
	"net/http"

	"gamon/database"
	"gamon/handler"
	"gamon/monitor"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	hub := handler.NewHub()
	go hub.Run()

	engine := monitor.NewEngine(hub)

	deviceHandler := handler.NewDeviceHandler(db, engine, hub)
	alertHandler := handler.NewAlertHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)
	monitoringHandler := handler.NewMonitoringHandler(db)
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
	mux.HandleFunc("/api/health", legacyAPI.Health)

	wrapped := corsMiddleware(mux)

	fmt.Println("===========================================")
	fmt.Println("   GAMON - Garda Monitoring v0.3")
	fmt.Println("   Web Backend Server")
	fmt.Println("   Running on http://localhost:8080")
	fmt.Println("===========================================")

	log.Fatal(http.ListenAndServe(":8080", wrapped))
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
