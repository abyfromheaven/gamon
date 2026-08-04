package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AlertHandler struct {
	db *sql.DB
}

func NewAlertHandler(db *sql.DB) *AlertHandler {
	return &AlertHandler{db: db}
}

func (h *AlertHandler) HandleAlerts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listAlerts(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *AlertHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
	idStr = strings.Split(idStr, "/")[0]

	if idStr == "count" {
		h.alertCount(w, r)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/alerts/"+idStr), "/")

	if len(parts) > 1 {
		switch parts[1] {
		case "resolve":
			h.resolveAlert(w, r, id)
			return
		case "acknowledge":
			h.acknowledgeAlert(w, r, id)
			return
		}
	}

	switch r.Method {
	case http.MethodGet:
		h.getAlert(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *AlertHandler) listAlerts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT a.id, a.device_id, d.name, d.type, d.ip, d.method,
		a.title, a.status, a.severity,
		a.started_at, a.resolved_at, a.description,
		a.acknowledged, a.acknowledged_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE 1=1`
	args := []interface{}{}

	if status := r.URL.Query().Get("status"); status != "" {
		query += " AND a.status = ?"
		args = append(args, status)
	}
	if severity := r.URL.Query().Get("severity"); severity != "" {
		query += " AND a.severity = ?"
		args = append(args, severity)
	}
	if deviceType := r.URL.Query().Get("device_type"); deviceType != "" {
		query += " AND d.type = ?"
		args = append(args, deviceType)
	}

	query += " ORDER BY a.started_at DESC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error listing alerts: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list alerts")
		return
	}
	defer rows.Close()

	type AlertResponse struct {
		ID             int     `json:"id"`
		DeviceID       int     `json:"device_id"`
		DeviceName     string  `json:"device_name"`
		DeviceType     string  `json:"device_type"`
		DeviceIP       string  `json:"device_ip"`
		Method         string  `json:"method"`
		Title          string  `json:"title"`
		Status         string  `json:"status"`
		Severity       string  `json:"severity"`
		StartedAt      string  `json:"started_at"`
		ResolvedAt     *string `json:"resolved_at"`
		Description    string  `json:"description"`
		Acknowledged   bool    `json:"acknowledged"`
		AcknowledgedAt *string `json:"acknowledged_at"`
	}

	var alerts []AlertResponse
	for rows.Next() {
		var a AlertResponse
		var startedAt string
		var resolvedAt *string
		var acknowledgedAt *string
		if err := rows.Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.DeviceType, &a.DeviceIP, &a.Method,
			&a.Title, &a.Status, &a.Severity,
			&startedAt, &resolvedAt, &a.Description,
			&a.Acknowledged, &acknowledgedAt); err != nil {
			log.Printf("Error scanning alert: %v", err)
			continue
		}
		a.StartedAt = startedAt
		a.ResolvedAt = resolvedAt
		a.AcknowledgedAt = acknowledgedAt
		alerts = append(alerts, a)
	}

	if alerts == nil {
		alerts = []AlertResponse{}
	}
	respondData(w, alerts)
}

func (h *AlertHandler) getAlert(w http.ResponseWriter, _ *http.Request, id int) {
	type AlertResponse struct {
		ID             int     `json:"id"`
		DeviceID       int     `json:"device_id"`
		DeviceName     string  `json:"device_name"`
		DeviceType     string  `json:"device_type"`
		DeviceIP       string  `json:"device_ip"`
		Title          string  `json:"title"`
		Status         string  `json:"status"`
		Severity       string  `json:"severity"`
		StartedAt      string  `json:"started_at"`
		ResolvedAt     *string `json:"resolved_at"`
		Description    string  `json:"description"`
		Acknowledged   bool    `json:"acknowledged"`
		AcknowledgedAt *string `json:"acknowledged_at"`
	}

	var a AlertResponse
	var startedAt string
	var resolvedAt *string
	var acknowledgedAt *string
	err := h.db.QueryRow(`SELECT a.id, a.device_id, d.name, d.type, d.ip, a.title, a.status, a.severity,
		a.started_at, a.resolved_at, a.description, a.acknowledged, a.acknowledged_at
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE a.id = ?`, id).
		Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.DeviceType, &a.DeviceIP, &a.Title, &a.Status, &a.Severity,
			&startedAt, &resolvedAt, &a.Description, &a.Acknowledged, &acknowledgedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Alert not found")
		return
	}
	if err != nil {
		log.Printf("Error getting alert: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get alert")
		return
	}
	a.StartedAt = startedAt
	a.ResolvedAt = resolvedAt
	a.AcknowledgedAt = acknowledgedAt
	respondData(w, a)
}

func (h *AlertHandler) resolveAlert(w http.ResponseWriter, _ *http.Request, id int) {
	result, err := h.db.Exec("UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP WHERE id = ? AND status = 'ongoing'", id)
	if err != nil {
		log.Printf("Error resolving alert: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to resolve alert")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Alert not found or already resolved")
		return
	}

	log.Printf("Alert resolved: %d", id)
	respondSuccess(w, "Alert resolved successfully")
}

func (h *AlertHandler) acknowledgeAlert(w http.ResponseWriter, _ *http.Request, id int) {
	result, err := h.db.Exec("UPDATE alerts SET acknowledged = TRUE, acknowledged_at = CURRENT_TIMESTAMP WHERE id = ? AND acknowledged = FALSE", id)
	if err != nil {
		log.Printf("Error acknowledging alert: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to acknowledge alert")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Alert not found or already acknowledged")
		return
	}

	log.Printf("Alert acknowledged: %d", id)
	respondSuccess(w, "Alert acknowledged successfully")
}

func (h *AlertHandler) alertCount(w http.ResponseWriter, _ *http.Request) {
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM alerts WHERE status = 'ongoing'").Scan(&count)
	if err != nil {
		log.Printf("Error counting alerts: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to count alerts")
		return
	}

	respondData(w, map[string]int{"ongoing": count})
}
