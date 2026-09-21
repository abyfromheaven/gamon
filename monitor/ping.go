package monitor

// Module Ping bertanggung jawab menjalankan Perintah Ping ICMP ke alamat IP target
// dan mengukur waktu respon (latensi ms) serta status keberadaan perangkat (online/offline).

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Konstanta status perangkat
const (
	StatusOnline  = "online"  // Perangkat merespon ping dengan baik
	StatusOffline = "offline" // Perangkat tidak merespon / terputus
)

// CheckResult adalah struktur standar hasil pengecekan 1 kali ping.
type CheckResult struct {
	DeviceID  int            `json:"device_id"`  // ID Perangkat
	IP        string         `json:"ip"`         // Alamat IP Perangkat
	Method    string         `json:"method"`     // Metode pengecekan (ICMP Ping)
	Status    string         `json:"status"`     // Status hasil pengecekan (online/offline)
	LatencyMs float64        `json:"latency_ms"` // Latensi / waktu respon dalam milidetik (ms)
	TTL       int            `json:"ttl"`        // Nilai TTL (Time To Live) dari paket ICMP
	Seq       int            `json:"seq"`        // Nomor urut percobaan ping
	Timestamp string         `json:"timestamp"`  // Waktu pengecekan dilakukan
	Details   map[string]any `json:"details"`    // Detail tambahan hasil ping
}

// PingOnce melakukan 1 kali pengiriman paket ICMP Ping ke Alamat IP target.
// Fungsi ini mendeteksi Sistem Operasi (Windows vs Linux) untuk menyesuaikan perintah ping terminal OS.
func PingOnce(ip string, seq int) CheckResult {
	result := CheckResult{
		IP:        ip,
		Status:    StatusOffline,
		Seq:       seq,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   map[string]any{},
	}

	var cmd *exec.Cmd
	// Menyesuaikan perintah CLI ping berdasarkan Sistem Operasi (OS)
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "-w", "3000", ip) // 1 paket, timeout 3000ms
	} else {
		cmd = exec.Command("ping", "-c", "1", "-W", "3", ip)    // 1 paket, timeout 3 detik
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Jika perintah ping error / tidak ada respon, kembalikan status offline
		return result
	}

	// Membaca baris demi baris teks hasil respon ping terminal
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		line = strings.TrimSpace(line)
		
		// Cek apakah baris menandakan balasan ping (Linux atau Windows English/Indonesian)
		isReply := strings.Contains(line, "icmp_seq=") ||
			strings.Contains(line, "Reply from") ||
			strings.Contains(line, "Balasan dari") ||
			(strings.Contains(line, "time=") && strings.Contains(line, "TTL=")) ||
			(strings.Contains(line, "waktu=") && strings.Contains(line, "TTL="))

		if !isReply {
			continue
		}

		// Jika ditemukan tanda balasan, berarti perangkat merespon (online)
		result.Status = StatusOnline
		for _, field := range strings.Fields(line) {
			field = strings.Trim(field, ",")
			if strings.HasPrefix(field, "time<") || strings.HasPrefix(field, "waktu<") {
				result.LatencyMs = 0.5
				continue
			}

			kv := strings.SplitN(field, "=", 2)
			if len(kv) != 2 {
				continue
			}
			
			key := strings.ToLower(kv[0])
			val := strings.TrimSuffix(strings.TrimSuffix(kv[1], "ms"), "s")

			switch key {
			case "ttl":
				// Ambil nilai TTL dari baris respon
				if value, err := strconv.Atoi(val); err == nil {
					result.TTL = value
					result.Details["ttl"] = value
				}
			case "time", "waktu":
				// Ambil nilai waktu latensi (ms) dari baris respon
				if strings.HasPrefix(val, "<") {
					result.LatencyMs = 0.5
				} else if value, err := strconv.ParseFloat(val, 64); err == nil {
					result.LatencyMs = value
				}
			}
		}
	}

	return result
}
