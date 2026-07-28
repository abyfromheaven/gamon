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

type Result struct {
	DeviceID  int     `json:"device_id"`
	IP        string  `json:"ip"`
	Method    string  `json:"method"`
	Status    string  `json:"status"`
	Latency   float64 `json:"latency"`
	TTL       int     `json:"ttl"`
	Seq       int     `json:"seq"`
	Timestamp string  `json:"timestamp"`
}

func PingOnce(ip string, seq int) Result {
	result := Result{
		IP:        ip,
		Status:    StatusOffline,
		Latency:   0,
		TTL:       0,
		Seq:       seq,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "3", ip)
	}

	out, err := cmd.CombinedOutput()
	raw := strings.TrimSpace(string(out))

	if err != nil {
		return result
	}

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "icmp_seq=") {
			continue
		}

		result.Status = StatusOnline

		for _, field := range strings.Fields(line) {
			if kv := strings.SplitN(field, "=", 2); len(kv) == 2 {
				switch kv[0] {
				case "ttl":
					if v, err := strconv.Atoi(kv[1]); err == nil {
						result.TTL = v
					}
				case "time":
					if v, err := strconv.ParseFloat(kv[1], 64); err == nil {
						result.Latency = v
					}
				}
			}
		}
	}

	if result.Status == StatusOnline && result.Latency >= LatencyWarningThreshold {
		result.Status = StatusWarning
	}

	return result
}
