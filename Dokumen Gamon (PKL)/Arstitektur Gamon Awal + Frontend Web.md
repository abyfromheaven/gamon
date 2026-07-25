Oke, ini arsitektur dan konsep lengkapnya:
GAMON Web v0.2 — Arsitektur & Konsep Sistem
1. High-Level Architecture
┌─────────────────────────────────────────────────────────┐
│                    BROWSER (Client)                     │
│                  React + TailwindCSS                    │
│                     Port :5173                          │
│                                                         │
│  ┌──────────┐   ┌──────────────┐   ┌────────────────┐  │
│  │ Input IP │──▶│ WebSocket    │──▶│ Dashboard Panel│  │
│  │ Component│   │ Client Hook  │   │ (Status/Chart) │  │
│  └──────────┘   └──────┬───────┘   └────────────────┘  │
│                        │                                │
└────────────────────────┼────────────────────────────────┘
                         │ WebSocket (ws://localhost:8080/ws)
                         │
┌────────────────────────┼────────────────────────────────┐
│                   GO BACKEND                            │
│                   Port :8080                             │
│                                                         │
│  ┌──────────┐   ┌──────┴───────┐   ┌────────────────┐  │
│  │ REST API │   │ WebSocket    │   │ Ping Engine    │  │
│  │ /api/*   │   │ Hub          │   │ (Goroutine)    │  │
│  │ (CORS)   │   │ (Broadcast)  │   │                │  │
│  └──────────┘   └──────────────┘   └────────────────┘  │
│                                                         │
│  ┌──────────────────────────────────────────────────┐   │
│  │              OS Ping Command                     │   │
│  │         (via os/exec — ICMP Echo Request)        │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
2. Alur Komunikasi (Data Flow)
USER INPUT IP
     │
     ▼
[Browser] ──POST /api/monitor──▶ [Go Backend]
                                    │
                                    ▼
                              Simpan target IP
                              di memory (map)
                                    │
                                    ▼
                              Spawn Goroutine
                              tiap target IP
                                    │
                          ┌─────────┴─────────┐
                          │   Loop Setiap 3s   │
                          │   exec ping -c 1   │
                          └─────────┬─────────┘
                                    │
                                    ▼
                              Parse stdout
                              Hitung latency
                                    │
                                    ▼
                              Marshal to JSON
                                    │
                                    ▼
                         ┌──────────┴──────────┐
                         │   WebSocket Hub      │
                         │   broadcast(result)  │
                         └──────────┬──────────┘
                                    │
                                    ▼
                         [Browser React via WS]
                                    │
                                    ▼
                         Update Dashboard
                         Realtime (no refresh)
3. Struktur File yang Akan Dibuat
/home/aby/gamon/
├── main.go                    # Go entry point, HTTP + WS server
├── handler/
│   ├── websocket.go           # WebSocket Hub (upgrade, broadcast, clients)
│   └── api.go                 # REST API endpoints (POST /api/monitor)
├── monitor/
│   ├── ping.go                # ICMP ping logic (refactor dari CLI)
│   └── engine.go              # Goroutine manager (start/stop per target)
├── go.mod
├── go.sum
│
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── index.html
│   └── src/
│       ├── main.tsx
│       ├── App.tsx
│       ├── components/
│       │   ├── Dashboard.tsx    # Panel utama monitoring
│       │   ├── StatusCard.tsx   # Card status 1 device (UP/DOWN + latency)
│       │   └── InputForm.tsx    # Form input IP
│       ├── hooks/
│       │   └── useWebSocket.ts  # Custom hook koneksi WebSocket
│       └── types/
│           └── index.ts         # TypeScript type definitions
4. Protokol Data (WebSocket JSON)
Backend → Frontend (setiap hasil ping):
{
  "type": "ping_result",
  "data": {
    "ip": "192.168.1.1",
    "status": "online",
    "latency": 12.4,
    "ttl": 64,
    "seq": 42,
    "timestamp": "2026-07-25T16:01:15Z"
  }
}
Backend → Frontend (saat status berubah):
{
  "type": "status_change",
  "data": {
    "ip": "192.168.1.1",
    "old_status": "online",
    "new_status": "down",
    "timestamp": "2026-07-25T16:02:00Z"
  }
}
Frontend → Backend (start monitoring):
{
  "action": "start",
  "ip": "192.168.1.1"
}
Frontend → Backend (stop monitoring):
{
  "action": "stop",
  "ip": "192.168.1.1"
}
5. Go Backend — Konsep Tiap Komponen
main.go
- Inisialisasi HTTP server di :8080
- Register route: /ws (WebSocket upgrade), /api/* (REST)
- Aktifkan CORS untuk dev React di :5173
- Jalankan server
handler/websocket.go — WebSocket Hub
- Struct Hub: managing connected clients, broadcast channel
- HandleWebSocket(): upgrade HTTP → WebSocket, register client
- readPump(): listen pesan dari client (start/stop action)
- writePump(): kirim data ke client dari broadcast channel
- Pattern: Fan-out broadcast — 1 pesan dari backend → semua client terima
handler/api.go — REST Endpoints
- POST /api/monitor — terima IP dari frontend, validasi, forward ke engine
- CORS middleware untuk dev environment
monitor/ping.go — ICMP Ping
- Refactor dari CLI, tapi output-nya bukan string console tapi struct JSON
- PingOnce(ip string) Result → return JSON-ready struct
- Parsing sama: exec ping -c 1 -W 3 <ip>, parse stdout
monitor/engine.go — Goroutine Manager
- Struct Engine: simpan map target → goroutine handle
- StartMonitor(ip, wsHub): spawn goroutine yang loop ping tiap 3 detik
- StopMonitor(ip): kirim signal stop ke goroutine, cleanup
- Concurrency: pakai context.Context per goroutine untuk graceful stop
- Setiap hasil ping → marshal JSON → kirim ke Hub broadcast channel
6. Frontend React — Konsep Tiap Komponen
App.tsx
- Render InputForm + Dashboard
- Manage state: list of monitored IPs
InputForm.tsx
- Text input + button "Monitor"
- Validate IP format (regex)
- Kirim IP ke backend via REST POST /api/monitor
- Trigger WebSocket start action
Dashboard.tsx
- Render list StatusCard berdasarkan data dari WebSocket
- Auto-update tanpa refresh
StatusCard.tsx
- Tampilkan: IP, Status (UP/DOWN), Latency (ms), TTL, Seq, Timestamp
- Visual indicator: hijau = Online, merah = Down
- Update realtime via props dari parent
hooks/useWebSocket.ts
- Custom hook: buka koneksi WebSocket ke ws://localhost:8080/ws
- Listen message, parse JSON, dispatch ke state
- Handle reconnect jika koneksi putus
- Return: { data, isConnected }
7. Sequence Diagram — Satu Siklus Monitoring
Browser          Go Backend          Goroutine         OS Ping
   │                 │                   │                │
   │── start(ip) ──▶│                   │                │
   │                 │── spawn() ──────▶│                │
   │                 │                   │── ping -c 1 ─▶│
   │                 │                   │◀── reply ─────│
   │                 │                   │                │
   │                 │◀── result JSON ───│                │
   │◀── WS push ────│                   │                │
   │                 │                   │   (sleep 3s)   │
   │                 │                   │── ping -c 1 ─▶│
   │                 │                   │◀── reply ─────│
   │                 │◀── result JSON ───│                │
   │◀── WS push ────│                   │                │
   │                 │                   │                │
8. Pertimbangan Teknis
Aspek	Keputusan
WS Library (Go)	gorilla/webnet atau nhooyr.io/websocket — tapi untuk zero-dependency, bisa pakai standar golang.org/x/net/websocket
CORS	Manual header di Go (Access-Control-Allow-Origin: *), tanpa library
WS Library (React)	Vanilla WebSocket API native browser — tanpa library, simpel
Ping Interval	3 detik default, cukup ringan untuk prototype
State Management	React useState/useReducer — belum perlu Redux untuk prototype
TailwindCSS	Via CDN atau install di Vite project
9. Risk & Mitigation
Risk	Mitigasi
Root permission untuk raw ICMP	Pakai os/exec ke system ping (sudah resolved di CLI)
Goroutine leak	context.Context + proper channel cleanup saat stop
WS reconnect	Frontend auto-reconnect dengan exponential backoff
CORS issue	Explicit allow localhost:5173 di Go backend
10. Urutan Pengerjaan
11. Backend Go — Setup HTTP server + WebSocket hub (kosong, hanya terima koneksi)
12. Backend Go — Refactor ping.go + buat engine.go (goroutine manager)
13. Backend Go — Hubungkan: engine → WebSocket broadcast
14. Frontend — Scaffold Vite + React + TailwindCSS
15. Frontend — Build useWebSocket hook
16. Frontend — Build InputForm, Dashboard, StatusCard
17. Integration test — Full flow: input IP → WS → dashboard update