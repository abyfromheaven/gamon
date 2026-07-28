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

type HealthResponse struct {
	Status string `json:"status"`
}

func (a *API) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
}
