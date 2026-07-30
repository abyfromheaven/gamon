# Wiring — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian Wiring

Wiring dalam pengembangan perangkat lunak merujuk pada proses menyambung berbagai komponen sistem menjadi satu kesatuan yang dapat berjalan. Setiap komponen — database, handler, engine, middleware — dibuat secara terpisah pada fase-fase sebelumnya. Wiring memastikan komponen-komponen tersebut saling terhubung melalui dependency injection, registrasi route, dan startup logic sehingga aplikasi dapat berfungsi sebagai produk yang utuh.

Dalam Gamon, file `main.go` merupakan titik pusat wiring. File ini tidak berisi logika bisnis; ia hanya bertanggung jawab untuk membuat instance setiap komponen, menyuntikkan dependency yang dibutuhkan, mendaftarkan route HTTP, dan menjalankan server.

### 1.2 Peran Wiring dalam Gamon

| Peran | Penjelasan |
|---|---|
| **Inisialisasi Database** | Membuka koneksi SQLite, menjalankan migrasi, mengatur konfigurasi |
| **Pembuatan Komponen** | Membuat instance Hub, Engine, dan semua Handler |
| **Dependency Injection** | Menyuntikkan db, engine, dan hub ke handler yang membutuhkan |
| **Registrasi Route** | Menghubungkan URL path ke handler yang sesuai |
| **Middleware** | Membungkus route dengan CORS agar frontend dapat mengakses |
| **Auto-Start** | Memulai monitoring untuk semua device active saat server boot |
| **Server Start** | Menjalankan HTTP server di port 8080 |

### 1.3 Mengapa Wiring Diperlukan

Tanpa wiring, setiap komponen hanya berupa blueprint yang belum dihubungkan. Database belum dibuka, handler belum mendaftarkan route, engine belum menjalankan ping. Wiring mengubah kumpulan kode terpisah menjadi aplikasi yang berjalan.

| Tanpa Wiring | Dengan Wiring |
|---|---|
| Database belum terbuka | Koneksi SQLite aktif dan siap digunakan |
| Handler belum mendaftarkan route | Semua endpoint API terdaftar di http.ServeMux |
| Engine belum menjalankan monitoring | Goroutine monitoring aktif untuk semua device |
| WebSocket belum menerima client | Hub siap menerima dan menyiarkan pesan |
| Aplikasi tidak bisa dijalankan | Server berjalan di port 8080, siap melayani request |

---

## 2. Arsitektur Aplikasi

### 2.1 Dependency Graph

Berikut diagram ketergantungan seluruh komponen dalam Gamon:

```text
main()
  │
  ├──► database.NewDB()
  │         │
  │         ▼
  │    *sql.DB ──────────────────────────────────────────────┐
  │         │                                                 │
  │         │                                                 │
  ├──► handler.NewHub(db)                                     │
  │         │                                                 │
  │         ▼                                                 │
  │    *Hub ─────────────────────────────────────────────┐    │
  │         │                                             │    │
  │         │                                             │    │
  ├──► monitor.NewEngine(hub, db)                         │    │
  │         │                                             │    │
  │         ▼                                             │    │
  │    *Engine ──────────────────────────────────────┐    │    │
  │         │                                         │    │    │
  │         │                                         │    │    │
  ├──► handler.NewDeviceHandler(db, engine, hub)      │    │    │
  │         │                                         │    │    │
  ├──► handler.NewAlertHandler(db)                    │    │    │
  │         │                                         │    │    │
  ├──► handler.NewDashboardHandler(db)                │    │    │
  │         │                                         │    │    │
  ├──► handler.NewMonitoringHandler(db)               │    │    │
  │         │                                         │    │    │
  │         ▼                                         ▼    ▼    │
  │    ┌─────────────────────────────────────────────────┐     │
  │    │              http.ServeMux (Router)              │     │
  │    │                                                  │     │
  │    │  /api/devices    → DeviceHandler.HandleDevices  │     │
  │    │  /api/devices/   → DeviceHandler.HandleDevice   │     │
  │    │  /api/alerts     → AlertHandler.HandleAlerts    │     │
  │    │  /api/alerts/    → AlertHandler.HandleAlert     │     │
  │    │  /api/dashboard  → DashboardHandler             │     │
  │    │  /api/monitoring → MonitoringHandler            │     │
  │    │  /api/monitoring/→ MonitoringHandler            │     │
  │    │  /api/health     → API.Health                   │     │
  │    │  /ws             → HandleWebSocket              │     │
  │    └─────────────────────────────────────────────────┘     │
  │                     │                                      │
  │                     ▼                                      │
  │    ┌───────────────────────────────────────┐               │
  │    │         corsMiddleware(mux)            │               │
  │    └───────────────────────────────────────┘               │
  │                     │                                      │
  │                     ▼                                      │
  │    ┌───────────────────────────────────────┐               │
  │    │    http.ListenAndServe(":8080")        │               │
  │    └───────────────────────────────────────┘               │
  │                                                             │
  └────────────────────────────────────────────────────────────┘
```

