# Backend API Layer — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian REST API

REST (Representational State Transfer) API merupakan antarmuka pemrograman aplikasi berbasis web yang menggunakan metode HTTP standar untuk pertukaran data antara client dan server. Dalam arsitektur REST, setiap sumber daya (resource) diidentifikasi oleh URL unik, dan operasi dilakukan melalui method HTTP: `GET` (membaca), `POST` (membuat), `PUT` (memperbarui), dan `DELETE` (menghapus).

REST API menjadi standar de facto dalam pengembangan aplikasi web modern karena beberapa keunggulan:

| Keunggulan | Penjelasan |
|---|---|
| **Stateless** | Setiap request berdiri sendiri tanpa menyimpan status session di server |
| **Scalable** | Mudah diskalakan karena tidak ada ketergantungan state antar request |
| **Platform Independent** | Client dapat menggunakan bahasa pemrograman apa pun selama menggunakan HTTP |
| **Loosely Coupled** | Frontend dan backend dapat dikembangkan secara terpisah |
| **Standardized** | Menggunakan HTTP method yang sudah menjadi standar industri |

### 1.2 Peran Backend API dalam Gamon

Backend API dalam Gamon berfungsi sebagai jembatan antara frontend (React) dan database (SQLite). Seluruh interaksi pengguna — mulai dari melihat daftar perangkat, menambah perangkat baru, memeriksa status monitoring, hingga mengelola alert — dilakukan melalui permintaan (request) ke endpoint API.

Peran utama Backend API dalam Gamon:

| Peran | Penjelasan |
|---|---|
| **Menyediakan Data** | Mengambil data dari database dan mengirimkannya ke frontend dalam format JSON |
| **Memproses Input** | Menerima input dari pengguna, memvalidasi, dan menyimpan ke database |
| **Mengelola Monitoring** | Mengontrol engine monitoring (start/stop) berdasarkan aksi pengguna |
| **Menangani Real-time** | Menyediakan data real-time melalui WebSocket untuk pembaruan status |
| **Enforcing Business Logic** | Menerapkan aturan bisnis seperti validasi data dan otorisasi aksi |

### 1.3 Mengapa REST API Diperlukan

Sebelum adanya Backend API yang terstruktur,前端 dan backend Gamon berkomunikasi secara langsung melalui WebSocket dengan format pesan yang tidak terstandarisasi. Pendekatan ini memiliki beberapa keterbatasan:

| Aspek | Tanpa REST API | Dengan REST API |
|---|---|---|
| Struktur data | Format pesan bebas, sulit di-maintain | Format JSON terstandarisasi |
| HTTP methods | Hanya WebSocket (tanpa GET/POST/PUT/DELETE) | Menggunakan HTTP method yang tepat |
| CRUD operations | Harus implementasi manual di WebSocket | Endpoint khusus untuk setiap operasi |
| Error handling | Tidak terstruktur | Response error konsisten |
| Testing | Sulit diuji | Dapat diuji dengan curl/Postman |
| Caching | Tidak bisa di-cache | HTTP caching tersedia |

---

## 2. Pemilihan Pendekatan

### 2.1 Mengapa http.ServeMux (Bukan Router Eksternal)

Dalam ekosistem Go, terdapat beberapa pilihan router eksternal seperti Gorilla Mux, Chi Router, atau Echo Framework. Namun, Gamon memilih untuk menggunakan `http.ServeMux` yang merupakan router bawaan Go.

Evaluasi terhadap beberapa opsi router:

| Router | Kelebihan | Kekurangan | Kesesuaian dengan Gamon |
|---|---|---|---|
| **http.ServeMux** | Built-in, zero dependency, ringan | Fitur terbatas (path parameter) | Sangat sesuai untuk API sederhana |
| **Gorilla Mux** | Path parameter, middleware support | Dependensi tambahan | Terlalu berat untuk kebutuhan Gamon |
| **Chi Router** | Lightweight, compatible with std | Dependensi tambahan | Fitur yang tidak dimanfaatkan |
| **Echo Framework** | Full-featured, performa tinggi | Overhead besar | Terlalu berat untuk prototype |

**Alasan pemilihan http.ServeMux:**

1. **Zero Dependency** — `http.ServeMux` merupakan bagian dari standard library Go. Tidak perlu menambahkan dependency eksternal, yang menjaga ukuran binary tetap kecil dan dependensi minimal.

2. **Sudah Mencukupi untuk Gamon** — API Gamon memiliki struktur yang relatif sederhana: CRUD perangkat, alert management, dan monitoring status. Tidak membutuhkan fitur canggih seperti path parameter atau regex routing.

3. **Mudah Dipahami** — Karena merupakan bagian dari standard library, dokumentasi dan contoh penggunaannya sangat mudah ditemukan. Tim pengembang dapat langsung bekerja tanpa perlu belajar API router baru.

