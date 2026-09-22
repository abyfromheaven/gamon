package monitor

// ============================================================================
// MODUL EXEKUSI ICMP PING (monitor/ping.go)
// ============================================================================
// Modul ini bertanggung jawab menjalankan perintah CLI `ping` bawaan Sistem
// Operasi (Windows vs Linux) ke alamat IP target.
//
// Alur Kerja:
// 1. Memeriksa OS (runtime.GOOS).
// 2. Mengeksekusi perintah CLI ping (`ping -n 1` di Windows / `ping -c 1` di Linux).
// 3. Membaca teks keluaran terminal (stdout) dan mengekstrak latensi (ms) serta TTL.
// 4. Mengembalikan struct CheckResult berstatus 'online' atau 'offline'.
// ============================================================================

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
	LatencyMs float64        `json:"latency_ms"` // Waktu respon dalam milidetik (ms)
	TTL       int            `json:"ttl"`        // Time To Live dari paket ICMP
	Seq       int            `json:"seq"`        // Nomor urut percobaan ping
	Timestamp string         `json:"timestamp"`  // Waktu pengecekan dilakukan
	Details   map[string]any `json:"details"`    // Detail tambahan hasil ping
}

// PingOnce melakukan 1 kali pengiriman paket ICMP Ping ke Alamat IP target.
func PingOnce(ip string, urutan int) CheckResult {
	hasil := CheckResult{
		IP:        ip,
		Status:    StatusOffline,
		Seq:       urutan,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   map[string]any{},
	}

	var perintah *exec.Cmd
	// Menyesuaikan perintah CLI ping berdasarkan Sistem Operasi (OS)
	if runtime.GOOS == "windows" {
		perintah = exec.Command("ping", "-n", "1", "-w", "3000", ip) // 1 paket, timeout 3000ms
	} else {
		perintah = exec.Command("ping", "-c", "1", "-W", "3", ip)    // 1 paket, timeout 3 detik
	}

	keluaran, err := perintah.CombinedOutput()
	if err != nil {
		// Jika perintah ping error / tidak ada respon, kembalikan status offline
		return hasil
	}

	// Membaca baris demi baris teks hasil respon ping terminal
	for _, barisTeks := range strings.Split(strings.TrimSpace(string(keluaran)), "\n") {
		barisTeks = strings.TrimSpace(barisTeks)

		if runtime.GOOS == "windows" {
			// Output Windows ping: "Reply from 192.168.1.1: bytes=32 time=2ms TTL=64"
			if !strings.Contains(barisTeks, "Reply from") {
				continue
			}

			hasil.Status = StatusOnline
			for _, elemen := range strings.Fields(barisTeks) {
				pasangan := strings.SplitN(elemen, "=", 2)
				if len(pasangan) != 2 {
					// Format "time<1ms"
					if idx := strings.Index(elemen, "<"); idx > 0 {
						nilai := elemen[idx+1:]
						nilai = strings.TrimSuffix(nilai, "ms")
						if f, err := strconv.ParseFloat(nilai, 64); err == nil {
							hasil.LatencyMs = f
						}
					}
					continue
				}
				switch strings.ToLower(pasangan[0]) {
				case "ttl":
					if nilai, err := strconv.Atoi(pasangan[1]); err == nil {
						hasil.TTL = nilai
						hasil.Details["ttl"] = nilai
					}
				case "time":
					nilai := strings.TrimSuffix(pasangan[1], "ms")
					if f, err := strconv.ParseFloat(nilai, 64); err == nil {
						hasil.LatencyMs = f
					}
				}
			}
		} else {
			// Output Linux ping: "64 bytes from 192.168.1.1: icmp_seq=1 ttl=64 time=1.23 ms"
			if !strings.Contains(barisTeks, "icmp_seq=") {
				continue
			}

			hasil.Status = StatusOnline
			for _, elemen := range strings.Fields(barisTeks) {
				pasangan := strings.SplitN(elemen, "=", 2)
				if len(pasangan) != 2 {
					continue
				}
				switch pasangan[0] {
				case "ttl":
					if nilai, err := strconv.Atoi(pasangan[1]); err == nil {
						hasil.TTL = nilai
						hasil.Details["ttl"] = nilai
					}
				case "time":
					if nilai, err := strconv.ParseFloat(pasangan[1], 64); err == nil {
						hasil.LatencyMs = nilai
					}
				}
			}
		}
	}

	return hasil
}
