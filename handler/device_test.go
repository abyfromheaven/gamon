package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gamon/monitor"

	_ "modernc.org/sqlite"
)

type mockHub struct{}

func (m *mockHub) Broadcast(_ string, _ interface{}) {}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	for _, query := range []string{
		`CREATE TABLE devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			ip TEXT NOT NULL,
			url TEXT DEFAULT '',
			port INTEGER,
			method TEXT NOT NULL DEFAULT 'ICMP Ping',
			location TEXT DEFAULT '',
			check_interval INTEGER DEFAULT 3,
			status TEXT DEFAULT 'active',
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE ping_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			latency_ms REAL DEFAULT 0,
			ttl INTEGER DEFAULT 0,
			seq INTEGER DEFAULT 0,
			details TEXT DEFAULT '{}',
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT DEFAULT 'ongoing',
			severity TEXT DEFAULT 'low',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			resolved_at DATETIME,
			description TEXT DEFAULT ''
		)`,
	} {
		if _, err := db.Exec(query); err != nil {
			t.Fatal(err)
		}
	}
	return db
}

func TestStartStop_MustBePOST(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	engine := monitor.NewEngine(&mockHub{}, db)
	hub := NewHub(db)
	h := NewDeviceHandler(db, engine, hub)

	// Insert a test device
	_, err := db.Exec(`INSERT INTO devices (name, type, ip, method, status) VALUES ('R1', 'Router', '10.0.0.1', 'ICMP Ping', 'active')`)
	if err != nil {
		t.Fatal(err)
	}

	// GET on /start should return 405
	req := httptest.NewRequest(http.MethodGet, "/api/devices/1/start", nil)
	w := httptest.NewRecorder()
	h.HandleDevice(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /start: got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}

	// DELETE on /stop should return 405
	req = httptest.NewRequest(http.MethodDelete, "/api/devices/1/stop", nil)
	w = httptest.NewRecorder()
	h.HandleDevice(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("DELETE /stop: got %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestCreateActiveDevice_StartsMonitoring(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	started := make(chan int, 1)
	engine := monitor.NewEngine(&mockHub{}, db, monitor.WithCheckFunc(func(config monitor.DeviceConfig, _ int) monitor.CheckResult {
		started <- config.DeviceID
		return monitor.CheckResult{Status: monitor.StatusOnline, LatencyMs: 1, Seq: 1}
	}))
	hub := NewHub(db)
	h := NewDeviceHandler(db, engine, hub)

	body, _ := json.Marshal(CreateDeviceRequest{
		Name:          "Server-01",
		Type:          "Server",
		IP:            "192.168.1.10",
		Method:        "ICMP Ping",
		CheckInterval: 3,
		Status:        "active",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/devices", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.HandleDevices(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("create active device: got %d, want 200, body: %s", w.Code, w.Body.String())
	}

	var resp DataResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(resp.Data)
	var created struct {
		ID int `json:"id"`
	}
	json.Unmarshal(data, &created)

	if created.ID == 0 {
		t.Fatal("expected device ID > 0")
	}

	// Engine should have started monitoring for this device
	if !engine.IsMonitoring(created.ID) {
		t.Fatal("expected engine to be monitoring the newly created active device")
	}
}

func TestUpdateInactive_StopsMonitoring(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	engine := monitor.NewEngine(&mockHub{}, db)
	hub := NewHub(db)
	h := NewDeviceHandler(db, engine, hub)

	// Insert active device
	res, err := db.Exec(`INSERT INTO devices (name, type, ip, method, status, check_interval) VALUES ('SW1', 'Switch', '10.0.0.2', 'ICMP Ping', 'active', 3)`)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()

	// Start monitoring
	engine.Start(monitor.DeviceConfig{DeviceID: int(id), IP: "10.0.0.2", Method: "ICMP Ping", Interval: 3})
	if !engine.IsMonitoring(int(id)) {
		t.Fatal("expected monitoring to be active before update")
	}

	// Update to inactive
	status := "inactive"
	body, _ := json.Marshal(UpdateDeviceRequest{Status: &status})
	req := httptest.NewRequest(http.MethodPut, "/api/devices/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.HandleDevice(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("update to inactive: got %d, want 200, body: %s", w.Code, w.Body.String())
	}

	if engine.IsMonitoring(int(id)) {
		t.Fatal("expected monitoring to be stopped after setting device inactive")
	}
}

func TestDelete_StopsMonitoring(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	engine := monitor.NewEngine(&mockHub{}, db)
	hub := NewHub(db)
	h := NewDeviceHandler(db, engine, hub)

	// Insert active device
	res, err := db.Exec(`INSERT INTO devices (name, type, ip, method, status, check_interval) VALUES ('AP1', 'Access Point', '10.0.0.3', 'ICMP Ping', 'active', 3)`)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()

	// Start monitoring
	engine.Start(monitor.DeviceConfig{DeviceID: int(id), IP: "10.0.0.3", Method: "ICMP Ping", Interval: 3})
	if !engine.IsMonitoring(int(id)) {
		t.Fatal("expected monitoring to be active before delete")
	}

	// Delete device
	req := httptest.NewRequest(http.MethodDelete, "/api/devices/1", nil)
	w := httptest.NewRecorder()
	h.HandleDevice(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("delete device: got %d, want 200, body: %s", w.Code, w.Body.String())
	}

	if engine.IsMonitoring(int(id)) {
		t.Fatal("expected monitoring to be stopped after deleting device")
	}
}
