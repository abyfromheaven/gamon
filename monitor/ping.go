package monitor

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	StatusOnline  = "online"
	StatusWarning = "warning"
	StatusOffline = "offline"

	LatencyWarningThreshold = 200.0
)

// CheckResult is the normalized result emitted by every monitoring method.
type CheckResult struct {
	DeviceID  int            `json:"device_id"`
	IP        string         `json:"ip"`
	Method    string         `json:"method"`
	Status    string         `json:"status"`
	LatencyMs float64        `json:"latency_ms"`
	TTL       int            `json:"ttl"`
	Seq       int            `json:"seq"`
	Timestamp string         `json:"timestamp"`
	Details   map[string]any `json:"details"`
}

// PingOnce executes one ICMP ping and normalizes its output into CheckResult.
func PingOnce(ip string, seq int) CheckResult {
	result := CheckResult{
		IP:        ip,
		Status:    StatusOffline,
		Seq:       seq,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   map[string]any{},
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "3", ip)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return result
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "icmp_seq=") {
			continue
		}

		result.Status = StatusOnline
		for _, field := range strings.Fields(line) {
			kv := strings.SplitN(field, "=", 2)
			if len(kv) != 2 {
				continue
			}
			switch kv[0] {
			case "ttl":
				if value, err := strconv.Atoi(kv[1]); err == nil {
					result.TTL = value
					result.Details["ttl"] = value
				}
			case "time":
				if value, err := strconv.ParseFloat(kv[1], 64); err == nil {
					result.LatencyMs = value
				}
			}
		}
	}

	if result.Status == StatusOnline && result.LatencyMs >= LatencyWarningThreshold {
		result.Status = StatusWarning
	}

	return result
}