4. **Tidak Over-Engineered** — Menggunakan framework yang terlalu besar untuk kebutuhan sederhana merupakan bentuk over-engineering. `http.ServeMux` memberikan solusi yang tepat ukuran (right-sized) untuk kebutuhan Gamon.

### 2.2 Mengapa Raw database/sql (Bukan ORM)

Untuk berinteraksi dengan database, Gamon menggunakan package `database/sql` secara langsung tanpa ORM (Object-Relational Mapping) seperti GORM atau Ent.

Evaluasi terhadap beberapa opsi akses database:

| Pendekatan | Kelebihan | Kekurangan | Kesesuaian dengan Gamon |
|---|---|---|---|
| **Raw database/sql** | Full control, performa tinggi, transparan | Lebih banyak boilerplate | Sangat sesuai untuk prototype |
| **GORM** | CRUD otomatis, migration tool | Overhead besar, abstraksi tebal | Terlalu berat untuk kebutuhan |
| **Ent** | Type-safe, code generation | Learning curve tinggi | Kompleks untuk prototype |
| **sqlx** | Helper functions, tetap fleksibel | Dependensi tambahan | Tidak signifikan bedanya |

**Alasan pemilihan raw database/sql:**

1. **Full Control** — Dengan menggunakan `database/sql` langsung, pengembang memiliki kontrol penuh terhadap setiap query SQL yang dijalankan. Tidak ada abstraksi yang menyembunyikan perilaku database.

2. **Transparansi** — Setiap query dapat dilihat, dipahami, dan dimodifikasi secara eksplisit. Sangat berguna untuk debugging dan optimasi performa.

3. **Ringan** — Tidak ada dependency tambahan yang perlu di-load. Cocok untuk aplikasi yang membutuhkan overhead minimal.

4. **Belajar SQL yang Sebenarnya** — Dengan tidak menggunakan ORM, pengembang dipaksa untuk memahami SQL secara langsung. Ini merupakan keunggulan untuk proyek pembelajaran seperti PKL.

5. **Sudah Mencukupi** — Jumlah query SQL dalam Gamon tidak terlalu banyak (sekitar 15-20 query). Boilerplate tambahan dari raw SQL masih dapat ditoleransi.

### 2.3 Mengapa Struct-Per-Resource Pattern

Gamon menerapkan pola desain di mana setiap resource (perangkat, alert, monitoring) memiliki struct handler tersendiri. Pola ini dipilih karena beberapa alasan:

| Pola | Kelebihan | Kekurangan |
|---|---|---|
| **Struct-per-resource** | Organisasi kode jelas, mudah di-maintain | Lebih banyak file |
| **Single handler file** | Semua kode di satu tempat | File besar, sulit navigasi |
| **Functional approach** | Simple, tanpa struct | Sulit manage state dan dependency |

**Alasan pemilihan struct-per-resource:**

1. **Separation of Concerns** — Setiap resource memiliki handler tersendiri, sehingga kode untuk mengelola perangkat terpisah dari kode untuk mengelola alert. Hal ini memudahkan debugging dan maintenance.

2. **Dependency Injection** — Setiap struct handler menyimpan dependency yang dibutuhkannya (database, engine, hub) sebagai field. Dependency di-inject melalui constructor, bukan sebagai parameter global.

3. **Mudah Di-extend** — Jika suatu saat ditambahkan resource baru (misal: user management), cukup buat file handler baru tanpa mengganggu handler yang sudah ada.

4. **Testable** — Setiap handler dapat diuji secara terpisah dengan dependency yang di-mock.

---

## 3. Arsitektur Backend API

### 3.1 Posisi dalam Arsitektur Gamon

Backend API menempati posisi tengah dalam arsitektur Gamon, menjadi penghubung antara frontend (React) dan backend services (Database, Engine, WebSocket):

