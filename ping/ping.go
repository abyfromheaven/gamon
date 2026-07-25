package ping

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	IP      string
	Alive   bool
	Seq     int
	TTL     string
	Time    string
	RawLine string
}

func PingOnce(ip string) Result {
	result := Result{IP: ip}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", ip)
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "3", ip)
	}

	out, err := cmd.CombinedOutput()
	raw := strings.TrimSpace(string(out))
	result.RawLine = raw

	if err != nil {
		result.Alive = false
		return result
	}

	result.Alive = true

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "icmp_seq=") {
			for _, field := range strings.Fields(line) {
				if kv := strings.SplitN(field, "=", 2); len(kv) == 2 {
					switch kv[0] {
					case "icmp_seq":
						if v, err := strconv.Atoi(kv[1]); err == nil {
							result.Seq = v
						}
					case "ttl":
						result.TTL = kv[1]
					case "time":
						result.Time = kv[1]
					}
				}
			}
		}
	}

	return result
}

func FormatOutput(r Result) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	if r.Alive {
		return fmt.Sprintf("[%s] %s -> UP   | seq=%-3d ttl=%-3s time=%sms", ts, r.IP, r.Seq, r.TTL, r.Time)
	}
	return fmt.Sprintf("[%s] %s -> DOWN | seq=%-3d no reply (timeout)", ts, r.IP, r.Seq)
}
