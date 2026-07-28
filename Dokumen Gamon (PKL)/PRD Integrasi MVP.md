# PRD: Integrasi MVP Gamon (Full-Stack Working System)

## Overview

PRD ini mendeskripsikan rancangan integrasi MVP Gamon dari kondisi **frontend dummy + backend scaffold** menjadi **full-stack working system** di mana semua halaman menggunakan data real dari database, monitoring berjalan secara real-time, dan alert dihasilkan secara otomatis.

---

## Kondisi Saat Ini

| Komponen | Status |
|---|---|
| Frontend | UI lengkap (Dashboard, Device Management, Alert Center) tapi pakai data dummy/hardcoded |
| Backend | Engine ping + WebSocket jalan, tapi gak nyambung ke frontend |
| Database | Gak ada — semua data in-memory, hilang saat restart |
| Device CRUD | Frontend-only (local React state), gak persist ke backend |
| Alert | Dummy data, gak ada auto-generation dari status change |
| Monitoring method | Hardcoded ICMP Ping |

### Masalah Utama

1. Frontend **tidak pernah** memanggil REST API backend
2. Frontend **tidak pernah** menampilkan data live dari WebSocket
3. Backend **tidak memiliki** konsep "device" — cuma tahu IP address
4. Tidak ada persistence layer (database)
5. Tidak ada alert generation logic
6. WebSocket `readPump` mengabaikan pesan dari client

---

## Target MVP

Full-stack working system di mana:

- Semua halaman menggunakan data real dari database
- Device Management bisa CRUD beneran (simpan ke SQLite)
- Monitoring real-time (ping jalan, status update live via WebSocket)
- Alert otomatis terbuat saat device down/recover
- 4 halaman berfungsi: Dashboard, Device Management, Monitoring, Alert Center
- Sistem **fleksibel** untuk multi-method monitoring (ICMP Ping sekarang, HTTP/TCP/SNMP nanti)

### Prinsip Desain

> Sistem harus **tidak terpaku** pada satu metode monitoring. Struktur data, engine, dan frontend harus dirancang agar metode baru bisa ditambahkan tanpa mengubah fondasi.

---

## Tech Stack

| Komponen | Pilihan | Alasan |
|---|---|---|
| Database | SQLite | File-based, zero config, cocok untuk server NOC tunggal |
| Driver Go | `modernc.org/sqlite` | Pure Go, CGO-free, gak butuh gcc saat build |
| ORM | Tidak pakai | Raw SQL + `database/sql` biar simpel dan kontrol penuh |
| Backend | Go (existing) | Sudah terbangun, goroutines untuk concurrency |
| Frontend | React + Vite + TailwindCSS (existing) | Sudah terbangun |
| WebSocket | gorilla/websocket (existing) | Sudah terpakai |

---

## Arsitektur Target

```
┌──────────────┐   REST API (CRUD)  ┌─────────────────────┐
│   Frontend   │ ◄────────────────► │   Go Backend         │
│   (React)    │                    │                       │
│              │   WebSocket        │  ┌─────────────────┐  │
│              │ ◄════════════════► │  │  Hub (WS)       │  │
└──────────────┘                    │  └────────┬────────┘  │
                                    │           │           │
                                    │  ┌────────▼────────┐  │
                                    │  │  Engine          │  │
                                    │  │  ┌──────────┐   │  │
                                    │  │  │ ICMP     │   │  │
                                    │  │  │ HTTP     │   │  │
                                    │  │  │ TCP Port │   │  │
                                    │  │  │ SNMP     │   │  │
                                    │  │  └──────────┘   │  │
                                    │  └────────┬────────┘  │
                                    │           │           │
                                    │  ┌────────▼────────┐  │
                                    │  │  SQLite DB      │  │
                                    │  │  ├─ devices     │  │
                                    │  │  ├─ ping_history│  │
                                    │  │  └─ alerts      │  │
                                    │  └─────────────────┘  │
                                    └─────────────────────┘
```

### Pembagian Tugas Komunikasi

| REST API | WebSocket |
|---|---|
| CRUD devices (tambah, edit, hapus) | Data monitoring real-time |
| Start/stop monitoring | Status change events |
| Fetch dashboard summary | Initial state saat client connect |
| Fetch alerts, resolve alert | Live ping/check results |

---

## Database Schema

### Tabel `devices`