```
┌─────────────────────────────────────────────────────────┐
│                     BROWSER (Client)                    │
│                   React + TailwindCSS                   │
│                     Port :5173                          │
│                                                         │
│  ┌──────────┐   ┌──────────────┐   ┌────────────────┐  │
│  │ Dashboard│   │ Device Mgmt  │   │ Alert Center   │  │
│  │ Page     │   │ Page         │   │ Page           │  │
│  └────┬─────┘   └──────┬───────┘   └───────┬────────┘  │
│       │                │                    │           │
│       └────────────────┼────────────────────┘           │
│                        │                                │
│                  fetch() / axios                        │
│                  HTTP Requests                          │
└────────────────────────┼────────────────────────────────┘
                         │
                         │ HTTP (GET/POST/PUT/DELETE)
                         │
┌────────────────────────┼────────────────────────────────┐
│                        ▼                                │
│            ┌──────────────────────┐                     │
│            │    GORILLA/MUX?      │                     │
│            │   http.ServeMux      │ ◄── Router          │
│            │   (URL Dispatch)     │                     │
│            └──────────┬───────────┘                     │
│                       │                                 │
│         ┌─────────────┼─────────────┐                   │
│         │             │             │                   │
│         ▼             ▼             ▼                   │
│  ┌────────────┐ ┌───────────┐ ┌─────────────┐          │
│  │  Device    │ │  Alert    │ │ Dashboard   │          │
│  │  Handler   │ │  Handler  │ │ Handler     │          │
│  │            │ │           │ │             │          │
│  │  CRUD +    │ │  List +   │ │ Summary +   │          │
│  │  Start/Stop│ │  Resolve  │ │ Latest      │          │
│  └──────┬─────┘ └─────┬─────┘ └──────┬──────┘          │
│         │             │              │                  │
│         └─────────────┼──────────────┘                  │
│                       │                                 │
│                       ▼                                 │
│              ┌──────────────────┐                       │
│              │   DATABASE LAYER │                       │
│              │   (SQLite)       │                       │
│              └──────────────────┘                       │
│                                                         │
│                       │                                 │
│         ┌─────────────┼─────────────┐                   │
│         │             │             │                   │
│         ▼             ▼             ▼                   │
│  ┌────────────┐ ┌───────────┐ ┌─────────────┐          │
│  │  Monitor   │ │ WebSocket │ │  CORS       │          │
│  │  Engine    │ │ Hub       │ │  Middleware  │          │
│  │            │ │           │ │             │          │
│  │  Start/Stop│ │  Broadcast│ │  Allow      │          │
│  │  Ping      │ │  Status   │ │  Origins    │          │
│  └────────────┘ └───────────┘ └─────────────┘          │
│                                                         │
│                      GO BACKEND                         │
│                      Port :8080                         │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Komponen Backend API

Backend API Gamon terdiri dari beberapa komponen utama:

| Komponen | File | Fungsi |
|---|---|---|
| **Response Helpers** | `handler/helpers.go` | Fungsi bantu untuk response JSON |
| **Device Handler** | `handler/device.go` | CRUD perangkat + start/stop monitoring |
| **Alert Handler** | `handler/alert.go` | List, detail, dan resolve alert |
| **Dashboard Handler** | `handler/dashboard.go` | Ringkasan dashboard |
| **Monitoring Handler** | `handler/monitoring.go` | Status monitoring + riwayat |
| **Health API** | `handler/api.go` | Health check endpoint |
| **CORS Middleware** | `main.go` | Middleware untuk Cross-Origin Resource Sharing |
| **Router** | `main.go` | Registrasi route ke handler |

### 3.3 Dependency Graph

```
main.go
  │
  ├──► database.NewDB() ──► *sql.DB
  │
  ├──► monitor.NewEngine(hub) ──► *Engine
  │
  ├──► handler.NewHub() ──► *Hub
  │
  ├──► handler.NewDeviceHandler(db, engine, hub) ──► *DeviceHandler
  │
  ├──► handler.NewAlertHandler(db) ──► *AlertHandler
  │
  ├──► handler.NewDashboardHandler(db) ──► *DashboardHandler
  │
  └──► handler.NewMonitoringHandler(db) ──► *MonitoringHandler
