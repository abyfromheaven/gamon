package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gamon/monitor"
)

type DeviceHandler struct {
	db     *sql.DB
	engine *monitor.Engine
	hub    *Hub
}

func NewDeviceHandler(db *sql.DB, engine *monitor.Engine, hub *Hub) *DeviceHandler {
	return &DeviceHandler{db: db, engine: engine, hub: hub}
}

type CreateDeviceRequest struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	IP            string `json:"ip"`
	URL           string `json:"url"`
	Port          *int   `json:"port"`
	Method        string `json:"method"`
	Location      string `json:"location"`
	CheckInterval int    `json:"check_interval"`
	Status        string `json:"status"`
	Description   string `json:"description"`
}

type UpdateDeviceRequest struct {
	Name          *string `json:"name"`
	Type          *string `json:"type"`
	IP            *string `json:"ip"`
	URL           *string `json:"url"`
	Port          *int    `json:"port"`
	Method        *string `json:"method"`
	Location      *string `json:"location"`
	CheckInterval *int    `json:"check_interval"`
	Status        *string `json:"status"`
	Description   *string `json:"description"`
}

type ToggleStatusRequest struct {
	Status string `json:"status"`
}

func (h *DeviceHandler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDevices(w, r)
	case http.MethodPost:
		h.createDevice(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *DeviceHandler) HandleDevice(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	idStr = strings.Split(idStr, "/")[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid device ID")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/devices/"+idStr), "/")

	if len(parts) > 1 && parts[1] == "start" {
		h.startMonitoring(w, r, id)
		return
	}
	if len(parts) > 1 && parts[1] == "stop" {
		h.stopMonitoring(w, r, id)
		return
	}
	if len(parts) > 1 && parts[1] == "status" {
		h.toggleStatus(w, r, id)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getDevice(w, r, id)
	case http.MethodPut:
		h.updateDevice(w, r, id)
	case http.MethodDelete:
		h.deleteDevice(w, r, id)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *DeviceHandler) listDevices(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.db.Query("SELECT id, name, type, ip, url, port, method, location, check_interval, status, description, created_at, updated_at FROM devices ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Error listing devices: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list devices")
		return
	}
	defer rows.Close()

	type DeviceResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		IP            string `json:"ip"`
		URL           string `json:"url"`
		Port          *int   `json:"port"`
		Method        string `json:"method"`
		Location      string `json:"location"`
		CheckInterval int    `json:"check_interval"`
		Status        string `json:"status"`
		Description   string `json:"description"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}

	var devices []DeviceResponse
	for rows.Next() {
		var d DeviceResponse
		var createdAt, updatedAt string
		if err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.IP, &d.URL, &d.Port, &d.Method, &d.Location, &d.CheckInterval, &d.Status, &d.Description, &createdAt, &updatedAt); err != nil {
			log.Printf("Error scanning device: %v", err)
			continue
		}
		d.CreatedAt = createdAt
		d.UpdatedAt = updatedAt
		devices = append(devices, d)
	}

	if devices == nil {
		devices = []DeviceResponse{}
	}
	respondData(w, devices)
}

func (h *DeviceHandler) getDevice(w http.ResponseWriter, _ *http.Request, id int) {
	type DeviceResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		IP            string `json:"ip"`
		URL           string `json:"url"`
		Port          *int   `json:"port"`
		Method        string `json:"method"`
		Location      string `json:"location"`
		CheckInterval int    `json:"check_interval"`
		Status        string `json:"status"`
		Description   string `json:"description"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
	}

	var d DeviceResponse
	var createdAt, updatedAt string
	err := h.db.QueryRow("SELECT id, name, type, ip, url, port, method, location, check_interval, status, description, created_at, updated_at FROM devices WHERE id = ?", id).
		Scan(&d.ID, &d.Name, &d.Type, &d.IP, &d.URL, &d.Port, &d.Method, &d.Location, &d.CheckInterval, &d.Status, &d.Description, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Device not found")
		return
	}
	if err != nil {
		log.Printf("Error getting device: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get device")
		return
	}
	d.CreatedAt = createdAt
	d.UpdatedAt = updatedAt
	respondData(w, d)
}

func (h *DeviceHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Type == "" {
		respondError(w, http.StatusBadRequest, "Type is required")
		return
	}
	if req.IP == "" {
		respondError(w, http.StatusBadRequest, "IP is required")
		return
	}
	if req.Method == "" {
		req.Method = "ICMP Ping"
	}
	if req.CheckInterval <= 0 {
		req.CheckInterval = 3
	}
	if req.Status == "" {
		req.Status = "active"
	}
	if req.Status != "active" && req.Status != "inactive" {
		respondError(w, http.StatusBadRequest, "Status must be 'active' or 'inactive'")
		return
	}

	result, err := h.db.Exec(
		"INSERT INTO devices (name, type, ip, url, port, method, location, check_interval, status, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Name, req.Type, req.IP, req.URL, req.Port, req.Method, req.Location, req.CheckInterval, req.Status, req.Description,
	)
	if err != nil {
		log.Printf("Error creating device: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create device")
		return
	}

	id, _ := result.LastInsertId()

	type DeviceResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		IP            string `json:"ip"`
		URL           string `json:"url"`
		Port          *int   `json:"port"`
		Method        string `json:"method"`
		Location      string `json:"location"`
		CheckInterval int    `json:"check_interval"`
		Status        string `json:"status"`
		Description   string `json:"description"`
	}

	device := DeviceResponse{
		ID:            int(id),
		Name:          req.Name,
		Type:          req.Type,
		IP:            req.IP,
		URL:           req.URL,
		Port:          req.Port,
		Method:        req.Method,
		Location:      req.Location,
		CheckInterval: req.CheckInterval,
		Status:        req.Status,
		Description:   req.Description,
	}

	log.Printf("Device created: %s (%s)", req.Name, req.IP)
	respondData(w, device)
}