| Kolom            | Tipe                               | Keterangan                                |
| ---------------- | ---------------------------------- | ----------------------------------------- |
| `id`             | INTEGER PRIMARY KEY AUTOINCREMENT  | —                                         |
| `name`           | TEXT NOT NULL                      | Nama device                               |
| `type`           | TEXT NOT NULL                      | Server/Router/Switch/Access Point/Website |
| `ip`             | TEXT NOT NULL                      | Alamat IP                                 |
| `url`            | TEXT DEFAULT ''                    | URL untuk HTTP check (opsional)           |
| `port`           | INTEGER                            | Port untuk TCP port check (opsional)      |
| `method`         | TEXT NOT NULL DEFAULT 'ICMP Ping'  | Metode monitoring                         |
| `location`       | TEXT DEFAULT ''                    | Lokasi device (opsional, user bisa input atau tidak) |
| `check_interval` | INTEGER DEFAULT 3                  | Interval pengecekan (detik)               |
| `description`    | TEXT DEFAULT ''                    | Deskripsi device                          |
| `created_at`     | DATETIME DEFAULT CURRENT_TIMESTAMP | —                                         |
| `updated_at`     | DATETIME DEFAULT CURRENT_TIMESTAMP | —                                         |

**Catatan:**
- Field `location` bersifat opsional. User bebas mau menambahkan lokasi atau tidak saat mengisi form.
- Tidak ada field `status` (active/inactive). Semua device yang terdaftar otomatis akan dimonitor.

### Tabel `ping_history`

| Kolom        | Tipe                               | Keterangan                      |
| ------------ | ---------------------------------- | ------------------------------- |
| `id`         | INTEGER PRIMARY KEY AUTOINCREMENT  | —                               |
| `device_id`  | INTEGER NOT NULL (FK)              | Relasi ke devices.id            |
| `status`     | TEXT NOT NULL                      | online/offline/warning          |
| `latency_ms` | REAL DEFAULT 0                     | Latency dalam milidetik         |
| `ttl`        | INTEGER DEFAULT 0                  | Time-to-live (untuk ICMP)       |
| `seq`        | INTEGER DEFAULT 0                  | Sequence number                 |
| `details`    | TEXT DEFAULT '{}'                  | JSON untuk data metode spesifik |
| `timestamp`  | DATETIME DEFAULT CURRENT_TIMESTAMP | —                               |

**Catatan field `details`:**
- ICMP: `{"ttl": 64}`
- HTTP: `{"status_code": 200, "content_type": "text/html"}`
- TCP: `{"connected": true}`
- SNMP: `{"uptime": 86400, "cpu_usage": 45.2}`
- Field ini membuat schema fleksibel tanpa perlu tambah kolom tiap metode baru

### Tabel `alerts`

| Kolom         | Tipe                               | Keterangan                               |
| ------------- | ---------------------------------- | ---------------------------------------- |
| `id`          | INTEGER PRIMARY KEY AUTOINCREMENT  | —                                        |
| `device_id`   | INTEGER NOT NULL (FK)              | Relasi ke devices.id                     |
| `title`       | TEXT NOT NULL                      | Judul alert                              |
| `status`      | TEXT NOT NULL DEFAULT 'ongoing'    | ongoing/resolved                         |
| `severity`    | TEXT NOT NULL DEFAULT 'low'        | low/medium/high/critical                 |
| `started_at`  | DATETIME DEFAULT CURRENT_TIMESTAMP | —                                        |
| `resolved_at` | DATETIME                           | Waktu resolved (null jika masih ongoing) |
| `description` | TEXT DEFAULT ''                    | Deskripsi alert                          |

### Cascade Rules

- Hapus device → otomatis hapus `ping_history` + `alerts` terkait
- SQLite WAL mode untuk performa concurrent read/write

---

## Backend API Endpoints

### Device CRUD — `handler/device.go`

| Method | Endpoint | Request Body | Response | Fungsi |
|---|---|---|---|---|
| `GET` | `/api/devices` | — | `{"success": true, "data": [...]}` | List semua device |
| `POST` | `/api/devices` | `{name, type, ip, method, port, url, ...}` | `{"success": true, "data": {...}}` | Tambah device |
| `PUT` | `/api/devices/{id}` | `{...fields}` | `{"success": true, "data": {...}}` | Update device |
| `DELETE` | `/api/devices/{id}` | — | `{"success": true, "message": "..."}` | Hapus device |
| `POST` | `/api/devices/{id}/start` | — | `{"success": true, "message": "..."}` | Mulai monitoring |
| `POST` | `/api/devices/{id}/stop` | — | `{"success": true, "message": "..."}` | Hentikan monitoring sementara |