```

---

## 4. Struktur File dan Package

### 4.1 Package `handler`

Package `handler` berisi seluruh logika untuk menangani request HTTP. Setiap file dalam package ini memiliki tanggung jawab spesifik:

#### `handler/helpers.go` — Response Helpers

File ini berisi fungsi-fungsi bantu untuk mengirim response JSON. Dengan adanya helper functions, seluruh handler menggunakan format response yang seragam.

**Fungsi yang tersedia:**

| Fungsi | Parameter | HTTP Status | Fungsi |
|---|---|---|---|
| `respondJSON` | `w, status, payload` | Caller-specified | Mengirim JSON response dengan status code tertentu |
| `respondError` | `w, status, message` | Caller-specified | Mengirim error response dengan format `{"success": false, "message": "..."}` |
| `respondSuccess` | `w, message` | 200 OK | Mengirim success response dengan format `{"success": true, "message": "..."}` |
| `respondData` | `w, data` | 200 OK | Mengirim data response dengan format `{"success": true, "data": {...}}` |

#### `handler/device.go` — Device Handler

File ini menangani seluruh operasi CRUD untuk resource perangkat (device), termasuk start/stop monitoring.

**Struct:**
```go
type DeviceHandler struct {
    db     *sql.DB
    engine *monitor.Engine
    hub    *Hub
}
```

**Method yang tersedia:**

| Method | Fungsi | HTTP Methods |
|---|---|---|
| `HandleDevices` | List semua device atau tambah device baru | GET, POST |
| `HandleDevice` | Get, update, hapus device, atau start/stop monitoring | GET, PUT, DELETE, POST |

#### `handler/alert.go` — Alert Handler

File ini menangani operasi terkait alert (gangguan), termasuk list, detail, dan resolve.

**Struct:**
```go
type AlertHandler struct {
    db *sql.DB
}
```

**Method yang tersedia:**

| Method | Fungsi | HTTP Methods |
|---|---|---|
| `HandleAlerts` | List semua alert dengan filter | GET |
| `HandleAlert` | Get detail alert atau resolve alert | GET, POST |

#### `handler/dashboard.go` — Dashboard Handler

File ini menangani pengambilan data ringkasan untuk dashboard.

**Struct:**
```go
type DashboardHandler struct {
    db *sql.DB
}
```

**Method yang tersedia:**

| Method | Fungsi | HTTP Methods |
|---|---|---|
| `HandleDashboard` | Ambil ringkasan dashboard (total device, online/offline, alert terbaru) | GET |

#### `handler/monitoring.go` — Monitoring Handler

File ini menangani data status monitoring dan riwayat pengecekan.

**Struct:**
```go
type MonitoringHandler struct {
    db *sql.DB
}
```

**Method yang tersedia:**

| Method | Fungsi | HTTP Methods |
|---|---|---|
| `HandleMonitoring` | List status semua device yang sedang dimonitor | GET |
| `HandleMonitoringDevice` | Ambil riwayat pengecekan device tertentu | GET |

#### `handler/api.go` — Health API (Legacy)

File ini berisi endpoint health check yang merupakan sisa dari implementasi sebelumnya.

**Struct:**
```go
type API struct {
    engine *monitor.Engine
    hub    *Hub
}
```

**Method yang tersedia:**

| Method | Fungsi | HTTP Methods |
|---|---|---|
| `Health` | Health check endpoint | GET |

---

## 5. Response Format Standard

### 5.1 Mengapa Response Format Diperlukan

Dalam pengembangan API, konsistensi format response merupakan hal yang sangat penting. Tanpa format yang seragam, frontend harus menangani berbagai macam format response yang berbeda-beda, yang menyebabkan kode yang rumit dan sulit di-maintain.

Gamon menerapkan dua jenis response envelope:

### 5.2 APIResponse — Success/Error Confirmation

Digunakan untuk konfirmasi operasi yang tidak mengembalikan data (create, update, delete, resolve, start/stop).

**Format:**
```json
{
    "success": true,
    "message": "Device deleted successfully"
}
```
```json
{
    "success": false,
    "message": "Invalid device ID"
}
```

### 5.3 DataResponse — Data Payload

Digunakan untuk response yang mengembalikan data (list, detail, summary).

**Format:**
```json
{
    "success": true,
    "data": [
        {"id": 1, "name": "Server A", "ip": "192.168.1.1"},
        {"id": 2, "name": "Router B", "ip": "192.168.1.2"}
    ]
}
```

### 5.4 Helper Functions

| Fungsi | Penggunaan | Contoh Output |
|---|---|---|
| `respondJSON(w, 200, data)` | Response dengan custom status | `{"key": "value"}` |
| `respondError(w, 400, "Invalid input")` | Error response | `{"success": false, "message": "Invalid input"}` |
| `respondSuccess(w, "Deleted")` | Success tanpa data | `{"success": true, "message": "Deleted"}` |
| `respondData(w, devices)` | Success dengan data | `{"success": true, "data": [...]}` |

---

## 6. Endpoint Reference

### 6.1 Device Endpoints

| Method | Endpoint | Fungsi | Request Body | Response |
|---|---|---|---|---|
| `GET` | `/api/devices` | List semua device | - | `{"success": true, "data": [...]}` |
| `POST` | `/api/devices` | Tambah device baru | `{name, type, ip, ...}` | `{"success": true, "data": {...}}` |
| `GET` | `/api/devices/{id}` | Detail device | - | `{"success": true, "data": {...}}` |
| `PUT` | `/api/devices/{id}` | Update device | `{name, ip, ...}` | `{"success": true, "data": {...}}` |
| `DELETE` | `/api/devices/{id}` | Hapus device | - | `{"success": true, "message": "..."}` |
| `POST` | `/api/devices/{id}/start` | Mulai monitoring | - | `{"success": true, "message": "..."}` |
| `POST` | `/api/devices/{id}/stop` | Hentikan monitoring | - | `{"success": true, "message": "..."}` |

#### POST /api/devices — Request Body

| Field | Tipe | Wajib | Default | Penjelasan |
|---|---|---|---|---|
| `name` | string | Ya | - | Nama perangkat |
| `type` | string | Ya | - | Tipe: server/router/switch/ap/website |
| `ip` | string | Ya | - | Alamat IP perangkat |
| `url` | string | Tidak | `""` | URL untuk HTTP check (opsional) |
| `port` | integer | Tidak | `null` | Port untuk TCP check (opsional) |
| `method` | string | Tidak | `"ICMP Ping"` | Metode monitoring |
| `location` | string | Tidak | `""` | Lokasi perangkat (opsional) |
| `check_interval` | integer | Tidak | `3` | Interval pengecekan dalam detik |
| `description` | string | Tidak | `""` | Deskripsi perangkat |

### 6.2 Alert Endpoints

| Method | Endpoint | Fungsi | Query Params | Response |
|---|---|---|---|---|
| `GET` | `/api/alerts` | List alerts | `status`, `severity`, `device_type` | `{"success": true, "data": [...]}` |
| `GET` | `/api/alerts/{id}` | Detail alert | - | `{"success": true, "data": {...}}` |
| `POST` | `/api/alerts/{id}/resolve` | Resolve alert | - | `{"success": true, "message": "..."}` |

#### GET /api/alerts — Query Parameters

| Parameter | Tipe | Default | Penjelasan |
|---|---|---|---|
| `status` | string | Semua | Filter: ongoing/resolved |
| `severity` | string | Semua | Filter: low/medium/high/critical |
| `device_type` | string | Semua | Filter: server/router/switch/ap/website |

### 6.3 Dashboard Endpoint

| Method | Endpoint | Fungsi | Response |
|---|---|---|---|
| `GET` | `/api/dashboard` | Ringkasan dashboard | `{"success": true, "data": {summary, latest_alerts}}` |

#### Response Structure

```json
{
    "success": true,
    "data": {
        "summary": {
            "total_devices": 10,
            "online_devices": 8,
            "offline_devices": 1,
            "warning_devices": 1
        },
        "latest_alerts": [
            {
                "id": 1,
                "device_name": "Server A",
                "title": "Device Offline",
                "severity": "critical",
                "status": "ongoing",
                "started_at": "2025-01-15T10:30:00Z"
            }
        ]
    }
}
```

### 6.4 Monitoring Endpoints

| Method | Endpoint | Fungsi | Response |
|---|---|---|---|
| `GET` | `/api/monitoring` | Status semua device | `{"success": true, "data": [...]}` |
| `GET` | `/api/monitoring/{id}/history` | Riwayat pengecekan device | `{"success": true, "data": [...]}` |

### 6.5 Health Endpoint

| Method | Endpoint | Fungsi | Response |
|---|---|---|---|
| `GET` | `/api/health` | Health check | `{"status": "ok"}` |

---

## 7. Routing Strategy

### 7.1 Cara Kerja http.ServeMux

`http.ServeMux` merupakan HTTP request multiplexer bawaan Go. Ia cocokkan URL path dari permintaan yang masuk dengan pola (pattern) yang sudah didaftarkan, lalu memanggil handler yang sesuai.

**Pola registrasi route:**

| Pattern | Perilaku |
|---|---|
| `/api/devices` | Cocok tepat: hanya `/api/devices` |
| `/api/devices/` | Cocok prefix: `/api/devices`, `/api/devices/5`, `/api/devices/5/start`, dll |

### 7.2 Trailing Slash Pattern

Gamon memanfaatkan perilaku prefix matching dari `http.ServeMux` dengan mendaftarkan dua pola untuk setiap resource:

```go
mux.HandleFunc("/api/devices", deviceHandler.HandleDevices)   // Collection endpoint
mux.HandleFunc("/api/devices/", deviceHandler.HandleDevice)   // Single item endpoint
```

- `/api/devices` (tanpa trailing slash) → Menangani operasi collection (list, create)
- `/api/devices/` (dengan trailing slash) → Menangani operasi single item (get, update, delete, sub-routes)

### 7.3 URL Path Extraction

Karena `http.ServeMux` tidak mendukung path parameter secara native, setiap handler harus melakukan parsing URL secara manual:

```go
// Contoh dari DeviceHandler
idStr := strings.TrimPrefix(r.URL.Path, "/api/devices/")
idStr = strings.Split(idStr, "/")[0]
id, _ := strconv.Atoi(idStr)
```

**Proses extraction:**

```
Input:  /api/devices/5/start
Step 1: TrimPrefix → "5/start"
Step 2: Split("/") → ["5", "start"]
Step 3: Take first → "5"
Step 4: Atoi → 5
```

### 7.4 Sub-Route Dispatch

Setelah ID diekstrak, handler memeriksa sisa URL path untuk menentukan sub-route:

```go
parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/devices/"+idStr), "/")
if len(parts) > 1 && parts[1] == "start" {
    // Handle start monitoring
}
if len(parts) > 1 && parts[1] == "stop" {
    // Handle stop monitoring
}
```

**Sub-route yang didukung:**

| Sub-Route | Fungsi |
|---|---|
| `/api/devices/{id}/start` | Mulai monitoring untuk device tertentu |
| `/api/devices/{id}/stop` | Hentikan monitoring untuk device tertentu |
| `/api/alerts/{id}/resolve` | Resolve alert tertentu |
| `/api/monitoring/{id}/history` | Riwayat pengecekan device tertentu |

---

## 8. Database Integration

### 8.1 Dependency Injection

Seluruh handler menerima database connection (`*sql.DB`) melalui constructor. Pendekatan ini disebut dependency injection, di mana dependency "diinjeksikan" dari luar bukan dibuat di dalam.

```go
// Di main.go
db, err := database.NewDB()
deviceHandler := handler.NewDeviceHandler(db, engine, hub)

