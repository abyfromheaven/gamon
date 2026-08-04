package monitor

import (
	"database/sql"
	"sync"
	"testing"

	_ "modernc.org/sqlite"
)

type recordingHub struct {
	mu       sync.Mutex
	messages []string
}

func (h *recordingHub) Broadcast(messageType string, _ interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = append(h.messages, messageType)
}

func TestEngine_ThreeOfflineChecksCreateAlertAndRecoveryResolvesIt(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for _, query := range []string{
		`CREATE TABLE devices (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`,
		`CREATE TABLE ping_history (device_id INTEGER, status TEXT, latency_ms REAL, ttl INTEGER, seq INTEGER, details TEXT, timestamp DATETIME)`,
		`CREATE TABLE alerts (device_id INTEGER, title TEXT, status TEXT, severity TEXT, description TEXT, resolved_at DATETIME, acknowledged BOOLEAN DEFAULT FALSE, acknowledged_at DATETIME)`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(`INSERT INTO devices (id, name) VALUES (1, 'Router-01')`); err != nil {
		t.Fatal(err)
	}

	statuses := []string{StatusOffline, StatusOffline, StatusOffline, StatusOnline, StatusOnline, StatusOnline}
	index := 0
	hub := &recordingHub{}
	engine := NewEngine(hub, db, WithCheckFunc(func(_ DeviceConfig, seq int) CheckResult {
		status := statuses[index]
		index++
		return CheckResult{Status: status, Seq: seq, Timestamp: "2026-07-30T00:00:00Z"}
	}))
	config := DeviceConfig{DeviceID: 1, IP: "192.168.1.1", Method: "ICMP Ping"}

	// 3x offline → should create 1 critical alert
	for seq := 1; seq <= 3; seq++ {
		engine.runCheck(config, seq)
	}
	var ongoing int
	if err := db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE status = 'ongoing' AND severity = 'critical'`).Scan(&ongoing); err != nil {
		t.Fatal(err)
	}
	if ongoing != 1 {
		t.Fatalf("ongoing critical alerts = %d, want 1", ongoing)
	}

	// 1x online → should auto-resolve the alert (no recovery alert created)
	engine.runCheck(config, 4)
	var resolved int
	if err := db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE status = 'resolved'`).Scan(&resolved); err != nil {
		t.Fatal(err)
	}
	if resolved != 1 {
		t.Fatalf("resolved alerts = %d, want 1", resolved)
	}

	// No recovery alerts should exist
	var totalAlerts int
	if err := db.QueryRow(`SELECT COUNT(*) FROM alerts`).Scan(&totalAlerts); err != nil {
		t.Fatal(err)
	}
	if totalAlerts != 1 {
		t.Fatalf("total alerts = %d, want 1 (only critical, no recovery)", totalAlerts)
	}

	// Check WebSocket messages: 2 status_changes + 1 alert_created (critical only)
	hub.mu.Lock()
	defer hub.mu.Unlock()
	statusChanges := 0
	for _, message := range hub.messages {
		if message == "status_change" {
			statusChanges++
		}
	}
	// Note: alert_created is no longer broadcast via WS in simplified system
	// We only have status_change events now
	if statusChanges != 2 {
		t.Fatalf("status_changes = %d, want 2", statusChanges)
	}
	if hub.messages[len(hub.messages)-1] != "check_result" {
		t.Fatalf("last message = %s, want check_result", hub.messages[len(hub.messages)-1])
	}
}
