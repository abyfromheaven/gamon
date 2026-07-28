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
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	h.listAlerts(w, r)
}

func (h *AlertHandler) HandleAlert(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/alerts/")
	idStr = strings.Split(idStr, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid alert ID")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/alerts/"+idStr), "/")

	if len(parts) > 1 && parts[1] == "resolve" {
		h.resolveAlert(w, r, id)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getAlert(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *AlertHandler) listAlerts(w http.ResponseWriter, r *http.Request) {
	query := `SELECT a.id, a.device_id, d.name, d.type, a.title, a.status, a.severity, a.started_at, a.resolved_at, a.description
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
		ID          int     `json:"id"`
		DeviceID    int     `json:"device_id"`
		DeviceName  string  `json:"device_name"`
		DeviceType  string  `json:"device_type"`
		Title       string  `json:"title"`
		Status      string  `json:"status"`
		Severity    string  `json:"severity"`
		StartedAt   string  `json:"started_at"`
		ResolvedAt  *string `json:"resolved_at"`
		Description string  `json:"description"`
	}

	var alerts []AlertResponse
	for rows.Next() {
		var a AlertResponse
		var startedAt string
		var resolvedAt *string
		if err := rows.Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.DeviceType, &a.Title, &a.Status, &a.Severity, &startedAt, &resolvedAt, &a.Description); err != nil {
			log.Printf("Error scanning alert: %v", err)
			continue
		}
		a.StartedAt = startedAt
		a.ResolvedAt = resolvedAt
		alerts = append(alerts, a)
	}

	if alerts == nil {
		alerts = []AlertResponse{}
	}
	respondData(w, alerts)
}

func (h *AlertHandler) getAlert(w http.ResponseWriter, _ *http.Request, id int) {
	type AlertResponse struct {
		ID          int     `json:"id"`
		DeviceID    int     `json:"device_id"`
		DeviceName  string  `json:"device_name"`
		DeviceType  string  `json:"device_type"`
		DeviceIP    string  `json:"device_ip"`
		Title       string  `json:"title"`
		Status      string  `json:"status"`
		Severity    string  `json:"severity"`
		StartedAt   string  `json:"started_at"`
		ResolvedAt  *string `json:"resolved_at"`
		Description string  `json:"description"`
	}

	var a AlertResponse
	var startedAt string
	var resolvedAt *string
	err := h.db.QueryRow(`SELECT a.id, a.device_id, d.name, d.type, d.ip, a.title, a.status, a.severity, a.started_at, a.resolved_at, a.description
		FROM alerts a
		JOIN devices d ON a.device_id = d.id
		WHERE a.id = ?`, id).
		Scan(&a.ID, &a.DeviceID, &a.DeviceName, &a.DeviceType, &a.DeviceIP, &a.Title, &a.Status, &a.Severity, &startedAt, &resolvedAt, &a.Description)
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
