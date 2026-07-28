# Engine Enhancement — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian Monitoring Engine

Monitoring engine merupakan komponen inti dalam sistem monitoring jaringan yang bertanggung jawab untuk secara aktif memeriksa status perangkat secara berkala. Engine ini menjalankan siklus pengecekan berulang (loop) untuk setiap perangkat yang didaftarkan, mengirimkan sinyal ping ke alamat IP perangkat, dan mencatat hasilnya.

Dalam konteks Gamon, monitoring engine berfungsi sebagai "jantung" sistem yang terus berdetak, memastikan setiap perangkat jaringan selalu dalam kondisi terpantau. Tanpa engine, sistem Gamon hanya dapat menampilkan data statis tanpa pembaruan status secara real-time.

### 1.2 Peran Monitoring Engine dalam Gamon

Monitoring engine dalam Gamon memiliki beberapa peran penting:

| Peran | Penjelasan |
|---|---|
| **Pengecekan Berkala** | Menjalankan ping ke setiap perangkat aktif secara otomatis sesuai interval yang ditentukan |
| **Penyimpanan Hasil** | Mencatat setiap hasil pengecekan ke dalam database untuk riwayat dan analisis |
| **Deteksi Status** | Mendeteksi perubahan status perangkat (online → offline, online → warning) |
| **Auto Alert** | Secara otomatis membuat alert saat perangkat mengalami gangguan |
| **Auto Resolve** | Secara otomatis menyelesaikan alert saat perangkat kembali normal |
| **Real-time Broadcast** | Mengirimkan hasil pengecekan ke semua client WebSocket secara real-time |

### 1.3 Mengapa Engine Perlu Ditingkatkan

Pada implementasi sebelumnya (Fase 2), monitoring engine memiliki keterbatasan:

| Aspek | Sebelumnya | Sesudahnya |
|---|---|---|
| **Penyimpanan Data** | Hanya broadcast ke WebSocket, tidak menyimpan ke database | Menyimpan setiap hasil ke `ping_history` |
| **Deteksi Status** | Tidak ada tracking perubahan status | Track perubahan status dengan in-memory cache |
| **Alert Otomatis** | Tidak ada alert saat device offline | Auto alert saat device offline/warning |
| **Status Threshold** | Hanya online/offline | Tambah status warning (latency >= 200ms) |
| **Active/Inactive** | Semua device selalu dimonitor | User bisa aktifkan/nonaktifkan monitoring |
| **Auto-Start** | User harus manual start setiap device | Semua device active otomatis dimonitor saat server start |
| **Initial State** | Client harus menunggu ping pertama | Client langsung dapat status terkini saat connect |

---

## 2. Arsitektur Engine

### 2.1 Komponen Engine

Engine enhancement terdiri dari beberapa komponen utama:

| Komponen | File | Fungsi |
|---|---|---|
| **Engine Core** | `monitor/engine.go` | Manajemen goroutine, status tracking, auto alert |
| **Ping Executor** | `monitor/ping.go` | Menjalankan ping system command, parsing output, threshold |
| **WebSocket Hub** | `handler/websocket.go` | Initial state, broadcast messages |
| **Device Handler** | `handler/device.go` | Active/inactive toggle, start/stop monitoring |
| **Main Server** | `main.go` | Auto-start monitoring saat boot |

