package handler

import (
	"encoding/json"
	"net/http"

	"gamon/monitor"
)

type API struct {
	engine *monitor.Engine
	hub    *Hub
}

func NewAPI(engine *monitor.Engine, hub *Hub) *API {
	return &API{engine: engine, hub: hub}
}

type MonitorRequest struct {
	IP string `json:"ip"`
}

type StopRequest struct {
	IP string `json:"ip"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (a *API) StartMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req MonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid request body"})
		return
	}

	if req.IP == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "IP address is required"})
		return
	}

	a.engine.Start(req.IP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Monitoring started for " + req.IP})
}

func (a *API) StopMonitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Method not allowed"})
		return
	}

	var req StopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Invalid request body"})
		return
	}

	a.engine.Stop(req.IP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{Success: true, Message: "Monitoring stopped for " + req.IP})
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
