package monitor

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"
)

const offlineFailureThreshold = 3

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

type StatusChange struct {
	DeviceID   int    `json:"device_id"`
	DeviceName string `json:"device_name"`
	OldStatus  string `json:"old_status"`
	NewStatus  string `json:"new_status"`
	Timestamp  string `json:"timestamp"`
}

type CheckFunc func(DeviceConfig, int) CheckResult

type EngineOption func(*Engine)

func WithCheckFunc(check CheckFunc) EngineOption {
	return func(engine *Engine) {
		if check != nil {
			engine.check = check
		}
	}
}

type Engine struct {
	hub HubInterface
	db  *sql.DB

	mu         sync.Mutex
	targets    map[int]context.CancelFunc
	lastStatus map[int]string
	failures   map[int]int
	check      CheckFunc
}

type alertSpec struct {
	title       string
	severity    string
	description string
}

func NewEngine(hub HubInterface, db *sql.DB, options ...EngineOption) *Engine {
	engine := &Engine{
		hub:        hub,
		db:         db,
		targets:    make(map[int]context.CancelFunc),
		lastStatus: make(map[int]string),
		failures:   make(map[int]int),
		check: func(config DeviceConfig, seq int) CheckResult {
			result := PingOnce(config.IP, seq)
			result.DeviceID = config.DeviceID
			result.Method = config.Method
			return result
		},
	}
	for _, option := range options {
		option(engine)
	}
	return engine
}

func (e *Engine) Start(config DeviceConfig) {
	interval := config.Interval
	if interval <= 0 {
		interval = 3
	}

	e.mu.Lock()
	if _, exists := e.targets[config.DeviceID]; exists {
		e.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.targets[config.DeviceID] = cancel
	e.mu.Unlock()

	go e.checkLoop(ctx, config, time.Duration(interval)*time.Second)
	log.Printf("started monitoring device %d (%s)", config.DeviceID, config.IP)
}

func (e *Engine) Stop(deviceID int) {
	e.mu.Lock()
	if cancel, exists := e.targets[deviceID]; exists {
		cancel()
		delete(e.targets, deviceID)
	}
	delete(e.lastStatus, deviceID)
	delete(e.failures, deviceID)
	e.mu.Unlock()
}

func (e *Engine) IsMonitoring(deviceID int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.targets[deviceID]
	return exists
}

func (e *Engine) checkLoop(ctx context.Context, config DeviceConfig, interval time.Duration) {
	seq := 1
	e.runCheck(config, seq)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			seq++
			e.runCheck(config, seq)
		}
	}
}

func (e *Engine) runCheck(config DeviceConfig, seq int) {
	result := e.check(config, seq)
	result.DeviceID = config.DeviceID
	result.IP = config.IP
	result.Method = config.Method
	if result.Timestamp == "" {
		result.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if result.Details == nil {
		result.Details = map[string]any{}
	}

	e.saveCheckResult(result)
	if change, alert, resolveAlert := e.trackStatus(result); change != nil {
		if alert != nil {
			e.createAlert(result.DeviceID, alert.title, alert.severity, alert.description)
		}
		if resolveAlert {
			e.resolveAlerts(result.DeviceID)
		}
		e.hub.Broadcast("status_change", *change)
	}
	e.hub.Broadcast("check_result", result)
}

func (e *Engine) saveCheckResult(result CheckResult) {
	details, err := json.Marshal(result.Details)
	if err != nil {
		details = []byte("{}")
	}
	_, err = e.db.Exec(`INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`, result.DeviceID, result.Status, result.LatencyMs, result.TTL, result.Seq, string(details))
	if err != nil {
		log.Printf("save check result for device %d: %v", result.DeviceID, err)
	}
}

func (e *Engine) trackStatus(result CheckResult) (*StatusChange, *alertSpec, bool) {
	e.mu.Lock()
	oldStatus := e.lastStatus[result.DeviceID]
	newStatus := oldStatus
	var alert *alertSpec
	resolveAlert := false

	switch result.Status {
	case StatusOffline:
		e.failures[result.DeviceID]++
		if e.failures[result.DeviceID] >= offlineFailureThreshold && oldStatus != StatusOffline {
			newStatus = StatusOffline
			alert = &alertSpec{title: "Device Offline", severity: "critical", description: "Device tidak merespons ping dari " + result.IP}
		}
	case StatusWarning:
		e.failures[result.DeviceID] = 0
		newStatus = StatusWarning
		if oldStatus == StatusOffline {
			resolveAlert = true
		}
		if oldStatus == StatusOnline {
			alert = &alertSpec{title: "Device High Latency", severity: "medium", description: "Latency perangkat melebihi 200ms"}
		}
	case StatusOnline:
		e.failures[result.DeviceID] = 0
		newStatus = StatusOnline
		if oldStatus == StatusOffline || oldStatus == StatusWarning {
			resolveAlert = true
		}
	default:
		e.mu.Unlock()
		return nil, nil, false
	}

	if newStatus == oldStatus {
		e.mu.Unlock()
		return nil, nil, false
	}
	e.lastStatus[result.DeviceID] = newStatus
	e.mu.Unlock()

	return &StatusChange{
		DeviceID:   result.DeviceID,
		DeviceName: e.getDeviceName(result.DeviceID),
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Timestamp:  result.Timestamp,
	}, alert, resolveAlert
}

func (e *Engine) createAlert(deviceID int, title, severity, description string) {
	_, err := e.db.Exec(`INSERT INTO alerts (device_id, title, status, severity, description)
		VALUES (?, ?, 'ongoing', ?, ?)`, deviceID, title, severity, description)
	if err != nil {
		log.Printf("create alert for device %d: %v", deviceID, err)
	}
}

func (e *Engine) resolveAlerts(deviceID int) {
	_, err := e.db.Exec(`UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
		WHERE device_id = ? AND status = 'ongoing'`, deviceID)
	if err != nil {
		log.Printf("resolve alerts for device %d: %v", deviceID, err)
	}
}

func (e *Engine) getDeviceName(deviceID int) string {
	var name string
	if err := e.db.QueryRow("SELECT name FROM devices WHERE id = ?", deviceID).Scan(&name); err != nil {
		return "Unknown"
	}
	return name
}
