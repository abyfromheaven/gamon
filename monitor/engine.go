package monitor

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"
)

type HubInterface interface {
	Broadcast(msgType string, data interface{})
}

type DeviceConfig struct {
	DeviceID int
	IP       string
	URL      string
	Port     int
	Method   string
	Interval int
}

type DeviceStatus struct {
	DeviceID  int     `json:"device_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	IP        string  `json:"ip"`
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	LastCheck string  `json:"last_check"`
}

type Engine struct {
	hub        HubInterface
	db         *sql.DB
	targets    map[int]context.CancelFunc
	lastStatus map[int]string
	mu         sync.Mutex
}

func NewEngine(hub HubInterface, db *sql.DB) *Engine {
	return &Engine{
		hub:        hub,
		db:         db,
		targets:    make(map[int]context.CancelFunc),
		lastStatus: make(map[int]string),
	}
}

func (e *Engine) Start(config DeviceConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.targets[config.DeviceID]; exists {
		log.Printf("Already monitoring device %d", config.DeviceID)
		return
	}

	interval := config.Interval
	if interval <= 0 {
		interval = 3
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.targets[config.DeviceID] = cancel

	go e.checkLoop(ctx, config, time.Duration(interval)*time.Second)

	log.Printf("Started monitoring device %d (%s)", config.DeviceID, config.IP)
}

func (e *Engine) Stop(deviceID int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, exists := e.targets[deviceID]; exists {
		cancel()
		delete(e.targets, deviceID)
		delete(e.lastStatus, deviceID)
		log.Printf("Stopped monitoring device %d", deviceID)
	}
}

func (e *Engine) IsMonitoring(deviceID int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.targets[deviceID]
	return exists
}

func (e *Engine) GetActiveDeviceStatuses() []DeviceStatus {
	e.mu.Lock()
	defer e.mu.Unlock()

	var statuses []DeviceStatus

	rows, err := e.db.Query(`
		SELECT d.id, d.name, d.type, d.ip,
			COALESCE(ph.status, 'unknown') as last_status,
			COALESCE(ph.latency_ms, 0) as last_latency,
			COALESCE(ph.timestamp, d.created_at) as last_check
		FROM devices d
		LEFT JOIN ping_history ph ON ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY id DESC LIMIT 1
		)
		WHERE d.status = 'active'
		ORDER BY d.name
	`)
	if err != nil {
		log.Printf("Error querying device statuses: %v", err)
		return statuses
	}
	defer rows.Close()

	for rows.Next() {
		var ds DeviceStatus
		var lastCheck string
		if err := rows.Scan(&ds.DeviceID, &ds.Name, &ds.Type, &ds.IP, &ds.Status, &ds.LatencyMs, &lastCheck); err != nil {
			log.Printf("Error scanning device status: %v", err)
			continue
		}
		ds.LastCheck = lastCheck
		statuses = append(statuses, ds)
	}

	return statuses
}

func (e *Engine) checkLoop(ctx context.Context, config DeviceConfig, interval time.Duration) {
	seq := 0
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	seq++
	result := PingOnce(config.IP, seq)
	result.DeviceID = config.DeviceID
	result.Method = config.Method

	e.savePingResult(result)
	e.checkStatusChange(config.DeviceID, config.IP)
	e.hub.Broadcast("ping_result", result)

	log.Printf("[device=%d] seq=%d status=%s latency=%.2fms", config.DeviceID, seq, result.Status, result.Latency)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Check loop stopped for device %d", config.DeviceID)
			return
		case <-ticker.C:
			seq++
			result := PingOnce(config.IP, seq)
			result.DeviceID = config.DeviceID
			result.Method = config.Method

			e.savePingResult(result)
			e.checkStatusChange(config.DeviceID, config.IP)
			e.hub.Broadcast("ping_result", result)

			log.Printf("[device=%d] seq=%d status=%s latency=%.2fms", config.DeviceID, seq, result.Status, result.Latency)
		}
	}
}

func (e *Engine) savePingResult(result Result) {
	_, err := e.db.Exec(
		`INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp) 
		 VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		result.DeviceID, result.Status, result.Latency, result.TTL, result.Seq, "{}",
	)
	if err != nil {
		log.Printf("Failed to save ping result for device %d: %v", result.DeviceID, err)
	}
}

func (e *Engine) checkStatusChange(deviceID int, ip string) {
	var currentStatus string
	err := e.db.QueryRow("SELECT status FROM ping_history WHERE device_id = ? ORDER BY id DESC LIMIT 1", deviceID).Scan(&currentStatus)
	if err != nil {
		e.lastStatus[deviceID] = "unknown"
		return
	}

	oldStatus, exists := e.lastStatus[deviceID]

	if exists && oldStatus != currentStatus {
		e.handleStatusChange(deviceID, ip, oldStatus, currentStatus)
	}

	e.lastStatus[deviceID] = currentStatus
}

func (e *Engine) handleStatusChange(deviceID int, ip, oldStatus, newStatus string) {
	log.Printf("Status changed for device %d: %s -> %s", deviceID, oldStatus, newStatus)

	if oldStatus == "online" && newStatus == "offline" {
		e.createAlert(deviceID, "Device Offline", "high", "Device tidak merespons ping dari "+ip)
	}

	if oldStatus == "online" && newStatus == "warning" {
		e.createAlert(deviceID, "Device High Latency", "medium", "Latency lebih dari 200ms")
	}

	if newStatus == "online" && (oldStatus == "offline" || oldStatus == "warning" || oldStatus == "unknown") {
		e.resolveAlert(deviceID)
	}
}

func (e *Engine) createAlert(deviceID int, title, severity, description string) {
	_, err := e.db.Exec(
		`INSERT INTO alerts (device_id, title, status, severity, description) 
		 VALUES (?, ?, 'ongoing', ?, ?)`,
		deviceID, title, severity, description,
	)
	if err != nil {
		log.Printf("Failed to create alert for device %d: %v", deviceID, err)
		return
	}

	alertData := map[string]interface{}{
		"device_id": deviceID,
		"title":     title,
		"severity":  severity,
		"description": description,
	}
	e.hub.Broadcast("alert_created", alertData)

	log.Printf("Alert created: device %d - %s (%s)", deviceID, title, severity)
}

func (e *Engine) resolveAlert(deviceID int) {
	_, err := e.db.Exec(
		`UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP 
		 WHERE device_id = ? AND status = 'ongoing'`,
		deviceID,
	)
	if err != nil {
		log.Printf("Failed to resolve alert for device %d: %v", deviceID, err)
		return
	}

	alertData := map[string]interface{}{
		"device_id": deviceID,
	}
	e.hub.Broadcast("alert_resolved", alertData)

	log.Printf("Alert resolved: device %d", deviceID)
}

func (e *Engine) getDeviceName(deviceID int) string {
	var name string
	err := e.db.QueryRow("SELECT name FROM devices WHERE id = ?", deviceID).Scan(&name)
	if err != nil {
		return "Unknown"
	}
	return name
}

func (e *Engine) marshalAndSend(msgType string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}
	e.hub.Broadcast(msgType, string(dataBytes))
}