**Catatan:**
- Semua device yang terdaftar otomatis dimonitor saat ditambahkan.
- Endpoint start/stop berguna untuk pause/resume monitoring sementara (misal saat maintenance).
- Tidak ada field `status` active/inactive di database.

### Alert Endpoints — `handler/alert.go`

| Method | Endpoint | Query Params | Response | Fungsi |
|---|---|---|---|---|
| `GET` | `/api/alerts` | `status`, `severity`, `device_type` | `{"success": true, "data": [...]}` | List alerts |
| `GET` | `/api/alerts/{id}` | — | `{"success": true, "data": {...}}` | Detail alert |
| `PUT` | `/api/alerts/{id}/resolve` | — | `{"success": true, "message": "..."}` | Resolve alert |

### Dashboard — `handler/dashboard.go`

| Method | Endpoint | Response | Fungsi |
|---|---|---|---|
| `GET` | `/api/dashboard` | `{summary: {total, online, offline, warning}, latestAlerts: [...]}` | Ringkasan sistem |

### Monitoring — `handler/monitoring.go`

| Method | Endpoint | Response | Fungsi |
|---|---|---|---|
| `GET` | `/api/monitoring` | `[{device, lastCheck, status, latency}]` | List device + status terkini |
| `GET` | `/api/monitoring/{id}/history` | `[{status, latency_ms, timestamp, ...}]` | Riwayat ping 50 terakhir |

---

## WebSocket Protocol

### Pesan dari Backend ke Frontend

```json
// Hasil check (setiap interval per device)
{
  "type": "check_result",
  "data": {
    "device_id": 1,
    "method": "ICMP Ping",
    "status": "online",
    "latency_ms": 12.3,
    "seq": 42,
    "timestamp": "2026-07-28T14:30:00Z",
    "details": {"ttl": 64}
  }
}

// Status berubah (online ↔ offline)
{
  "type": "status_change",
  "data": {
    "device_id": 1,
    "device_name": "Router-01",
    "old_status": "online",
    "new_status": "offline",
    "timestamp": "2026-07-28T14:30:00Z"
  }
}

// Status awal saat client pertama kali connect
{
  "type": "initial_state",
  "data": [
    {"device_id": 1, "name": "Router-01", "status": "online", "latency_ms": 12.3, ...},
    {"device_id": 2, "name": "Server-05", "status": "offline", "latency_ms": 0, ...}
  ]
}
```

### Pesan dari Frontend ke Backend

```json
{ "action": "start", "device_id": 1 }
{ "action": "stop", "device_id": 1 }
```

---

## Monitoring Engine (Multi-Method)

### Struktur Data

```go
type DeviceConfig struct {
    DeviceID int
    IP       string
    URL      string   // untuk HTTP check
    Port     int      // untuk TCP port check
    Method   string   // "ICMP Ping", "HTTP Check", "TCP Port"
    Interval int      // detik
}

type CheckResult struct {
    DeviceID  int                    `json:"device_id"`
    Method    string                 `json:"method"`
    Status    string                 `json:"status"`     // online/offline/warning
    Latency   float64                `json:"latency_ms"`
    Seq       int                    `json:"seq"`
    Timestamp string                 `json:"timestamp"`
    Details   map[string]interface{} `json:"details"`
}

type CheckFunc func(deviceID int, target string, port int, seq int) CheckResult
```

### Checker per Metode

| File | Fungsi | Cara Kerja |
|---|---|---|
| `monitor/ping.go` | `ICMPCheck()` | Jalankan `ping -c 1 -W 3 <ip>`, parse output |
| `monitor/http.go` (baru) | `HTTPCheck()` | `net/http.Get(url)`, ukur latency, cek status code |
| `monitor/tcp.go` (baru) | `TCPCheck()` | `net.DialTimeout("tcp", ip:port, 3s)`, cek connected |

**Untuk MVP:** ICMP Check sudah ada. HTTP dan TCP bisa ditambahkan belakangan. SNMP nanti lagi.

### Engine Logic

```go
type Engine struct {
    hub      HubInterface
    db       *sql.DB
    targets  map[int]context.CancelFunc  // key: device_id
    statuses map[int]string              // previousStatus per device
    fails    map[int]int                 // consecutiveFailures per device
    mu       sync.Mutex
}
```

### Status Change Detection + Auto Alert

```
Untuk setiap device, track:
  - previousStatus (online/offline)
  - consecutiveFailures (counter)

Logika pada setiap hasil check:
  Jika status == "offline":
    consecutiveFailures++
    Jika consecutiveFailures >= 3 DAN previousStatus == "online":
      → Buat alert di DB (severity: critical)
      → Broadcast "status_change" via WebSocket
  Jika status == "online":
    Jika previousStatus == "offline":
      → Resolve alert ongoing di DB
      → Broadcast "status_change" via WebSocket
    consecutiveFailures = 0
  previousStatus = status
```

