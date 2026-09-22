package handler

// ============================================================================
// MODUL HANDLER KESEHATAN API (handler/api.go)
// ============================================================================
// Modul ini mengelola endpoint pengecekan kesehatan server backend (/api/health).
// ============================================================================

import (
	"encoding/json"
	"net/http"

	"gamon/monitor"
)

// API struct pengelola rute legacy / health backend.
type API struct {
	engine *monitor.Engine
	hub    *Hub
}

// NewAPI membuat instansi baru pengelola API health.
func NewAPI(engine *monitor.Engine, hub *Hub) *API {
	return &API{engine: engine, hub: hub}
}

// HealthResponse format balasan tes kesehatan server.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health mengecek status server backend backend golang (merespons {"status": "ok"}).
func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
