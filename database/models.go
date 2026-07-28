package database

import "time"

type Device struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	IP             string    `json:"ip"`
	URL            string    `json:"url"`
	Port           *int      `json:"port"`
	Method         string    `json:"method"`
	Location       string    `json:"location"`
	CheckInterval  int       `json:"check_interval"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PingHistory struct {
	ID        int       `json:"id"`
	DeviceID  int       `json:"device_id"`
	Status    string    `json:"status"`
	LatencyMs float64   `json:"latency_ms"`
	TTL       int       `json:"ttl"`
	Seq       int       `json:"seq"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

type Alert struct {
	ID          int        `json:"id"`
	DeviceID    int        `json:"device_id"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	Severity    string     `json:"severity"`
	StartedAt   time.Time  `json:"started_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	Description string     `json:"description"`
}

type DeviceWithType struct {
	Device
	TypeName string `json:"type_name"`
}

type MonitoringStatus struct {
	DeviceID    int       `json:"device_id"`
	DeviceName  string    `json:"device_name"`
	DeviceType  string    `json:"device_type"`
	IP          string    `json:"ip"`
	Method      string    `json:"method"`
	Status      string    `json:"status"`
	LatencyMs   float64   `json:"latency_ms"`
	LastCheck   time.Time `json:"last_check"`
	Interval    int       `json:"interval"`
}

type DashboardSummary struct {
	TotalDevices  int    `json:"total_devices"`
	OnlineDevices int    `json:"online_devices"`
	OfflineDevices int   `json:"offline_devices"`
	WarningDevices int   `json:"warning_devices"`
	LatestAlerts  []Alert `json:"latest_alerts"`
}