### Initial State + Ping History

- Saat WebSocket client connect → langsung kirim `initial_state` (status semua device aktif)
- Setiap hasil check → simpan ke `ping_history` (dengan field `details` JSON)
- Broadcast `check_result` ke semua WebSocket clients

---

## Frontend Integration

### API Client — `src/lib/api.ts`

Fungsi-fungsi untuk memanggil semua REST API endpoint:

| Fungsi | Endpoint | Method |
|---|---|---|
| `fetchDevices()` | `/api/devices` | GET |
| `createDevice(data)` | `/api/devices` | POST |
| `updateDevice(id, data)` | `/api/devices/{id}` | PUT |
| `deleteDevice(id)` | `/api/devices/{id}` | DELETE |
| `startMonitor(id)` | `/api/devices/{id}/start` | POST |
| `stopMonitor(id)` | `/api/devices/{id}/stop` | POST |
| `fetchAlerts(filters)` | `/api/alerts` | GET |
| `resolveAlert(id)` | `/api/alerts/{id}/resolve` | PUT |
| `fetchDashboard()` | `/api/dashboard` | GET |
| `fetchMonitoring()` | `/api/monitoring` | GET |
| `fetchDeviceHistory(id)` | `/api/monitoring/{id}/history` | GET |

### WebSocket Hook — `src/hooks/useWebSocket.ts`

- Handle `initial_state` → populate semua device statuses
- Handle `check_result` → update status device spesifik
- Handle `status_change` → update status + trigger alert refresh
- Simpan `Map<number, MonitorResult>` (bukan cuma 1 result)

### Halaman yang Dimodifikasi

| Halaman | Perubahan |
|---|---|
| **Dashboard** | Fetch dari `GET /api/dashboard`, real-time update via WebSocket |
| **Device Management** | Real CRUD via API, method dropdown (bukan disabled input), conditional fields |
| **Monitoring** (halaman baru) | List device + status real-time, filter by status/type, detail panel + riwayat ping |
| **Alert Center** | Fetch dari `GET /api/alerts`, resolve via API |

### Device Form — Multi-Method Ready

Form tambah/edit device harus:

- **Method dropdown**: ICMP Ping, HTTP Check, TCP Port (SNMP nanti)
- **Kondisional fields**:
  - ICMP Ping: cuma IP
  - HTTP Check: IP/URL + Path (opsional, default `/`)
  - TCP Port: IP + Port (wajib)
- Port wajib diisi hanya untuk TCP Port
- **Location**: opsional, user bisa isi lokasi device atau biarkan kosong

---

## Halaman Monitoring (Halaman Baru)

### Tanggung Jawab

Monitoring menjawab pertanyaan: **"Perangkat mana yang sedang mengalami kondisi tersebut?"**

Monitoring **tidak** bertugas untuk:
- Menampilkan ringkasan sistem (itu Dashboard)
- Mengelola konfigurasi perangkat (itu Device Management)
- Menampilkan riwayat alert (itu Alert Center)

### Sections

| Section | Konten |
|---|---|
| **Monitoring Summary** | Total Device, Online, Offline, Warning, Unknown |
| **Filter & Search** | Search device, filter by status, filter by device type |
| **Device List** | Device Name, Device Type, Monitoring Status, Last Check Time |
| **Device Detail** | Panel detail: name, type, IP, status, method, last check, interval, deskripsi |

### User Flow

1. Pengguna membuka halaman Monitoring
2. Melihat Monitoring Summary (kondisi keseluruhan)
3. Filter/search perangkat yang ingin dilihat
4. Pilih perangkat → lihat Device Detail
5. Jika perangkat gangguan → bisa ke Alert Center

---

## Implementasi (6 Fase)

### Fase 1: Database Layer

| Aktivitas | File |
|---|---|
| Add dependency `go get modernc.org/sqlite` | `go.mod` |
| Koneksi SQLite + auto-migration | `database/db.go` |
| Struct definitions | `database/models.go` |

### Fase 2: Backend API

| Aktivitas | File |
|---|---|
| Device CRUD endpoints | `handler/device.go` |
| Alert endpoints | `handler/alert.go` |
| Dashboard summary endpoint | `handler/dashboard.go` |
| Monitoring endpoints | `handler/monitoring.go` |

### Fase 3: Engine Enhancement

