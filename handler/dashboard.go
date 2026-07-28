package handler

import (
	"database/sql"
	"log"
	"net/http"
)

type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	type Summary struct {
		TotalDevices   int `json:"total_devices"`
		OnlineDevices  int `json:"online_devices"`
		OfflineDevices int `json:"offline_devices"`
		WarningDevices int `json:"warning_devices"`
	}

	type LatestAlert struct {
		ID         int    `json:"id"`
		DeviceName string `json:"device_name"`
		Title      string `json:"title"`
		Severity   string `json:"severity"`
		Status     string `json:"status"`
		StartedAt  string `json:"started_at"`
	}

	type DashboardData struct {
		Summary      Summary      `json:"summary"`
		LatestAlerts []LatestAlert `json:"latest_alerts"`
	}

	var summary Summary
	err := h.db.QueryRow("SELECT COUNT(*) FROM devices").Scan(&summary.TotalDevices)
	if err != nil {
		log.Printf("Error counting total devices: %v", err)
	}

	err = h.db.QueryRow(`SELECT COUNT(DISTINCT d.id)
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		WHERE ph.status = 'online'`).Scan(&summary.OnlineDevices)
	if err != nil {
		log.Printf("Error counting online devices: %v", err)
	}

	err = h.db.QueryRow(`SELECT COUNT(DISTINCT d.id)
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		WHERE ph.status = 'offline'`).Scan(&summary.OfflineDevices)
	if err != nil {
		log.Printf("Error counting offline devices: %v", err)
	}

	err = h.db.QueryRow(`SELECT COUNT(DISTINCT d.id)
		FROM devices d
		LEFT JOIN ping_history ph ON d.id = ph.device_id AND ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY timestamp DESC LIMIT 1
		)
		WHERE ph.status = 'warning'`).Scan(&summary.WarningDevices)
	if err != nil {
		log.Printf("Error counting warning devices: %v", err)
	}

	rows, err := h.db.Query(`SELECT a.id, d.name, a.title, a.severity, a.status, a.started_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		ORDER BY a.started_at DESC
		LIMIT 5`)
	if err != nil {
		log.Printf("Error listing latest alerts: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get dashboard")
		return
	}
	defer rows.Close()

	var latestAlerts []LatestAlert
	for rows.Next() {
		var a LatestAlert
		var startedAt string
		if err := rows.Scan(&a.ID, &a.DeviceName, &a.Title, &a.Severity, &a.Status, &startedAt); err != nil {
			log.Printf("Error scanning alert: %v", err)
			continue
		}
		a.StartedAt = startedAt
		latestAlerts = append(latestAlerts, a)
	}

	if latestAlerts == nil {
		latestAlerts = []LatestAlert{}
	}

	data := DashboardData{
		Summary:      summary,
		LatestAlerts: latestAlerts,
	}

	respondData(w, data)
}