### 2.2 Urutan Inisialisasi

Startup sequence dalam `main.go` berjalan secara berurutan:

```text
┌─────────────────────────────────────────────────────────────┐
│                    URUTAN STARTUP                            │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. database.NewDB()                                         │
│     ├── Buat direktori "data/" jika belum ada                │
│     ├── Buka koneksi SQLite (data/gamon.db)                 │
│     ├── Atur WAL mode + busy timeout                        │
│     ├── Set MaxOpenConns(1)                                 │
│     └── Jalankan migrate() → buat tabel devices,            │
│         ping_history, alerts + index                         │
│                                                              │
│  2. handler.NewHub(db)                                      │
│     └── Buat WebSocket Hub dengan database reference         │
│                                                              │
│  3. go hub.Run()                                            │
│     └── Jalankan goroutine untuk register/unregister/        │
│         broadcast clients                                   │
│                                                              │
│  4. monitor.NewEngine(hub, db)                              │
│     └── Buat Engine dengan Hub dan Database reference        │
│                                                              │
│  5. handler.New*Handler(db, engine, hub)                    │
│     └── Buat semua Handler dengan dependency injection       │
│                                                              │
│  6. Registrasi route ke http.ServeMux                        │
│     └── 10 HTTP route + 1 WebSocket endpoint                │
│                                                              │
│  7. corsMiddleware(mux)                                     │
│     └── Bungkus mux dengan CORS headers                     │
│                                                              │
│  8. go autoStartMonitoring(db, engine)                      │
│     └── Goroutine: tunggu 2 detik, query device active,     │
│         engine.Start() per device                            │
│                                                              │
│  9. http.ListenAndServe(":8080", wrapped)                   │
│     └── Mulai HTTP server, blocking call                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Inisialisasi Database

### 3.1 Pemanggilan NewDB()

```go
db, err := database.NewDB()
if err != nil {
    log.Fatalf("Failed to initialize database: %v", err)
}
defer db.Close()
```

`database.NewDB()` merupakan fungsi yang dieksekusi pertama kali. Fungsi ini membuka koneksi ke SQLite dan menjalankan migrasi.

**Yang dilakukan NewDB():**

| Langkah | Fungsi | Keterangan |
|---|---|---|
| 1 | `os.MkdirAll("data", 0755)` | Buat direktori penyimpanan jika belum ada |
| 2 | `sql.Open("sqlite", "data/gamon.db?_journal_mode=WAL&_busy_timeout=5000")` | Buka koneksi SQLite dengan WAL mode |
| 3 | `db.SetMaxOpenConns(1)` | Batasi hanya 1 koneksi aktif (SQLite limitation) |
| 4 | `db.Ping()` | Test koneksi |
| 5 | `migrate(db)` | Buat tabel dan index jika belum ada |

**Mengapa `log.Fatalf` pada error?**

Jika database tidak bisa dibuka, tidak ada gunanya melanjutkan server. `log.Fatalf` mencetak error dan menghentikan program dengan exit code 1. Ini adalah fail-fast pattern yang mencegah server berjalan tanpa database.

### 3.2 Migrasi Otomatis

Fungsi `migrate()` menggunakan `CREATE TABLE IF NOT EXISTS` sehingga aman dijalankan berulang kali:

```go
func migrate(db *sql.DB) {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS devices (...)`,
        `CREATE TABLE IF NOT EXISTS ping_history (...)`,
        `CREATE TABLE IF NOT EXISTS alerts (...)`,
        `CREATE INDEX IF NOT EXISTS idx_ping_history_device_id ...`,
        `CREATE INDEX IF NOT EXISTS idx_ping_history_timestamp ...`,
        `CREATE INDEX IF NOT EXISTS idx_alerts_device_id ...`,
        `CREATE INDEX IF NOT EXISTS idx_alerts_status ...`,
    }
    for _, query := range queries {
        db.Exec(query)
    }
}
```

**Mengapa migrasi dijalankan di main.go?**

Migrasi harus dijalankan sebelum komponen lain mencoba mengakses database. Dengan menempatkan migrasi di awal `NewDB()`, dijamin database sudah memiliki struktur tabel sebelum handler atau engine menggunakannya.

---

## 4. Inisialisasi WebSocket Hub

### 4.1 Pembuatan Hub

```go
hub := handler.NewHub(db)
go hub.Run()
```

Hub adalah pusat manajemen koneksi WebSocket. Ia menyimpan daftar semua client yang terhubung dan menyediakan mekanisme broadcast.

**Field Hub:**

```go
type Hub struct {
    clients    map[*Client]bool   // Daftar client aktif
    broadcast  chan []byte         // Channel untuk pesan broadcast
    register   chan *Client        // Channel untuk client baru
    unregister chan *Client        // Channel untuk client yang disconnect
    mu         sync.RWMutex       // Mutex untuk melindungi clients map
    db         *sql.DB            // Database reference untuk initial state
}
```

### 4.2 Goroutine Run()

`hub.Run()` dijalankan sebagai goroutine terpisah agar tidak memblokir inisialisasi komponen lain:

```go
go hub.Run()
```

**Yang dilakukan Run():**

```text
hub.Run() — infinite loop
    │
    ├─── case client := <-h.register:
    │       Tambah client ke map
    │       Kirim initial_state ke client (goroutine baru)
    │
    ├─── case client := <-h.unregister:
    │       Hapus client dari map
    │       Tutup channel client.send
    │
    └─── case message := <-h.broadcast:
            Kirim pesan ke SEMUA client
            Tutup client yang buffer-nya penuh
```

**Mengapa menggunakan goroutine terpisah?**

Hub harus berjalan terus-menerus untuk melayani register, unregister, dan broadcast. Jika dijalankan di main goroutine, program akan blocking dan tidak bisa menerima request HTTP.

---

## 5. Inisialisasi Engine

### 5.1 Pembuatan Engine

```go
engine := monitor.NewEngine(hub, db)
```

Engine merupakan komponen yang menjalankan monitoring loop untuk setiap device. Ia membutuhkan Hub (untuk broadcast) dan Database (untuk menyimpan hasil dan alert).

**Field Engine:**

```go
type Engine struct {
    hub        HubInterface          // WebSocket Hub
    db         *sql.DB               // Database
    mu         sync.Mutex            // Mutex untuk melindungi state
    targets    map[int]context.CancelFunc  // deviceID → cancel function
    lastStatus map[int]string        // deviceID → status terakhir
    failures   map[int]int           // deviceID → consecutive failures
    check      CheckFunc             // Fungsi pengecekan (injected untuk testing)
}
```

### 5.2 Default Check Function

Engine memiliki check function default yang memanggil `PingOnce()`:

```go
check: func(config DeviceConfig, seq int) CheckResult {
    result := PingOnce(config.IP, seq)
    result.DeviceID = config.DeviceID
    result.Method = config.Method
    return result
},
```

**Mengapa check function di-inject?**

Pola ini memungkinkan pengujian (testing) dengan check function palsu yang tidak menjalankan ping sesungguhnya:

```go
// Dalam test
engine := NewEngine(hub, db, WithCheckFunc(func(config DeviceConfig, seq int) CheckResult {
    return CheckResult{Status: "online", LatencyMs: 10}
}))
```

---

## 6. Dependency Injection

### 6.1 Pengertian Dependency Injection

Dependency injection (DI) adalah pola desain di mana dependency suatu komponen disediakan dari luar, bukan dibuat di dalam komponen itu sendiri. Ini kebalikan dari tight coupling di mana komponen membuat sendiri dependency-nya.

### 6.2 Perbandingan Pendekatan

| Pendekatan | Contoh | Kelebihan | Kekurangan |
|---|---|---|---|
| **Tight Coupling** | `handler := NewHandler()` (handler membuat db sendiri) | Sederhana | Tidak testable, sulit ganti dependency |
| **Dependency Injection** | `handler := NewHandler(db, engine, hub)` | Testable, loosely coupled | Lebih banyak parameter |

Gamon memilih DI karena:

1. **Testability** — Dalam pengujian, `*sql.DB` dapat diganti dengan database in-memory.
2. **Single Source of Truth** — Hanya ada satu instance `*sql.DB` yang dibagikan ke semua komponen.
3. **Loosely Coupled** — Handler tidak perlu tahu cara membuat database atau engine.

### 6.3 Contoh DI dalam Gamon

```go
// main.go — Komponen dibuat di sini
db, _ := database.NewDB()
hub := handler.NewHub(db)
engine := monitor.NewEngine(hub, db)

// Dependency di-inject ke handler
deviceHandler := handler.NewDeviceHandler(db, engine, hub)
alertHandler := handler.NewAlertHandler(db)
dashboardHandler := handler.NewDashboardHandler(db)
monitoringHandler := handler.NewMonitoringHandler(db)
```

**Diagram aliran dependency:**

```text
main.go membuat:
    │
    ├── db (*sql.DB)
    │     │
    │     ├──► di-inject ke Hub
    │     ├──► di-inject ke Engine
    │     ├──► di-inject ke DeviceHandler
    │     ├──► di-inject ke AlertHandler
    │     ├──► di-inject ke DashboardHandler
    │     └──► di-inject ke MonitoringHandler
    │
    ├── hub (*Hub)
    │     │
    │     ├──► di-inject ke Engine
    │     ├──► di-inject ke DeviceHandler
    │     └──► di-pass ke WebSocket handler
    │
    └── engine (*Engine)
          │
          ├──► di-inject ke DeviceHandler
          └──► di-inject ke API (legacy)
```

---

## 7. Registrasi Route

### 7.1 http.ServeMux

Go menggunakan `http.ServeMux` sebagai HTTP request multiplexer bawaan. Mux mencocokkan URL path dari request yang masuk dengan pattern yang sudah didaftarkan, lalu memanggil handler yang sesuai.

```go
mux := http.NewServeMux()

mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
    handler.HandleWebSocket(hub, w, r)
})

mux.HandleFunc("/api/devices", deviceHandler.HandleDevices)
mux.HandleFunc("/api/devices/", deviceHandler.HandleDevice)
mux.HandleFunc("/api/alerts", alertHandler.HandleAlerts)
mux.HandleFunc("/api/alerts/", alertHandler.HandleAlert)
mux.HandleFunc("/api/dashboard", dashboardHandler.HandleDashboard)
mux.HandleFunc("/api/monitoring", monitoringHandler.HandleMonitoring)
mux.HandleFunc("/api/monitoring/", monitoringHandler.HandleMonitoringDevice)
mux.HandleFunc("/api/health", legacyAPI.Health)
```

### 7.2 Pola Trailing Slash

Gamon mendaftarkan dua pola untuk setiap resource:

| Pattern | Perilaku | Contoh URL yang Cocok |
|---|---|---|
| `/api/devices` | Exact match | `/api/devices` saja |
| `/api/devices/` | Prefix match | `/api/devices/1`, `/api/devices/1/start` |

**Mengapa dua pola?**

`http.ServeMux` membedakan antara exact match (tanpa trailing slash) dan prefix match (dengan trailing slash). Pola tanpa trailing slash menangani operasi collection (list, create), sedangkan pola dengan trailing slash menangani operasi single item (get, update, delete, sub-routes).

### 7.3 URL Path Extraction

Karena `http.ServeMux` tidak mendukung path parameter secara native, setiap handler mengekstrak ID dari URL secara manual:

```go
// Input: /api/devices/5/start
idStr := strings.TrimPrefix(r.URL.Path, "/api/devices/")  // → "5/start"
idStr = strings.Split(idStr, "/")[0]                        // → "5"
id, _ := strconv.Atoi(idStr)                               // → 5
```

### 7.4 Sub-Route Dispatch

Setelah ID diekstrak, handler memeriksa sisa URL path:

```go
parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/devices/"+idStr), "/")

if len(parts) > 1 && parts[1] == "start" {
    h.startMonitoring(w, r, id)
    return
}
if len(parts) > 1 && parts[1] == "stop" {
    h.stopMonitoring(w, r, id)
    return
}
```

---

## 8. CORS Middleware

### 8.1 Pengertian CORS

CORS (Cross-Origin Resource Sharing) adalah mekanisme keamanan browser yang membatasi akses antar origin berbeda. Dalam pengembangan Gamon, frontend berjalan di port `:5173` (Vite) sedangkan backend di port `:8080`. Karena port berbeda, browser menganggap ini cross-origin.

### 8.2 Implementasi

```go
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### 8.3 Alur Request dengan CORS

```text
Frontend mengirim request
        │
        ▼
Browser: "Ini cross-origin? Kirim OPTIONS preflight dulu"
        │
        ▼
OPTIONS /api/devices
        │
        ▼
CORS Middleware: "Set headers, return 200 OK"
        │
        ▼
Browser: "Oke, sekarang kirim request sebenarnya"
        │
        ▼
GET /api/devices
        │
        ▼
CORS Middleware: "Set headers, lanjutkan ke handler"
        │
        ▼
DeviceHandler: "Proses dan kirim response"
```

---

## 9. Auto-Start Monitoring

### 9.1 Pengertian

Auto-start monitoring adalah mekanisme di mana server secara otomatis memulai monitoring untuk semua perangkat yang berstatus `active` saat pertama kali dijalankan.

### 9.2 Implementasi

```go
go autoStartMonitoring(db, engine)
```

Fungsi ini dijalankan sebagai goroutine dengan delay 2 detik:

```go
func autoStartMonitoring(db *sql.DB, engine *monitor.Engine) {
    time.Sleep(2 * time.Second)  // Tunggu server siap

    rows, _ := db.Query("SELECT id, ip, method, url, port, check_interval FROM devices WHERE status = 'active'")
    defer rows.Close()

    for rows.Next() {
        var id, interval int
        var ip, method, url string
        var port *int
        rows.Scan(&id, &ip, &method, &url, &port, &interval)

        config := monitor.DeviceConfig{DeviceID: id, IP: ip, URL: url, Method: method, Interval: interval}
        if port != nil {
            config.Port = *port
        }
        engine.Start(config)
    }
}
```

### 9.3 Mengapa Delay 2 Detik?

| Tanpa Delay | Dengan Delay |
|---|---|
| Auto-start berbarengan dengan server start | Server sudah siap menerima WebSocket client |
| Race condition potential | Tidak ada race condition |
| Client mungkin belum terhubung | Client sudah bisa menerima initial_state |

Delay 2 detik memberikan waktu bagi HTTP server untuk mulai mendengarkan dan WebSocket Hub untuk memulai goroutine-nya.

### 9.4 Idempotent Check

`engine.Start()` memiliki mekanisme idempotent — jika device sudah dimonitor, panggilan berikutnya akan diabaikan:

```go
func (e *Engine) Start(config DeviceConfig) {
    e.mu.Lock()
    if _, exists := e.targets[config.DeviceID]; exists {
        e.mu.Unlock()
        return  // Sudah dimonitor, skip
    }
    // ... mulai monitoring ...
}
```

---

## 10. Lifecycle Aplikasi

### 10.1 Diagram Lifecycle

```text
┌─────────────────────────────────────────────────────────────┐
│                    LIFECYCLE APLIKASI                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐                                            │
│  │   STARTUP    │  main() mulai dieksekusi                  │
│  │              │                                            │
│  │ 1. Init DB   │  database.NewDB()                         │
│  │ 2. Init Hub  │  handler.NewHub(db) + go hub.Run()        │
│  │ 3. Init Eng  │  monitor.NewEngine(hub, db)               │
│  │ 4. Init Hdlr │  handler.New*(db, engine, hub)            │
│  │ 5. Reg Route │  mux.HandleFunc(...)                      │
│  │ 6. Wrap CORS │  corsMiddleware(mux)                      │
│  │ 7. Auto-Start│  go autoStartMonitoring(db, engine)       │
│  │ 8. Listen    │  http.ListenAndServe(":8080", wrapped)    │
│  └──────┬───────┘                                            │
│         │                                                    │
│         ▼                                                    │
│  ┌──────────────┐                                            │
│  │   RUNNING    │  Server melayani request                   │
│  │              │                                            │
│  │  • Accept    │  HTTP requests dari frontend               │
│  │  • WebSocket │  Koneksi real-time dari client             │
│  │  • Monitor   │  Engine menjalankan ping loop              │
│  │  • Broadcast │  Hub mengirim update ke semua client       │
│  └──────┬───────┘                                            │
│         │                                                    │
│         │  Ctrl+C / SIGTERM                                  │
│         ▼                                                    │
│  ┌──────────────┐                                            │
│  │   SHUTDOWN   │  Graceful shutdown                         │
│  │              │                                            │
│  │  • db.Close()│  Tutup koneksi database                    │
│  │  • Stop eng  │  Hentikan semua monitoring loop            │
│  │  • Close WS  │  Tutup semua koneksi WebSocket             │
│  └──────────────┘                                            │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### 10.2 Graceful Shutdown (Belum Diimplementasikan)

Saat ini, Gamon menggunakan `log.Fatal(http.ListenAndServe(...))` yang akan menghentikan program secara langsung jika server error. Belum ada mekanisme graceful shutdown yang menunggu semua request selesai sebelum berhenti.

**Pengembangan lanjutan:**

```go
// Contoh graceful shutdown (belum diimplementasikan)
srv := &http.Server{Addr: ":8080", Handler: wrapped}
go func() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    srv.Shutdown(ctx)
    db.Close()
}()
log.Fatal(srv.ListenAndServe())
```

---

## 11. Struktur File Final

### 11.1 Peta File Gamon

```
/home/aby/gamon/
├── main.go                    # Entry point + wiring (Fase 6)
├── go.mod                     # Module definition
├── go.sum                     # Dependency checksums
│
├── database/
│   ├── db.go                  # Koneksi + migrasi (Fase 1)
│   └── models.go              # Struct definitions (Fase 1)
│
├── handler/
│   ├── helpers.go             # Response helpers
│   ├── device.go              # CRUD device + start/stop (Fase 2)
│   ├── alert.go               # Alert CRUD (Fase 2)
│   ├── dashboard.go           # Dashboard summary (Fase 2)
│   ├── monitoring.go          # Monitoring status + history (Fase 2)
│   ├── websocket.go           # WebSocket Hub + broadcast (Fase 2)
│   ├── api.go                 # Legacy health check
│   └── device_test.go         # Handler tests (Fase 5)
│
├── monitor/
│   ├── engine.go              # Monitoring engine (Fase 3)
│   ├── ping.go                # ICMP ping checker (Fase 3)
│   └── engine_test.go         # Engine tests (Fase 3)
│
├── data/
│   └── gamon.db               # SQLite database (otomatis dibuat)
│
└── frontend/                  # React application
    ├── src/
    │   ├── App.tsx            # Root component + routing
    │   ├── main.tsx           # React entry point
    │   ├── types/index.ts     # TypeScript types
    │   ├── hooks/
    │   │   └── useWebSocket.ts # WebSocket hook (Fase 5)
    │   ├── lib/
    │   │   ├── api.ts         # REST API client (Fase 4)
    │   │   └── presenters.ts  # Data mapper (Fase 4)
    │   ├── pages/
    │   │   ├── DashboardPage.tsx
    │   │   ├── DeviceManagementPage.tsx
    │   │   ├── MonitoringPage.tsx
    │   │   └── AlertCenterPage.tsx
    │   └── components/
    │       ├── Sidebar.tsx
    │       ├── TopBar.tsx
    │       ├── DeviceFormModal.tsx
    │       └── ... (20+ components)
    └── package.json
```

### 11.2 Kontribusi per Fase

| Fase | File Baru | File Dimodifikasi |
|---|---|---|
| 1. Database Layer | `database/db.go`, `database/models.go` | — |
| 2. Backend API | `handler/device.go`, `alert.go`, `dashboard.go`, `monitoring.go`, `websocket.go`, `helpers.go` | — |
| 3. Engine Enhancement | `monitor/engine.go`, `monitor/ping.go` | — |
| 4. Frontend API Layer | `frontend/src/lib/api.ts`, `presenters.ts` | — |
| 5. Frontend Integration | `frontend/src/pages/MonitoringPage.tsx` | `types/index.ts`, `useWebSocket.ts`, `App.tsx`, `DashboardPage.tsx`, `DeviceManagementPage.tsx`, `AlertCenterPage.tsx`, `DeviceFormModal.tsx`, `Sidebar.tsx`, `handler/websocket.go` |
| 6. Wiring | — | `main.go` (sudah lengkap) |

---

## 12. Pengujian

### 12.1 Smoke Test

Untuk memastikan wiring berfungsi, jalankan aplikasi dan lakukan smoke test:

```bash
# Build dan jalankan backend
go run .

# Jalankan frontend (terminal terpisah)
cd frontend && npm run dev
```

**Checklist smoke test:**

| No | Aksi | Endpoint/Halaman | Hasil yang Diharapkan |
|---|---|---|---|
| 1 | Buka browser | `http://localhost:5173` | Halaman Dashboard muncul |
| 2 | Cek health | `GET /api/health` | `{"status": "ok"}` |
| 3 | Ambil devices | `GET /api/devices` | `{"success": true, "data": [...]}` |
| 4 | Ambil dashboard | `GET /api/dashboard` | `{"success": true, "data": {summary, latest_alerts}}` |
| 5 | Ambil monitoring | `GET /api/monitoring` | `{"success": true, "data": [...]}` |
| 6 | Ambil alerts | `GET /api/alerts` | `{"success": true, "data": [...]}` |
| 7 | WebSocket connect | `ws://localhost:8080/ws` | Menerima `initial_state` |
| 8 | Tambah device | `POST /api/devices` | Device tersimpan, monitoring mulai |
| 9 | Hapus device | `DELETE /api/devices/{id}` | Device terhapus, monitoring berhenti |
| 10 | Lihat status live | Monitoring Page | Status update real-time |

### 12.2 Verifikasi Endpoint

```bash
# Health check
curl http://localhost:8080/api/health

# List devices
curl http://localhost:8080/api/devices

# Dashboard
curl http://localhost:8080/api/dashboard

# Monitoring
curl http://localhost:8080/api/monitoring
```

---

## 13. Kesimpulan

### 13.1 Prinsip Desain Wiring Gamon

1. **Separation of Concerns** — `main.go` hanya bertanggung jawab untuk wiring, bukan logika bisnis. Setiap komponen memiliki file dan package tersendiri.

2. **Dependency Injection** — Seluruh dependency disuntikkan melalui constructor, bukan dibuat di dalam komponen. Ini membuat aplikasi testable dan loosely coupled.

3. **Fail-Fast** — Jika database tidak bisa dibuka, program langsung berhenti (`log.Fatalf`). Tidak ada gunanya menjalankan server tanpa database.

4. **Concurrency** — Hub dan Engine dijalankan sebagai goroutine terpisah agar tidak memblokir HTTP server.

5. **Convention over Configuration** — Tidak ada file konfigurasi eksternal. Semua konfigurasi (port, database path, interval) di-hardcode untuk kesederhanaan MVP.

### 13.2 Capaian MVP

Dengan selesainya Fase 6, Gamon telah mencapai target MVP sebagai **full-stack working system**:

| Target MVP | Status |
|---|---|
| Semua halaman menggunakan data real dari database | ✅ |
| Device Management bisa CRUD beneran | ✅ |
| Monitoring real-time (ping jalan, status update live) | ✅ |
| Alert otomatis terbuat saat device down/recover | ✅ |
| 4 halaman berfungsi: Dashboard, Device Management, Monitoring, Alert Center | ✅ |
| Sistem fleksibel untuk multi-method monitoring | ✅ |

### 13.3 Poin Jawaban untuk Sidang

| Pertanyaan penguji | Jawaban ringkas |
|---|---|
| Apa itu wiring dalam konteks Gamon? | Proses menyambung seluruh komponen (database, handler, engine, middleware) di main.go agar dapat berjalan sebagai satu kesatuan. |
| Mengapa menggunakan dependency injection? | Agar komponen testable (database bisa diganti mock) dan loosely coupled (handler tidak tahu cara membuat database). |
| Mengapa menggunakan http.ServeMux, bukan router eksternal? | Kebutuhan Gamon sederhana. http.ServeMux sudah mencukupi dan tidak menambah dependency. |
| Bagaimana cara kerja CORS middleware? | Menambahkan header Access-Control-Allow-Origin agar frontend di port berbeda dapat mengakses backend. |
| Apa yang terjadi saat server pertama kali dijalankan? | Database diinisialisasi, semua komponen dibuat, route didaftarkan, dan auto-start monitoring dimulai untuk semua device active. |
| Bagaimana jika database tidak bisa dibuka? | Program langsung berhenti dengan log.Fatalf (fail-fast). Tidak ada gunanya menjalankan server tanpa database. |
| Apakah MVP sudah selesai? | Ya. Semua target tercapai: CRUD real, monitoring real-time, alert otomatis, 4 halaman berfungsi. |
| Apa yang belum diimplementasikan? | Graceful shutdown, authentication, config file, structured logging — semua di luar scope MVP. |

---

## Daftar Pustaka

[^1]: Go. *net/http Package*. https://pkg.go.dev/net/http (diakses 30 Juli 2026).

[^2]: Go. *database/sql Package*. https://pkg.go.dev/database/sql (diakses 30 Juli 2026).

[^3]: Gorilla WebSocket. *Introduction*. https://pkg.go.dev/github.com/gorilla/websocket (diakses 30 Juli 2026).

[^4]: MDN Web Docs. *CORS*. https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS (diakses 30 Juli 2026).

[^5]: Martin, R. C. *Clean Architecture*. Prentice Hall, 2017.

---

*Dokumentasi ini merupakan bagian dari laporan Praktek Kerja Lapangan (PKL) Gamon — Garda Monitoring.*