| Aktivitas | File |
|---|---|
| Generalisasi CheckResult | `monitor/checker.go` (baru) |
| Multi-method dispatch | `monitor/engine.go` (modifikasi) |
| Adapt ICMP ke CheckResult baru | `monitor/ping.go` (modifikasi) |
| HTTP checker (opsional, bisa nanti) | `monitor/http.go` (baru) |
| TCP checker (opsional, bisa nanti) | `monitor/tcp.go` (baru) |

### Fase 4: Frontend API Layer

| Aktivitas | File |
|---|---|
| API client functions | `src/lib/api.ts` (baru) |

### Fase 5: Frontend Integration

| Aktivitas | File |
|---|---|
| Generalisasi types | `src/types/index.ts` |
| Multi-device WebSocket state | `src/hooks/useWebSocket.ts` |
| Dashboard real data | `src/pages/DashboardPage.tsx` |
| Device Management real CRUD | `src/pages/DeviceManagementPage.tsx` |
| Alert Center real data | `src/pages/AlertCenterPage.tsx` |
| Monitoring page (baru) | `src/pages/MonitoringPage.tsx` |
| Multi-method device form | `src/components/DeviceFormModal.tsx` |
| Add monitoring nav | `src/components/Sidebar.tsx` |
| Add monitoring route | `src/App.tsx` |

### Fase 6: Wiring

| Aktivitas | File |
|---|---|
| Init DB, register routes, startup logic | `main.go` |

---

## Ringkasan File

| Fase | File Baru | File Dimodifikasi |
|---|---|---|
| 1 | `database/db.go`, `database/models.go` | — |
| 2 | `handler/device.go`, `handler/alert.go`, `handler/dashboard.go`, `handler/monitoring.go` | — |
| 3 | `monitor/checker.go`, `monitor/http.go`, `monitor/tcp.go` | `monitor/engine.go`, `monitor/ping.go` |
| 4 | `src/lib/api.ts` | — |
| 5 | `src/pages/MonitoringPage.tsx` | `types/index.ts`, `useWebSocket.ts`, `DashboardPage.tsx`, `DeviceManagementPage.tsx`, `AlertCenterPage.tsx`, `DeviceFormModal.tsx`, `Sidebar.tsx`, `App.tsx` |
| 6 | — | `main.go` |

**Total: ~11 file baru, ~11 file dimodifikasi**

---

## Urutan Build yang Disarankan

```
1. database/db.go              (fondasi)
2. monitor/checker.go          (generalisasi result)
3. monitor/engine.go           (multi-method dispatch)
4. monitor/ping.go             (adapt ke CheckResult baru)
5. handler/device.go           (CRUD API)
6. handler/dashboard.go        (summary API)
7. handler/alert.go            (alert API)
8. handler/monitoring.go       (monitoring API)
9. main.go                     (wiring)
10. src/lib/api.ts             (frontend API client)
11. src/types/index.ts         (generalisasi types)
12. src/hooks/useWebSocket.ts  (multi-device state)
13. src/pages/*                (ganti dummy → real data)
14. Sidebar + App.tsx          (routing)
```

---

## Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| SQLite write contention (banyak ping nulis ke DB) | WAL mode + batch insert setiap 1 detik |
| WebSocket reconnect kehilangan state | Re-fetch dari REST API setelah reconnect |
| Gak ada auth untuk MVP | Acceptable untuk jaringan internal NOC |
| Ping butuh root di beberapa Linux | Pakai binary ping system (sudah punya suid bit) |
| Frontend build error sebelum integrasi | Fix pre-existing errors (Dashboard.tsx StatusCard, Sidebar.tsx JSX namespace) |

---

## Fitur yang Tidak Termasuk dalam MVP

- User authentication (multi-user, role)
- Email/SMS/Telegram alerting
- Configuration file (config.yaml) — masih hardcoded
- Structured logging (slog)
- Auto-discovery (Nmap)
- SNMP polling
- CPU/RAM monitoring
- Grafik latency jangka panjang
- Docker/deployment config

---

## Catatan untuk Pengembangan Selanjutnya

Setelah MVP berfungsi, metode monitoring baru bisa ditambahkan dengan:

1. Buat file `monitor/http.go` dengan fungsi `HTTPCheck()`
2. Tambahkan case `"HTTP Check"` di switch engine
3. Tambahkan option di frontend DeviceFormModal
4. Tambahkan conditional fields (URL, path)

Pola yang sama berlaku untuk TCP Port, SNMP, dan metode lainnya. Fondasi yang dibangun di MVP ini dirancang untuk扩展 tanpa mengubah arsitektur inti.