// Di handler
type DeviceHandler struct {
    db     *sql.DB       // ← Database connection
    engine *monitor.Engine
    hub    *Hub
}

func NewDeviceHandler(db *sql.DB, engine *monitor.Engine, hub *Hub) *DeviceHandler {
    return &DeviceHandler{db: db, engine: engine, hub: hub}
}
```

**Keuntungan dependency injection:**

1. **Testable** — Dalam pengujian, `*sql.DB` dapat diganti dengan mock database.
2. **Loosely Coupled** — Handler tidak perlu tahu cara membuat koneksi database.
3. **Single Source of Truth** — Hanya ada satu koneksi database yang dibagikan ke seluruh handler.

### 8.2 Query Patterns

`database/sql` menyediakan beberapa fungsi untuk menjalankan query:

#### Multi-Row Read — `db.Query()`

Digunakan untuk membaca multiple rows (list devices, list alerts):

```go
rows, err := h.db.Query("SELECT id, name, type, ip FROM devices")
if err != nil {
    respondError(w, 500, "Failed to fetch devices")
    return
}
defer rows.Close()

var devices []Device
for rows.Next() {
    var d Device
    err := rows.Scan(&d.ID, &d.Name, &d.Type, &d.IP)
    if err != nil {
        continue
    }
    devices = append(devices, d)
}
```

#### Single-Row Read — `db.QueryRow()`

Digunakan untuk membaca satu row (get device by ID, get alert detail):

```go
var device Device
err := h.db.QueryRow("SELECT id, name, ip FROM devices WHERE id = ?", id).
    Scan(&device.ID, &device.Name, &device.IP)