func (h *DeviceHandler) updateDevice(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var currentName, currentType, currentIP, currentURL, currentMethod, currentLocation, currentStatus, currentDescription string
	var currentPort *int
	var currentInterval int
	err := h.db.QueryRow("SELECT name, type, ip, url, port, method, location, check_interval, status, description FROM devices WHERE id = ?", id).
		Scan(&currentName, &currentType, &currentIP, &currentURL, &currentPort, &currentMethod, &currentLocation, &currentInterval, &currentStatus, &currentDescription)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Device not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get device")
		return
	}

	name := currentName
	typ := currentType
	ip := currentIP
	url := currentURL
	port := currentPort
	method := currentMethod
	location := currentLocation
	interval := currentInterval
	status := currentStatus
	desc := currentDescription

	if req.Name != nil {
		name = *req.Name
	}
	if req.Type != nil {
		typ = *req.Type
	}
	if req.IP != nil {
		ip = *req.IP
	}
	if req.URL != nil {
		url = *req.URL
	}
	if req.Port != nil {
		port = req.Port
	}
	if req.Method != nil {
		method = *req.Method
	}
	if req.Location != nil {
		location = *req.Location
	}
	if req.CheckInterval != nil {
		interval = *req.CheckInterval
	}
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "inactive" {
			respondError(w, http.StatusBadRequest, "Status must be 'active' or 'inactive'")
			return
		}
		status = *req.Status
	}
	if req.Description != nil {
		desc = *req.Description
	}

	_, err = h.db.Exec(
		"UPDATE devices SET name = ?, type = ?, ip = ?, url = ?, port = ?, method = ?, location = ?, check_interval = ?, status = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		name, typ, ip, url, port, method, location, interval, status, desc, id,
	)
	if err != nil {
		log.Printf("Error updating device: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to update device")
		return
	}

	type DeviceResponse struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Type          string `json:"type"`
		IP            string `json:"ip"`
		URL           string `json:"url"`
		Port          *int   `json:"port"`
		Method        string `json:"method"`
		Location      string `json:"location"`
		CheckInterval int    `json:"check_interval"`
		Status        string `json:"status"`
		Description   string `json:"description"`
	}

	device := DeviceResponse{
		ID:            id,
		Name:          name,
		Type:          typ,
		IP:            ip,
		URL:           url,
		Port:          port,
		Method:        method,
		Location:      location,
		CheckInterval: interval,
		Status:        status,
		Description:   desc,
	}

	log.Printf("Device updated: %d", id)
	respondData(w, device)
}

func (h *DeviceHandler) deleteDevice(w http.ResponseWriter, _ *http.Request, id int) {
	result, err := h.db.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting device: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to delete device")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Device not found")
		return
	}

	h.engine.Stop(id)
	log.Printf("Device deleted: %d", id)
	respondSuccess(w, "Device deleted successfully")
}

func (h *DeviceHandler) startMonitoring(w http.ResponseWriter, _ *http.Request, id int) {
	var ip, method, url, status string
	var port *int
	var interval int
	err := h.db.QueryRow("SELECT ip, method, url, port, check_interval, status FROM devices WHERE id = ?", id).
		Scan(&ip, &method, &url, &port, &interval, &status)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Device not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get device")
		return
	}

	if status != "active" {
		respondError(w, http.StatusBadRequest, "Device is inactive. Activate it first.")
		return
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

	h.engine.Start(config)
	log.Printf("Monitoring started for device %d (%s)", id, ip)
	respondSuccess(w, "Monitoring started")
}

func (h *DeviceHandler) stopMonitoring(w http.ResponseWriter, _ *http.Request, id int) {
	h.engine.Stop(id)
	log.Printf("Monitoring stopped for device %d", id)
	respondSuccess(w, "Monitoring stopped")
}

func (h *DeviceHandler) toggleStatus(w http.ResponseWriter, r *http.Request, id int) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req ToggleStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Status != "active" && req.Status != "inactive" {
		respondError(w, http.StatusBadRequest, "Status must be 'active' or 'inactive'")
		return
	}

	var currentStatus string
	err := h.db.QueryRow("SELECT status FROM devices WHERE id = ?", id).Scan(&currentStatus)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Device not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get device")
		return
	}

	if currentStatus == req.Status {
		respondError(w, http.StatusBadRequest, "Device is already "+req.Status)
		return
	}

	_, err = h.db.Exec("UPDATE devices SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Status, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update device status")
		return
	}

	if req.Status == "active" {
		var ip, method, url string
		var port *int
		var interval int
		h.db.QueryRow("SELECT ip, method, url, port, check_interval FROM devices WHERE id = ?", id).
			Scan(&ip, &method, &url, &port, &interval)

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
		h.engine.Start(config)
		log.Printf("Monitoring started for device %d (%s)", id, ip)
	} else {
		h.engine.Stop(id)
		log.Printf("Monitoring stopped for device %d", id)
	}

	type StatusResponse struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	msg := "Device activated"
	if req.Status == "inactive" {
		msg = "Device deactivated"
	}

	respondData(w, StatusResponse{
		ID:      id,
		Status:  req.Status,
		Message: msg,
	})
}