### 2.2 Diagram Arsitektur

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         MONITORING ENGINE                                   │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                        Engine Struct                                 │   │
│  │                                                                      │   │
│  │  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐               │   │
│  │  │ hub         │  │ db           │  │ targets      │               │   │
│  │  │ (WebSocket) │  │ (SQLite)     │  │ (cancel funcs│               │   │
│  │  └──────┬──────┘  └──────┬───────┘  └──────┬───────┘               │   │
│  │         │                │                  │                        │   │
│  │         │                │                  │                        │   │
│  │  ┌──────┴──────┐  ┌──────┴───────┐  ┌──────┴───────┐               │   │
│  │  │ Broadcast   │  │ Save/Query   │  │ Start/Stop   │               │   │
│  │  │ Messages    │  │ Ping History │  │ Monitoring   │               │   │
│  │  └─────────────┘  └──────────────┘  └──────────────┘               │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │                  Check Loop (per device)                    │    │   │
│  │  │                                                             │    │   │
│  │  │  ┌─────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │    │   │
│  │  │  │ PingOnce│→ │ SavePing │→ │ Check    │→ │ Broadcast │  │    │   │
│  │  │  │         │  │ Result   │  │ Status   │  │ Result    │  │    │   │
│  │  │  └─────────┘  └──────────┘  └──────────┘  └───────────┘  │    │   │
│  │  │       ↑                                               │    │    │   │
│  │  │       └───────────────────────────────────────────────┘    │    │   │
│  │  │                    (every interval seconds)                │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│  ┌──────────────────────────────────────────────────────────────────────┐   │
│  │                    Status Change Handler                             │   │
│  │                                                                      │   │
│  │  ┌─────────────────────────────────────────────────────────────┐    │   │
│  │  │  handleStatusChange(deviceID, oldStatus, newStatus)         │    │   │
│  │  │                                                             │    │   │
│  │  │  online → offline  → createAlert (severity: high)          │    │   │
│  │  │  online → warning  → createAlert (severity: medium)        │    │   │
│  │  │  offline → online  → resolveAlert                          │    │   │
│  │  │  warning → online  → resolveAlert                          │    │   │
│  │  └─────────────────────────────────────────────────────────────┘    │   │
│  └──────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Hubungan dengan Komponen Lain

```
┌─────────────┐
│   Browser   │
│   (React)   │
└──────┬──────┘
       │
       │ WebSocket (initial_state, ping_result, alert_created, alert_resolved)
       │
┌──────▼──────┐
│  WebSocket  │
│    Hub      │
└──────┬──────┘
       │
       │ Broadcast()
       │
┌──────▼──────┐
│   Engine    │ ──► PingOnce() ──► System Ping Command
│  (Monitor)  │
└──────┬──────┘
       │
       │ savePingResult() / createAlert() / resolveAlert()
       │
┌──────▼──────┐
│  Database   │ ──► ping_history, alerts, devices
│  (SQLite)   │
└─────────────┘
```

---

## 3. Active/Inactive System

### 3.1 Pengertian

Active/Inactive system merupakan mekanisme untuk mengontrol apakah sebuah perangkat harus dimonitor atau tidak. Saat user menambahkan perangkat baru, ia dapat memilih untuk mengaktifkan atau menonaktifkan monitoring segera.

### 3.2 Status Device

| Status | Penjelasan | Dampak |
|---|---|---|
| `active` | Perangkat aktif dan dimonitor | Engine menjalankan ping secara berkala |
| `inactive` | Perangkat tidak dimonitor | Engine tidak menjalankan ping |

### 3.3 Endpoint Toggle Status

| Method | Endpoint | Fungsi |
|---|---|---|
| `PUT` | `/api/devices/{id}/status` | Toggle status active/inactive |

**Request Body:**
```json
{
    "status": "active"
}
```

**Response:**
```json
{
    "success": true,
    "data": {
        "id": 1,
        "status": "active",
        "message": "Device activated"
    }
}
```

### 3.4 Flow Toggle Status

```
User klik "Activate"
  │
  ▼
PUT /api/devices/1/status
  │
  ▼
Handler: Update devices.status = 'active'
  │
  ▼
Handler: engine.Start(config)
  │
  ▼
Engine: Mulai checkLoop untuk device ini

────────────────────────────────────

User klik "Deactivate"
  │
  ▼
PUT /api/devices/1/status
  │
  ▼
Handler: Update devices.status = 'inactive'
  │
  ▼
Handler: engine.Stop(id)
  │
  ▼
Engine: Hentikan checkLoop untuk device ini
```

### 3.5 Proteksi Start Monitoring

Saat user mencoba memulai monitoring untuk device inactive:

```json
POST /api/devices/5/start
```

Response:
```json
{
    "success": false,
    "message": "Device is inactive. Activate it first."
}
```

---

## 4. Database Write Integration

### 4.1 Mengapa Perlu Menyimpan ke Database

Sebelumnya, engine hanya mengirim hasil ping ke WebSocket tanpa menyimpannya. Hal ini memiliki beberapa keterbatasan:

| Tanpa Database | Dengan Database |
|---|---|
| Riwayat hilang saat WebSocket disconnect | Riwayat tersimpan permanen |
| Tidak bisa analisis historis | Query SQL untuk analisis |
| Dashboard tidak bisa tampilkan data | Dashboard ambil data dari database |
| Alert tidak ada dasar data | Alert didukung oleh riwayat pengecekan |

### 4.2 Cara Kerja Database Write

Setiap hasil ping disimpan ke tabel `ping_history`:

```sql
INSERT INTO ping_history (device_id, status, latency_ms, ttl, seq, details, timestamp)
VALUES (?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
```

### 4.3 Flow Database Write

```
PingOnce(ip, seq)
  │
  ▼
Result{Status: "online", Latency: 12.5, TTL: 64, Seq: 1}
  │
  ▼
savePingResult(result)
  │
  ▼
INSERT INTO ping_history
  │
  ▼
Data tersimpan di SQLite
```

### 4.4 Fire-and-Forget Design

Operasi database write menggunakan pola fire-and-forget:

```go
func (e *Engine) savePingResult(result Result) {
    _, err := e.db.Exec(...)
    if err != nil {
        log.Printf("Failed to save ping result: %v", err)
    }
}
```

**Alasan:**
1. **Tidak Blok Check Loop** — Jika database write gagal, check loop tetap berlanjut
2. **Tolerance** — Satu failed write tidak mengganggu monitoring keseluruhan
3. **Logging** — Error tetap dicatat untuk debugging

---

## 5. Status Tracking

### 5.1 Pengertian Status Tracking

Status tracking adalah mekanisme untuk mendeteksi perubahan status perangkat dari satu kondisi ke kondisi lain. Engine menyimpan status terakhir setiap perangkat dalam memori dan membandingkannya dengan status terbaru dari database.

### 5.2 In-Memory Cache

Engine menggunakan map `lastStatus` untuk menyimpan status terakhir:

```go
type Engine struct {
    lastStatus map[int]string  // deviceID → status terakhir
}
```

### 5.3 Flow Status Tracking

```
Setiap kali ping selesai:
  │
  ▼
savePingResult() → INSERT status baru ke database
  │
  ▼
checkStatusChange(deviceID, ip)
  │
  ▼
Query database: Ambil status terakhir device ini
  │
  ▼
Bandingkan dengan lastStatus[deviceID] (in-memory cache)
  │
  ├─ Sama? → Tidak ada perubahan
  │
  └─ Berbeda? → handleStatusChange()
       │
       ▼
     Update lastStatus[deviceID] = status baru
```

### 5.4 Keuntungan In-Memory Cache

| Pendekatan | Kelebihan | Kekurangan |
|---|---|---|
| Query DB setiap kali | Selalu up-to-date | Lambat, overhead I/O |
| In-memory cache | Cek cepat, minimal I/O | Bisa stale jika crash |
| **Kombinasi** | **Cek cache, validasi via DB** | **Balance speed & accuracy** |

Gamon menggunakan kombinasi: cache untuk perbandingan cepat, database untuk validasi.

---

## 6. Status Threshold

### 6.1 Definisi Status

| Status | Kondisi | Keterangan |
|---|---|---|
| `online` | Ping berhasil, latency < 200ms | Device normal, merespons dengan baik |
| `warning` | Ping berhasil, latency >= 200ms | Device lambat tapi masih merespons |
| `offline` | Ping gagal (timeout/error) | Device tidak merespons sama sekali |
| `unknown` | Belum ada data ping | Device baru ditambahkan, belum pernah di-ping |

### 6.2 Threshold Logic

```go
const LatencyWarningThreshold = 200.0

func PingOnce(ip string, seq int) Result {
    // ... jalankan ping ...
    
    if ping gagal {
        return Result{Status: "offline"}
    }
    
    if ping berhasil {
        parse latency
        
        if latency < 200ms {
            return Result{Status: "online", Latency: latency}
        } else {
            return Result{Status: "warning", Latency: latency}
        }
    }
}
```

### 6.3 Diagram Threshold

```
                    ┌─────────────────┐
                    │   Ping Result   │
                    └────────┬────────┘
                             │
                ┌────────────▼────────────┐
                │  Ping berhasil?         │
                └────────────┬────────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
              ▼                             ▼
        ┌─────────┐                   ┌─────────┐
        │  YA     │                   │  TIDAK  │
        └────┬────┘                   └────┬────┘
             │                             │
             ▼                             ▼
    ┌────────────────┐              ┌────────────┐
    │ Latency < 200ms│              │  OFFLINE   │
    └───────┬────────┘              └────────────┘
            │
   ┌────────┴────────┐
   │                 │
   ▼                 ▼
┌──────┐       ┌──────────┐
│ONLINE│       │ WARNING  │
└──────┘       └──────────┘
```

---

## 7. Auto Alert Generation

### 7.1 Pengertian

Auto alert generation adalah mekanisme di mana engine secara otomatis membuat alert saat mendeteksi perubahan status perangkat yang mengindikasikan gangguan.

### 7.2 Trigger Alert

| Transisi | Alert Title | Severity | Description |
|---|---|---|---|
| `online → offline` | Device Offline | `high` | Device tidak merespons ping dari [IP] |
| `online → warning` | Device High Latency | `medium` | Latency lebih dari 200ms |

### 7.3 Flow Auto Alert

```
checkStatusChange()
  │
  ▼
Detect: online → offline
  │
  ▼
handleStatusChange(deviceID, "online", "offline")
  │
  ▼
createAlert(deviceID, "Device Offline", "high", "Device tidak merespons...")
  │
  ▼
INSERT INTO alerts (device_id, title, status, severity, description)
  │
  ▼
hub.Broadcast("alert_created", alertData)
  │
  ▼
Semua WebSocket client menerima alert
```

### 7.4 SQL Insert Alert

```sql
INSERT INTO alerts (device_id, title, status, severity, description)
VALUES (?, ?, 'ongoing', ?, ?)
```

### 7.5 WebSocket Broadcast

```json
{
    "type": "alert_created",
    "data": {
        "device_id": 1,
        "title": "Device Offline",
        "severity": "high",
        "description": "Device tidak merespons ping dari 192.168.1.1"
    }
}
```

---

## 8. Auto Alert Resolve

### 8.1 Pengertian

Auto alert resolve adalah mekanisme di mana engine secara otomatis menyelesaikan (resolve) alert yang masih ongoing saat perangkat kembali ke status normal.

### 8.2 Trigger Resolve

| Transisi | Aksi |
|---|---|
| `offline → online` | Resolve semua alert ongoing untuk device ini |
| `warning → online` | Resolve semua alert ongoing untuk device ini |
| `unknown → online` | Resolve semua alert ongoing untuk device ini |

### 8.3 Flow Auto Resolve

```
checkStatusChange()
  │
  ▼
Detect: offline → online
  │
  ▼
handleStatusChange(deviceID, "offline", "online")
  │
  ▼
resolveAlert(deviceID)
  │
  ▼
UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP
WHERE device_id = ? AND status = 'ongoing'
  │
  ▼
hub.Broadcast("alert_resolved", alertData)
  │
  ▼
Semua WebSocket client menerima notifikasi
```

### 8.4 SQL Update Alert

```sql
UPDATE alerts 
SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP 
WHERE device_id = ? AND status = 'ongoing'
```

### 8.5 Bulk Resolve

Semua alert ongoing untuk device yang sama di-resolve sekaligus:

```
Device A punya 3 alert ongoing:
  - Alert 1: "Device Offline" (high)
  - Alert 2: "Device High Latency" (medium)
  - Alert 3: "Device Offline" (high)

Device A kembali online:
  → Semua 3 alert di-resolve sekaligus
```

---

## 9. Initial State

### 9.1 Pengertian

Initial state adalah mekanisme untuk mengirim status terkini semua perangkat aktif saat WebSocket client pertama kali terhubung. Tanpa initial state, client harus menunggu ping pertama sebelum bisa menampilkan status.

### 9.2 Mengapa Initial State Diperlukan

| Tanpa Initial State | Dengan Initial State |
|---|---|
| Client blank saat pertama kali connect | Langsung tampilkan status terkini |
| Harus menunggu ping pertama (3+ detik) | Tidak perlu menunggu |
| UX buruk | UX lebih baik |

### 9.3 Flow Initial State

```
WebSocket Client Connect
  │
  ▼
hub.register <- client
  │
  ▼
go sendInitialState(client)
  │
  ▼
Tunggu 100ms (pastikan registration selesai)
  │
  ▼
getAllDeviceStatuses()
  │
  ▼
Query database:
  SELECT d.id, d.name, d.type, d.ip,
         COALESCE(ph.status, 'unknown'),
         COALESCE(ph.latency_ms, 0),
         COALESCE(ph.timestamp, d.created_at)
  FROM devices d
  LEFT JOIN ping_history ph ON ...
  WHERE d.status = 'active'
  │
  ▼
Kirim ke client:
  {
    "type": "initial_state",
    "data": [
      {device_id: 1, name: "Server A", status: "online", latency_ms: 12.5, ...},
      {device_id: 2, name: "Router B", status: "offline", latency_ms: 0, ...}
    ]
  }
```

### 9.4 SQL Query Initial State

```sql
SELECT d.id, d.name, d.type, d.ip,
    COALESCE(ph.status, 'unknown') as last_status,
    COALESCE(ph.latency_ms, 0) as last_latency,
    COALESCE(ph.timestamp, d.created_at) as last_check
FROM devices d
LEFT JOIN ping_history ph ON ph.id = (
    SELECT id FROM ping_history WHERE device_id = d.id ORDER BY id DESC LIMIT 1
)
WHERE d.status = 'active'
ORDER BY d.name
```

### 9.5 Correlated Subquery

Query menggunakan correlated subquery untuk mengambil hanya record terakhir:

```sql
ph.id = (
    SELECT id FROM ping_history 
    WHERE device_id = d.id 
    ORDER BY id DESC 
    LIMIT 1
)
```

**Mengapa correlated subquery:**
- Mengambil hanya 1 record terakhir per device
- Lebih efektif dibandingkan JOIN semua record
- COALESCE menangani device yang belum pernah di-ping

---

## 10. Auto-Start Monitoring

### 10.1 Pengertian

Auto-start monitoring adalah mekanisme di mana server secara otomatis memulai monitoring untuk semua perangkat active saat pertama kali dijalankan.

### 10.2 Mengapa Auto-Start Diperlukan

| Tanpa Auto-Start | Dengan Auto-Start |
|---|---|
| User harus manual start setiap device | Semua device active langsung dimonitor |
| Ribet jika banyak device | Praktis, zero-config |
| Lupa start = tidak termonitor | Selalu termonitor |

### 10.3 Flow Auto-Start

```
Server Start (main.go)
  │
  ▼
go autoStartMonitoring(db, engine)
  │
  ▼
Tunggu 2 detik (pastikan server siap)
  │
  ▼
Query semua device active:
  SELECT id, ip, method, url, port, check_interval
  FROM devices
  WHERE status = 'active'
  │
  ▼
Untuk setiap device:
  │
  ├─ Buat DeviceConfig
  │
  ├─ engine.Start(config)
  │
  └─ Log: "Auto-started monitoring for device X"
```

### 10.4 Idempotent Check

Engine.Start() memiliki check idempotent:

```go
func (e *Engine) Start(config DeviceConfig) {
    e.mu.Lock()
    defer e.mu.Unlock()

    if _, exists := e.targets[config.DeviceID]; exists {
        log.Printf("Already monitoring device %d", config.DeviceID)
        return
    }
    // ... mulai monitoring ...
}
```

**Mengapa idempotent:**
- Mencegah duplicate monitoring jika auto-start dan manual start berbarengan
- Aman dipanggil berulang kali tanpa efek samping

---

## 11. WebSocket Messages

### 11.1 Tipe Pesan

| Type | Kapan Dikirim | Data |
|---|---|---|
| `initial_state` | Client pertama kali connect | Array semua device active + status terkini |
| `ping_result` | Setiap kali ping selesai | `{device_id, status, latency_ms, ttl, seq, ...}` |
| `alert_created` | Status berubah ke offline/warning | `{device_id, title, severity, description}` |
| `alert_resolved` | Device kembali online | `{device_id}` |

### 11.2 Format Pesan

Semua pesan menggunakan format envelope:

```json
{
    "type": "pesan_type",
    "data": { ... }
}
```

### 11.3 Contoh Pesan

**Initial State:**
```json
{
    "type": "initial_state",
    "data": [
        {
            "device_id": 1,
            "name": "Server A",
            "type": "server",
            "ip": "192.168.1.1",
            "status": "online",
            "latency_ms": 12.5,
            "last_check": "2025-01-15T10:30:00Z"
        }
    ]
}
```

**Ping Result:**
```json
{
    "type": "ping_result",
    "data": {
        "device_id": 1,
        "ip": "192.168.1.1",
        "method": "ICMP Ping",
        "status": "online",
        "latency": 12.5,
        "ttl": 64,
        "seq": 5,
        "timestamp": "2025-01-15T10:30:00Z"
    }
}
```

**Alert Created:**
```json
{
    "type": "alert_created",
    "data": {
        "device_id": 1,
        "title": "Device Offline",
        "severity": "high",
        "description": "Device tidak merespons ping dari 192.168.1.1"
    }
}
```

**Alert Resolved:**
```json
{
    "type": "alert_resolved",
    "data": {
        "device_id": 1
    }
}
```

---

## 12. Status Transition Matrix

### 12.1 Tabel Transisi Lengkap

| Old Status | New Status | Aksi | Alert |
|---|---|---|---|
| `online` | `offline` | createAlert | severity: high |
| `online` | `warning` | createAlert | severity: medium |
| `offline` | `online` | resolveAlert | - |
| `warning` | `online` | resolveAlert | - |
| `unknown` | `online` | resolveAlert | - |
| `online` | `online` | tidak ada | - |
| `offline` | `warning` | tidak ada | - |
| `offline` | `unknown` | tidak ada | - |
| `warning` | `offline` | tidak ada | - |
| `warning` | `warning` | tidak ada | - |
| `unknown` | `offline` | tidak ada | - |
| `unknown` | `warning` | tidak ada | - |

### 12.2 Prinsip Transisi

1. **Hanya degraded dari online yang buat alert** — Alert hanya dibuat saat perangkat turun dari kondisi baik (online) ke kondisi buruk (offline/warning). Transisi lain tidak membuat alert untuk mencegah duplikat.

2. **Semua recovery resolve** — Saat perangkat kembali online dari kondisi apapun, semua alert ongoing di-resolve.

3. **No cascading alerts** — Transisi offline → warning atau warning → offline tidak membuat alert baru.

### 12.3 Diagram Transisi

```
                    ┌─────────────────────────────────────┐
                    │                                     │
                    │              ONLINE                 │
                    │                                     │
                    └───────┬─────────────────┬───────────┘
                            │                 │
                 latency    │                 │  ping gagal
                 >= 200ms   │                 │
                            │                 │
                            ▼                 ▼
                    ┌───────────────┐ ┌───────────────┐
                    │   WARNING     │ │   OFFLINE     │
                    │               │ │               │
                    └───────┬───────┘ └───────┬───────┘
                            │                 │
                            │   ping gagal    │   latency < 200ms
                            │                 │
                            ▼                 ▼
                    ┌───────────────┐ ┌───────────────┐
                    │   OFFLINE     │ │   WARNING     │
                    │               │ │               │
                    └───────────────┘ └───────────────┘

        Semua panah ke ONLINE = resolveAlert
        Panah dari ONLINE ke WARNING/OFFLINE = createAlert
```

---

## 13. Kesimpulan

### 13.1 Prinsip Desain Engine Enhancement

Engine enhancement Gamon dirancang dengan beberapa prinsip utama:

1. **Persistence** — Setiap hasil pengecekan disimpan ke database, bukan hanya di-cache sementara. Ini menjamin riwayat lengkap dan mendukung fitur dashboard serta analisis.

2. **Reactive Alerting** — Alert dibuat dan di-resolve secara otomatis berdasarkan perubahan status. Tidak perlu intervensi manual untuk alert rutin.

3. **Efficiency** — Status tracking menggunakan in-memory cache untuk perbandingan cepat, dengan database sebagai source of truth.

4. **User Control** — Sistem active/inactive memberikan kendali penuh kepada user untuk mengontrol perangkat mana yang perlu dimonitor.

5. **Real-time** — Semua perubahan status langsung di-broadcast ke semua client melalui WebSocket.

### 13.2 Capaian Fase 3

Dengan selesainya Engine Enhancement, Gamon kini memiliki:

| Komponen | Status | Fungsi |
|---|---|---|
| **Database Write** | ✅ Selesai | Setiap hasil ping tersimpan di `ping_history` |
| **Status Tracking** | ✅ Selesai | Deteksi perubahan status dengan in-memory cache |
| **Status Threshold** | ✅ Selesai | online (< 200ms), warning (>= 200ms), offline |
| **Auto Alert** | ✅ Selesai | Alert otomatis saat device offline/warning |
| **Auto Resolve** | ✅ Selesai | Alert otomatis resolve saat device kembali online |
| **Active/Inactive** | ✅ Selesai | User bisa aktifkan/nonaktifkan monitoring |
| **Initial State** | ✅ Selesai | Client langsung dapat status terkini saat connect |
| **Auto-Start** | ✅ Selesai | Semua device active otomatis dimonitor saat server start |

### 13.3 Lanjutan ke Fase Berikutnya

Dengan engine yang sudah lengkap, langkah selanjutnya adalah:

1. **Phase 4 — Frontend Integration**: Menghubungkan frontend React dengan backend API yang sebenarnya, menggantikan dummy data dengan data real dari engine.

---

*Dokumentasi ini merupakan bagian dari laporan PKL (Praktek Kerja Lapangan) untuk proyek Gamon (Garda Monitoring) — Sistem Monitoring Jaringan Berbasis Web.*