if err == sql.ErrNoRows {
    respondError(w, 404, "Device not found")
    return
}
```

#### Write Operations — `db.Exec()`

Digunakan untuk INSERT, UPDATE, DELETE:

```go
result, err := h.db.Exec(
    "INSERT INTO devices (name, type, ip) VALUES (?, ?, ?)",
    name, deviceType, ip,
)
if err != nil {
    respondError(w, 500, "Failed to create device")
    return
}
id, _ := result.LastInsertId()
```

### 8.3 Parameterized Queries (SQL Injection Prevention)

Seluruh query dalam Gamon menggunakan parameterized queries (parameter `?`) bukan string concatenation:

```go
// ✅ BENAR — Parameterized query
h.db.Query("SELECT * FROM devices WHERE id = ?", id)

// ❌ SALAH — SQL Injection vulnerable
h.db.Query("SELECT * FROM devices WHERE id = " + id)
```

**Mengapa parameterized queries penting:**

| Pendekatan | Keamanan | Penjelasan |
|---|---|---|
| Parameterized query | Aman | Parameter di-escape oleh database driver |
| String concatenation | Rentan | Pengguna dapat menyisipkan SQL code berbahaya |

### 8.4 NULL Handling

Beberapa kolom di database memiliki nilai NULL. Gamon menangani NULL dengan dua cara:

**1. COALESCE di SQL:**
```sql
SELECT COALESCE(ph.status, 'unknown') FROM ping_history ph WHERE ph.device_id = d.id
```

**2. Pointer Types di Go:**
```go
type UpdateDeviceRequest struct {
    Name     *string  // nil = tidak diupdate
    IP       *string  // nil = tidak diupdate
    Port     *int     // nil = tidak diupdate
}
```

---

## 9. CORS Middleware

### 9.1 Pengertian CORS

CORS (Cross-Origin Resource Sharing) merupakan mekanisme keamanan yang diterapkan oleh browser untuk membatasi akses antar origin yang berbeda. Origin terdiri dari kombinasi protocol, hostname, dan port.

**Contoh cross-origin:**

| Frontend | Backend | Status |
|---|---|---|
| `http://localhost:5173` | `http://localhost:8080` | Cross-origin (port berbeda) |
| `http://localhost:5173` | `http://localhost:5173` | Same-origin |
| `https://gamon.example.com` | `https://api.gamon.example.com` | Cross-origin (hostname berbeda) |

### 9.2 Mengapa CORS Diperlukan

Dalam pengembangan Gamon, frontend berjalan di port `:5173` (Vite dev server) sedangkan backend berjalan di port `:8080`. Karena port berbeda, browser menganggap ini sebagai cross-origin request dan akan memblokir request tersebut tanpa header CORS yang tepat.

**Tanpa CORS middleware:**

```
Browser: "Frontend di port 5173 mencoba akses backend di port 8080"
Browser: "Ini cross-origin! Request diblokir!"
```

**Dengan CORS middleware:**

```
Backend: "Saya izinkan semua origin untuk mengakses API ini"
Backend: Header: Access-Control-Allow-Origin: *
Browser: "Oke, request diizinkan!"
```

### 9.3 Konfigurasi CORS Gamon

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

| Header | Nilai | Penjelasan |
|---|---|---|
| `Access-Control-Allow-Origin` | `*` | Izinkan semua origin (development mode) |
| `Access-Control-Allow-Methods` | `GET, POST, PUT, DELETE, OPTIONS` | Izinkan HTTP method ini |
| `Access-Control-Allow-Headers` | `Content-Type` | Izinkan header Content-Type |

