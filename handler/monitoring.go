package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type MonitoringHandler struct {
	db *sql.DB
}

func NewMonitoringHandler(db *sql.DB) *MonitoringHandler {
	return &MonitoringHandler{db: db}
}

func (h *MonitoringHandler) HandleMonitoring(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.listMonitoringStatus(w, r)
}

func (h *MonitoringHandler) HandleMonitoringDevice(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/monitoring/")
	parts := strings.Split(path, "/")

	idStr := parts[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	if len(parts) > 1 && parts[1] == "history" {
		h.getDeviceHistory(w, r, id)
		return
	}

	respondError(w, http.StatusNotFound, "Endpoint not found")
}

func (h *MonitoringHandler) listMonitoringStatus(w http.ResponseWriter, _ *http.Request) {
	type MonitoringStatus struct {
		DeviceID   int     `json:"device_id"`
		DeviceName string  `json:"device_name"`
		DeviceType string  `json:"device_type"`
		IP         string  `json:"ip"`
		Method     string  `json:"method"`
		Status     string  `json:"status"`
		LatencyMs  float64 `json:"latency_ms"`
		LastCheck  *string `json:"last_check"`
		Interval   int     `json:"interval"`
	}

	rows, err := h.db.Query(`SELECT d.id, d.name, d.type, d.ip, d.method, d.check_interval,
		COALESCE(ph.status, 'unknown') as status,
		COALESCE(ph.latency_ms, 0) as latency_ms,
		COALESCE(ph.timestamp, '') as last_check
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		ORDER BY d.name`)
	if err != nil {
		log.Printf("Error listing monitoring status: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list monitoring status")
		return
	}
	defer rows.Close()

	var statuses []MonitoringStatus
	for rows.Next() {
		var s MonitoringStatus
		var lastCheck string
		if err := rows.Scan(&s.DeviceID, &s.DeviceName, &s.DeviceType, &s.IP, &s.Method, &s.Interval, &s.Status, &s.LatencyMs, &lastCheck); err != nil {
			log.Printf("Error scanning monitoring status: %v", err)
			continue
		}
		if lastCheck != "" {
			s.LastCheck = &lastCheck
		}
		statuses = append(statuses, s)
	}

	if statuses == nil {
		statuses = []MonitoringStatus{}
	}
	respondData(w, statuses)
}

func (h *MonitoringHandler) getDeviceHistory(w http.ResponseWriter, _ *http.Request, deviceID int) {
	type PingRecord struct {
		ID        int     `json:"id"`
		Status    string  `json:"status"`
		LatencyMs float64 `json:"latency_ms"`
		TTL       int     `json:"ttl"`
		Seq       int     `json:"seq"`
		Details   string  `json:"details"`
		Timestamp string  `json:"timestamp"`
	}

	rows, err := h.db.Query(`SELECT id, status, latency_ms, ttl, seq, details, timestamp
		FROM ping_history
		WHERE device_id = ?
		ORDER BY timestamp DESC
		LIMIT 50`, deviceID)
	if err != nil {
		log.Printf("Error getting device history: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get device history")
		return
	}
	defer rows.Close()

	var history []PingRecord
	for rows.Next() {
		var p PingRecord
		var timestamp string
		if err := rows.Scan(&p.ID, &p.Status, &p.LatencyMs, &p.TTL, &p.Seq, &p.Details, &timestamp); err != nil {
			log.Printf("Error scanning ping record: %v", err)
			continue
		}
		p.Timestamp = timestamp
		history = append(history, p)
	}

	if history == nil {
		history = []PingRecord{}
	}
	respondData(w, history)
}
