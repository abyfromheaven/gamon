package main

import (
	"fmt"
	"log"
	"net/http"

	"gamon/handler"
	"gamon/monitor"
)

func main() {
	hub := handler.NewHub()
	go hub.Run()

	engine := monitor.NewEngine(hub)

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleWebSocket(hub, w, r)
	})

	api := handler.NewAPI(engine, hub)
	mux.HandleFunc("/api/monitor", api.StartMonitor)
	mux.HandleFunc("/api/stop", api.StopMonitor)
	mux.HandleFunc("/api/health", api.Health)

	wrapped := corsMiddleware(mux)

	fmt.Println("===========================================")
	fmt.Println("   GAMON - Garda Monitoring v0.2")
	fmt.Println("   Web Backend Server")
	fmt.Println("   Running on http://localhost:8080")
	fmt.Println("===========================================")

	log.Fatal(http.ListenAndServe(":8080", wrapped))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