**Catatan:** Untuk production, `Access-Control-Allow-Origin` sebaiknya diganti dengan origin spesifik (misal: `https://gamon.example.com`) demi keamanan.

---

## 10. Flow Data

### 10.1 Flow Tambah Device

```
┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐
│ Browser │       │  CORS   │       │ Router  │       │ Device  │       │Database │
│ (React) │       │Middleware│       │ServeMux │       │ Handler │       │ (SQLite)│
└────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘
     │                 │                 │                 │                 │
     │ POST /api/devices                │                 │                 │
     │ {name, type, ip}                 │                 │                 │
     │────────────────►│                │                 │                 │
     │                 │                │                 │                 │
     │                 │►─ Check Origin ─►               │                 │
     │                 │                 │                │                 │
     │                 │►─ OPTIONS ──────►               │                 │
     │                 │  (preflight)    │                │                 │
     │                 │◄────────────────│                │                 │
     │                 │                │                │                 │
     │                 │►─ POST ────────►│               │                 │
     │                 │                 │                │                 │
     │                 │                 │►─ Match ──────►│                │
     │                 │                 │  /api/devices  │                │
     │                 │                 │                │                 │
     │                 │                 │                │►─ Validate ────►│
     │                 │                 │                │                 │
     │                 │                 │                │►─ INSERT ──────►│
     │                 │                 │                │  INTO devices   │
     │                 │                 │                │                 │
     │                 │                 │                │◄── Return ID ──│
     │                 │                 │                │                 │
     │                 │                 │                │►─ SELECT ──────►│
     │                 │                 │                │  (get new device)│
     │                 │                 │                │                 │
     │                 │                 │                │◄── Device data ─│
     │                 │                 │                │                 │
     │                 │                 │◄── respondData─│                 │
     │                 │                 │                │                 │
     │◄────────────────│◄────────────────│                │                 │
     │  {"success":true, "data": {...}}  │                │                 │
     │                 │                │                │                 │
```

**Penjelasan Flow:**

1. Browser mengirim `POST /api/devices` dengan body `{name, type, ip}`
2. CORS middleware memeriksa origin dan menambahkan header yang diperlukan
3. Jika request OPTIONS (preflight), langsung return 200 OK
4. `http.ServeMux` mencocokkan pattern `/api/devices` ke `DeviceHandler.HandleDevices`
5. Handler memvalidasi input (name, type, wajib diisi)
6. Handler menjalankan SQL INSERT ke tabel `devices`
7. Handler mengambil data device yang baru dibuat
8. Handler mengirim response `{"success": true, "data": {...}}`

### 10.2 Flow Ambil Dashboard Summary

```
┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐
│ Browser │       │  CORS   │       │ Router  │       │Dashboard│       │Database │
│ (React) │       │Middleware│       │ServeMux │       │ Handler │       │ (SQLite)│
└────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘
     │                 │                 │                 │                 │
     │ GET /api/dashboard               │                 │                 │
     │────────────────►│                │                 │                 │
     │                 │                │                 │                 │
     │                 │►─ GET ────────►│                │                 │
     │                 │                │                 │                 │
     │                 │                 │►─ Match ──────►│                │
     │                 │                 │  /api/dashboard│                │
     │                 │                 │                │                 │
     │                 │                 │                │►─ COUNT ──────►│
     │                 │                 │                │  total_devices │
     │                 │                 │                │                 │
     │                 │                 │                │◄── Count ──────│
     │                 │                 │                │                 │
     │                 │                 │                │►─ Subquery ───►│
     │                 │                 │                │  online/offline │
     │                 │                 │                │                 │
     │                 │                 │                │◄── Counts ─────│
     │                 │                 │                │                 │
     │                 │                 │                │►─ SELECT ──────►│
     │                 │                 │                │  latest_alerts │
     │                 │                 │                │                 │
     │                 │                 │                │◄── Alerts ─────│
     │                 │                 │                │                 │
     │                 │                 │                │►─ Build ──────►│
     │                 │                 │                │  DashboardSummary
     │                 │                 │                │                 │
     │                 │                 │◄── respondData─│                 │
     │                 │                 │                │                 │
     │◄────────────────│◄────────────────│                │                 │
     │  {"success":true, "data": {summary, latest_alerts}}                │
     │                 │                │                │                 │
```

**Penjelasan Flow:**

1. Browser mengirim `GET /api/dashboard`
2. Request melewati CORS middleware dan router
3. `DashboardHandler.HandleDashboard` dipanggil
4. Handler menjalankan beberapa query:
   - `COUNT(*) FROM devices` — total perangkat
   - Subquery untuk menghitung online/offline/warning berdasarkan ping_history terakhir
   - `SELECT` untuk mengambil 5 alert terbaru
5. Hasil digabung menjadi struct `DashboardSummary`
6. Response dikirim dengan format `{"success": true, "data": {summary, latest_alerts}}`

### 10.3 Flow Resolve Alert

```
┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐       ┌─────────┐
│ Browser │       │  CORS   │       │ Router  │       │  Alert  │       │Database │
│ (React) │       │Middleware│       │ServeMux │       │ Handler │       │ (SQLite)│
└────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘       └────┬────┘
     │                 │                 │                 │                 │
     │ POST /api/alerts/1/resolve       │                 │                 │
     │────────────────►│                │                 │                 │
     │                 │                │                 │                 │
     │                 │►─ POST ───────►│                │                 │
     │                 │                │                 │                 │
     │                 │                 │►─ Match ──────►│                │
     │                 │                 │  /api/alerts/  │                │
     │                 │                 │                │                 │
     │                 │                 │                │►─ Parse ID ───►│
     │                 │                 │                │  (1)            │
     │                 │                 │                │                 │
     │                 │                 │                │►─ Check ──────►│
     │                 │                 │                │  alert exists   │
     │                 │                 │                │                 │
     │                 │                 │                │◄── Alert found ─│
     │                 │                 │                │                 │
     │                 │                 │                │►─ UPDATE ──────►│
     │                 │                 │                │  SET status =  │
     │                 │                 │                │  'resolved'    │
     │                 │                 │                │                 │
     │                 │                 │                │◄── RowsAffected─│
     │                 │                 │                │                 │
     │                 │                 │                │►─ Check ──────►│
     │                 │                 │                │  RowsAffected > 0?
     │                 │                 │                │                 │
     │                 │                 │◄── respondSuccess               │
     │                 │                 │                │                 │
     │◄────────────────│◄────────────────│                │                 │
     │  {"success":true, "message": "Alert resolved"}    │                 │
     │                 │                │                │                 │
```

**Penjelasan Flow:**

1. Browser mengirim `POST /api/alerts/1/resolve`
2. Request melewati CORS middleware dan router
3. `AlertHandler.HandleAlert` dipanggil
4. Handler mengekstrak alert ID (1) dan path "resolve" dari URL
5. Handler memeriksa apakah alert dengan ID tersebut ada
6. Handler menjalankan SQL UPDATE: `UPDATE alerts SET status = 'resolved', resolved_at = CURRENT_TIMESTAMP WHERE id = 1 AND status = 'ongoing'`
7. Handler memeriksa `RowsAffected` — jika > 0, resolve berhasil; jika 0, alert tidak ditemukan atau sudah resolved
8. Response success/error dikirim

---

## 11. Kesimpulan

### 11.1 Prinsip Desain Backend API Gamon

Backend API layer Gamon dirancang dengan beberapa prinsip utama:

1. **Simplicity** — Menggunakan built-in Go packages (`http.ServeMux`, `database/sql`) tanpa dependency tambahan yang tidak perlu. Prinsip ini memastikan kode tetap mudah dipahami dan di-maintain.

2. **Consistency** — Seluruh endpoint menggunakan format response yang seragam (`APIResponse` dan `DataResponse`). Pengembang frontend dapat menulis kode handling response yang konsisten.

3. **Separation of Concerns** — Setiap resource memiliki handler tersendiri dengan tanggung jawab yang jelas. Device handler tidak campur aduk dengan alert handler.

4. **Dependency Injection** — Database connection di-inject melalui constructor, bukan dibuat di dalam handler. Hal ini memudahkan pengujian dan penggantian dependency.

5. **Security** — Parameterized queries digunakan untuk mencegah SQL Injection. CORS middleware dikonfigurasi untuk mengontrol akses cross-origin.

### 11.2 Capaian Fase 2

Dengan selesainya Backend API layer, Gamon kini memiliki:

| Komponen | Status | Fungsi |
|---|---|---|
| **Device CRUD** | ✅ Selesai | Create, Read, Update, Delete perangkat |
| **Start/Stop Monitoring** | ✅ Selesai | Mengontrol engine monitoring per perangkat |
| **Alert Management** | ✅ Selesai | List, detail, dan resolve alert |
| **Dashboard Summary** | ✅ Selesai | Ringkasan status untuk dashboard |
| **Monitoring Status** | ✅ Selesai | Status dan riwayat pengecekan |
| **CORS Support** | ✅ Selesai | Mengizinkan akses dari frontend |
| **Health Check** | ✅ Selesai | Endpoint untuk memeriksa status server |

### 11.3 Lanjutan ke Fase Berikutnya

Dengan backend API yang sudah berfungsi, langkah selanjutnya adalah:

1. **Phase 3 — Engine Enhancement**: Mengintegrasikan monitoring engine dengan backend API yang baru, termasuk auto-alert generation dan status tracking.

2. **Phase 4 — Frontend Integration**: Menghubungkan frontend React dengan backend API yang sebenarnya (menggantikan dummy data).

---

*Dokumentasi ini merupakan bagian dari laporan PKL (Praktek Kerja Lapangan) untuk proyek Gamon (Garda Monitoring) — Sistem Monitoring Jaringan Berbasis Web.*
