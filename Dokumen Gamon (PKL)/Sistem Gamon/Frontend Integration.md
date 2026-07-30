# Frontend Integration — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian Frontend Integration

Frontend integration merupakan proses menghubungkan antarmuka pengguna (UI) yang sebelumnya bekerja dengan data dummy atau hardcoded menjadi sebuah sistem yang secara penuh berkomunikasi dengan backend melalui REST API dan WebSocket. Dalam konteks Gamon, fase ini mengubah seluruh halaman — Dashboard, Device Management, Alert Center, dan Monitoring — dari tampilan statis menjadi aplikasi yang menampilkan data real dari SQLite, mengoperasikan CRUD yang sesungguhnya, dan memperbarui status perangkat secara real-time.

Istilah "integration" di sini merujuk pada penyatuan tiga lapisan sistem:

| Lapisan | Peran dalam Integrasi |
|---|---|
| **Frontend (React)** | Menampilkan data, merespons aksi pengguna, mengirim request |
| **Backend (Go REST API)** | Memproses permintaan, menjalankan logika bisnis, mengelola database |
| **WebSocket** | Menyalurkan data monitoring real-time dari engine ke seluruh client |

Sebelum fase ini, ketiga lapisan tersebut belum terhubung secara utuh. Frontend menampilkan data dari file dummy, backend menyediakan API tetapi belum dipanggil oleh frontend, dan WebSocket hanya mengirim data tanpa ada client yang mengelolanya secara multi-device.

### 1.2 Kondisi Sebelum Fase 5

Pada akhir Fase 4, fondasi komunikasi REST sudah tersedia di `frontend/src/lib/api.ts`, tetapi seluruh halaman React masih mengambil data dari file dummy di folder `frontend/src/data/`. Artinya, meskipun backend sudah memiliki CRUD device, alert generation, dan monitoring engine yang berjalan, pengguna tidak melihat dampaknya di antarmuka web.

| Aspek | Sebelum Fase 5 | Setelah Fase 5 |
|---|---|---|
| Sumber data halaman | File dummy/hardcoded | REST API + WebSocket real-time |
| CRUD perangkat | Frontend-only (local state) | Real CRUD via backend API + SQLite |
| Status monitoring | Dummy | Live dari WebSocket (Map per device) |
| Alert | Data statis | Dihasilkan otomatis oleh engine, diambil dari database |
| Halaman Monitoring | Tidak ada | Halaman baru dengan status live dan riwayat |
| WebSocket | Belum terpakai | Multi-device state, reconnect, re-fetch |
| Form device | Hanya ICMP Ping hardcoded | Method selector dengan field kondisional |

### 1.3 Tujuan Fase 5

Tujuan utama Fase 5 adalah mewujudkan **full-stack working system** di mana:

1. Seluruh halaman menggunakan data real dari database SQLite melalui REST API.
2. Device Management dapat melakukan CRUD beneran (create, read, update, delete) yang tersimpan di backend.
3. Monitoring berjalan secara real-time: hasil ping muncul langsung di UI tanpa refresh.
4. Alert dihasilkan secara otomatis oleh engine saat device down atau recover.
5. Empat halaman berfungsi penuh: Dashboard, Device Management, Monitoring, dan Alert Center.
6. Sistem tetap fleksibel untuk multi-method monitoring (ICMP Ping saat ini, HTTP/TCP nanti).

---

## 2. Posisi dalam Arsitektur Gamon

### 2.1 Arsitektur Full-Stack Terintegrasi

Setelah Fase 5, arsitektur Gamon mencapai bentuk targetnya:前端 React berkomunikasi dengan backend Go melalui dua jalur — REST API untuk operasi CRUD dan pengambilan data, serta WebSocket untuk pembaruan status real-time.

```text
┌─────────────────────────────────────────────────────────────┐
│                     BROWSER (Client)                         │
│                   React + TypeScript + Tailwind               │
│                     Port :5173                                │
│                                                               │
│  ┌──────────┐  ┌──────────────┐  ┌────────────┐  ┌────────┐ │
│  │ Dashboard│  │ Device Mgmt  │  │ Monitoring │  │ Alerts │ │
│  │ Page     │  │ Page         │  │ Page (new) │  │ Page   │ │
│  └────┬─────┘  └──────┬───────┘  └─────┬──────┘  └───┬────┘ │
│       │               │                │              │      │
│       │  ┌────────────┴────────────────┴──────────────┘      │
│       │  │                                                   │
│       │  ▼                                                   │
│  ┌────────────────┐     ┌─────────────────────────┐          │
│  │ api.ts         │     │ useWebSocket.ts          │          │
│  │ (REST client)  │     │ (multi-device state)     │          │
│  └───────┬────────┘     └────────────┬────────────┘          │
│          │                           │                       │
│          │ HTTP (GET/POST/PUT/DELETE) │ WebSocket             │
└──────────┼───────────────────────────┼───────────────────────┘
           │                           │
           │                           │
┌──────────▼───────────────────────────▼───────────────────────┐
│                     GO BACKEND :8080                           │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────┐  │
│  │ REST API     │  │ WebSocket    │  │ Monitoring Engine  │  │
│  │ Handlers     │  │ Hub          │  │ (ping loop)        │  │
│  │ (CRUD)       │  │ (broadcast)  │  │ (status tracking)  │  │
│  └──────┬───────┘  └──────┬───────┘  └────────┬───────────┘  │
│         │                 │                    │              │
│         └─────────────────┼────────────────────┘              │
│                           │                                   │
│                           ▼                                   │
│                  ┌──────────────────┐                         │
│                  │   SQLite DB      │                         │
│                  │   (WAL mode)     │                         │
│                  └──────────────────┘                         │
└──────────────────────────────────────────────────────────────┘
```

### 2.2 Pembagian Tanggung Jawab: REST vs WebSocket

Gamon menerapkan arsitektur hybrid di mana REST API dan WebSocket memiliki peran yang saling melengkapi. Pembagian ini bukan sekadar pilihan teknis, tetapi mencerminkan karakteristik berbeda dari data yang dipertukarkan.

| Saluran | Karakteristik | Penggunaan pada Gamon | Contoh |
|---|---|---|---|
| **REST API** | Request-response, transaksional, satu kali aksi | Operasi CRUD, pengambilan data awal, aksi user | Create device, fetch alerts, resolve alert |
| **WebSocket** | Koneksi persisten, event-driven, broadcast | Data monitoring real-time, notifikasi perubahan | Hasil ping, status change, initial state |

**Mengapa tidak semua menggunakan REST?**

REST dirancang untuk operasi request-response. Jika frontend harus meminta status device setiap 3 detik melalui REST, akan terjadi ratusan request per menit yang sebagian besar tidak menghasilkan perubahan. WebSocket mengatasi masalah ini dengan menjaga koneksi terbuka dan hanya mengirim data ketika ada perubahan.

**Mengapa tidak semua menggunakan WebSocket?**

WebSocket tidak memiliki konsep HTTP method (GET, POST, PUT, DELETE) atau status code. Operasi seperti "hapus device" atau "ambil riwayat 50 pengecekan terakhir" lebih tepat dijalankan melalui REST yang memiliki semantik method dan status code yang jelas.

### 2.3 Alur Data Real-Time

Berikut alur data dari engine hingga tampilan pengguna:

```text
Engine menjalankan ping ke device
        │
        ▼
PingOnce() menghasilkan CheckResult
(status, latency_ms, details)
        │
        ▼
Engine.saveCheckResult() → INSERT ke ping_history
        │
        ▼
Engine.trackStatus() → Deteksi perubahan status
        │
        ├─── Status berubah? → Broadcast "status_change"
        │                      Buat/resolve alert di DB
        │
        └─── Selalu → Broadcast "check_result"
                       │
                       ▼
              WebSocket Hub menerima pesan
                       │
                       ▼
              Hub.broadcast → Kirim ke semua client
                       │
                       ▼
              useWebSocket hook menerima pesan
              │
              ├── initial_state → Populate Map<device_id, MonitorResult>
              ├── check_result  → Update device terkait di Map
              └── status_change → Simpan event, trigger refresh alert
                       │
                       ▼
              React re-render dengan data baru
              │
              ├── Dashboard: update summary & metrics
              ├── Monitoring: update status device di tabel
              └── Alert Center: refresh daftar alert
```

---

## 3. Integrasi WebSocket Multi-Device

### 3.1 useWebSocket Hook

Hook `useWebSocket` merupakan pusat manajemen koneksi WebSocket di frontend. Hook ini bertanggung jawab untuk:

1. Membuka dan memelihara koneksi WebSocket ke backend.
2. Mengelola state `Map<deviceId, MonitorResult>` yang berisi status terkini semua device.
3. Menangani tiga tipe pesan: `initial_state`, `check_result`, dan `status_change`.
4. Melakukan reconnect otomatis setelah koneksi terputus.
5. Memberitahu parent component saat reconnect agar dapat re-fetch data REST.

**Lokasi file:** `frontend/src/hooks/useWebSocket.ts`

### 3.2 State Map<deviceId, MonitorResult>

Sebelum Fase 5, state monitoring hanya menyimpan satu hasil ping. Setelah Fase 5, state berubah menjadi `Map<number, MonitorResult>` yang menyimpan status semua device aktif secara bersamaan.

```ts
// Sebelum: hanya satu result
const [monitorResult, setMonitorResult] = useState<MonitorResult | null>(null);

// Sesudah: Map untuk semua device
const [monitorResults, setMonitorResults] = useState<Map<number, MonitorResult>>(new Map());
```

**Mengapa menggunakan Map?**

| Struktur | Kelebihan | Kekurangan |
|---|---|---|
| Array | Mudah di-map dan di-filter | Pencarian by ID membutuhkan iterasi O(n) |
| Object | Pencarian by ID cepat O(1) | Tidak memiliki urutan default |
| **Map** | **Pencarian by ID cepat O(1), iterable, preserves insertion order** | Syntax sedikit berbeda |

Map dipilih karena operasi paling sering dilakukan adalah "update status device tertentu" yang membutuhkan pencarian by device ID. Dengan Map, operasi ini berjalan dalam waktu konstan.

**Struktur MonitorResult:**

```ts
interface MonitorResult {
  device_id: number;
  ip: string;
  status: 'online' | 'offline' | 'warning' | 'unknown';
  latency_ms: number;
  ttl: number;
  seq: number;
  timestamp: string;
  method: string;        // 'ICMP Ping', 'HTTP Check', 'TCP Port'
  details: Record<string, unknown>;
}
```

### 3.3 Penanganan Pesan WebSocket

WebSocket menerima tiga tipe pesan dari backend:

**1. initial_state — Status awal saat client connect**

```ts
if (message.type === 'initial_state' && Array.isArray(message.data)) {
  const results = message.data.map((item: InitialStateItem) => ({
    device_id: item.device_id,
    ip: item.ip,
    status: item.status,
    latency_ms: item.latency_ms,
    ttl: 0,
    seq: 0,
    timestamp: item.last_check,
    method: item.method,   // ← dibaca dari server, bukan hardcoded
    details: {},
  }));
  setMonitorResults(new Map(results.map((item) => [item.device_id, item])));
}
```

**Mengapa initial_state diperlukan?**

Tanpa initial state, client harus menunggu hasil ping pertama (minimal 3 detik) sebelum bisa menampilkan status. Dengan initial state, halaman langsung menampilkan data terkini saat pertama kali dibuka.

**2. check_result — Hasil pengecekan setiap interval**

```ts
if (message.type === 'check_result') {
  const result = message.data as MonitorResult;
  setMonitorResults((current) => new Map(current).set(result.device_id, result));
}
```

Pesan ini dikirim setiap kali engine selesai melakukan ping ke satu device. State Map diperbarui dengan status terbaru device tersebut.

**3. status_change — Perubahan status online↔offline**

```ts
if (message.type === 'status_change') {
  setLastStatusChange(message.data as StatusChange);
}
```

Pesan ini tidak langsung mengubah Map (karena `check_result` yang mengikutinya sudah melakukannya), tetapi disimpan sebagai event yang memicu halaman Alert Center untuk refresh data alert.

### 3.4 Mekanisme Reconnect dan Re-fetch

Koneksi WebSocket dapat terputus karena berbagai alasan: jaringan tidak stabil, server restart, atau timeout. Hook `useWebSocket` menangani hal ini dengan mekanisme reconnect otomatis.

```ts
socket.onclose = () => {
  setIsConnected(false);
  reconnectTimer.current = setTimeout(connect, 3000);  // reconnect setelah 3 detik
};
```

**Masalah: Setelah reconnect, data REST bisa stale**

Ketika WebSocket reconnect, ia menerima `initial_state` terbaru dari server. Namun, data REST (seperti daftar device di halaman Device Management atau data dashboard) tidak otomatis diperbarui. Jika user menambahkan device saat WebSocket sedang disconnect, device baru tidak akan muncul di halaman Dashboard sampai user melakukan refresh manual.

**Solusi: Reconnect callback**

```ts
// App.tsx
const [reconnectKey, setReconnectKey] = useState(0);
const handleReconnect = useCallback(() => setReconnectKey((k) => k + 1), []);
const { isConnected, monitorResults, lastStatusChange } = useWebSocket(handleReconnect);

// DashboardPage.tsx
useEffect(() => { void load(); }, [reconnectKey]);
// reconnectKey berubah → useEffect trigger → data REST di-fetch ulang
```

Dengan mekanisme ini, setiap kali WebSocket berhasil reconnect, `reconnectKey` berubah, dan semua halaman yang menerima prop ini akan mengambil data terbaru dari REST API.

---

## 4. Integrasi Halaman Dashboard

### 4.1 Pengambilan Data dari REST API

Dashboard mengambil dua sumber data saat pertama kali dimuat:

| Endpoint | Data yang Diambil | Fungsi |
|---|---|---|
| `GET /api/dashboard` | Summary (total, online, offline, warning) + 5 alert terbaru | Menampilkan ringkasan sistem |
| `GET /api/monitoring` | Status terakhir semua device | Menampilkan breakdown per tipe |

```ts
const load = async () => {
  const [nextDashboard, nextMonitoring] = await Promise.all([
    fetchDashboard(),
    fetchMonitoring()
  ]);
  setDashboard(nextDashboard);
  setMonitoring(nextMonitoring);
};
```

**Mengapa menggunakan Promise.all?**

Kedua request tidak memiliki ketergantungan satu sama lain. Dengan `Promise.all`, kedua request dijalankan secara paralel dan menunggu keduanya selesai. Ini lebih cepat daripada menjalankan secara berurutan (sequential).

### 4.2 Penggabungan Data API dengan WebSocket Live

Data dari REST API bersifat "snapshot" (status pada saat query), sedangkan WebSocket memberikan update real-time. Dashboard menggabungkan keduanya:

```ts
const liveMonitoring = useMemo(() => monitoring.map((record) => {
  const live = monitorResults.get(record.device_id);
  return live
    ? { ...record, status: live.status, latency_ms: live.latency_ms, last_check: live.timestamp }
    : record;
}), [monitoring, monitorResults]);
```

**Cara kerja:**

1. Ambil semua record dari REST API (`monitoring`).
2. Untuk setiap record, cek apakah ada data lebih baru di WebSocket (`monitorResults`).
3. Jika ada, gunakan data WebSocket (status, latency, timestamp terbaru).
4. Jika tidak ada, pertahankan data dari REST API.

**Mengapa menggunakan useMemo?**

Penggabungan data dilakukan di dalam `useMemo` agar hanya dihitung ulang ketika `monitoring` atau `monitorResults` berubah. Tanpa `useMemo`, penggabungan akan dijalankan pada setiap render, yang dapat memperlambat performa.

### 4.3 Komponen Dashboard

| Komponen | Fungsi | Data |
|---|---|---|
| `MetricsGrid` | Menampilkan 4 kartu metrik: Total, Online, Offline, Warning | `summary` |
| `DeviceSummary` | Breakdown device berdasarkan tipe (Server, Router, Switch, AP, Website) | `deviceBreakdown` |
| `LatestAlerts` | 5 alert terbaru dengan severity badge | `latestAlerts` |
| `SystemStatus` | Status monitoring, interval, last scan, notifikasi | `systemStatus` |
| `QuickActions` | Tombol aksi: Add Device, View Monitoring, View Alerts, Refresh | Event handlers |

---

## 5. Integrasi Halaman Device Management

### 5.1 CRUD Real melalui API

Device Management merupakan halaman yang paling banyak berubah dari Fase 4 ke Fase 5. Seluruh operasi CRUD yang sebelumnya hanya memanipulasi local state React sekarang benar-benar berkomunikasi dengan backend.

| Operasi | Fungsi API | Endpoint | Efek |
|---|---|---|---|
| Load devices | `fetchDevices()` | `GET /api/devices` | Mengambil daftar dari SQLite |
| Create device | `createDevice(data)` | `POST /api/devices` | Menyimpan ke SQLite, engine start jika active |
| Update device | `updateDevice(id, data)` | `PUT /api/devices/{id}` | Memperbarui di SQLite, restart monitoring |
| Delete device | `deleteDevice(id)` | `DELETE /api/devices/{id}` | Menghapus dari SQLite, stop monitoring |

**Flow Create Device:**

```text
User mengisi form dan klik "Add Device"
        │
        ▼
DeviceFormModal.onSave(deviceInput)
        │
        ▼
DeviceManagementPage.save(input)
        │
        ▼
createDevice(input)  ← api.ts
        │
        ▼
POST /api/devices
Content-Type: application/json
Body: { name, type, ip, method, check_interval, ... }
        │
        ▼
Backend: validasi → INSERT ke SQLite → engine.Start() jika active
        │
        ▼
Response: { success: true, data: { id, name, ip, ... } }
        │
        ▼
Frontend: update state devices → tutup modal
```

### 5.2 Form Device

Form device pada Fase 5 disesuaikan untuk metode ICMP Ping saja, sesuai keputusan Fase 5:

| Field | Tipe | Keterangan |
|---|---|---|
| Device Name | text | Wajib diisi |
| Device Type | select | Server, Router, Switch, Access Point, Website |
| Status | select | Active, Inactive |
| IP Address | text | Wajib diisi, format monospace |
| Check Interval | number | Minimal 1 detik, default 3 |
| Location | text | Opsional |
| Description | textarea | Opsional |

**Perubahan dari Fase 4:**

- Method hardcoded `ICMP Ping` (tidak ada dropdown method).
- Tidak ada field Port (karena ICMP tidak butuh port).
- Tidak ada field URL (karena ICMP tidak butuh URL).
- Location tidak wajib diisi (sebelumnya ada validasi location wajib).

### 5.3 Loading dan Error State

Setiap operasi CRUD memiliki state loading dan error:

```ts
const [isSaving, setSaving] = useState(false);
const [error, setError] = useState('');

const save = async (input: DeviceInput) => {
  setSaving(true);
  try {
    const saved = editing
      ? await updateDevice(editing.id, input)
      : await createDevice(input);
    // update state...
    setFormOpen(false);
  } catch (err) {
    setError(err instanceof Error ? err.message : 'Gagal menyimpan device.');
  } finally {
    setSaving(false);
  }
};
```

**Mengapa error ditangkap di sini?**

Error dari API (misal: network failure, validasi backend) perlu ditampilkan kepada pengguna sebagai pesan yang dapat dimengerti, bukan error teknis yang membingungkan. Error ditangkap dan dikonversi menjadi pesan lokal yang ditampilkan di dalam form.

---

## 6. Integrasi Halaman Alert Center

### 6.1 Pengambilan Alert dengan Filter Backend

Alert Center mengambil data dari backend dengan filter yang dikirim sebagai query parameter:

```ts
const load = useCallback(async () => {
  const data = await fetchAlerts({
    ...(status !== 'all' && { status }),
    ...(severity !== 'all' && { severity }),
    ...(deviceType !== 'all' && { device_type: deviceType }),
  });
  setAlerts(data.map(presentAlert));
}, [status, severity, deviceType]);
```

**Mengapa filter di backend, bukan frontend?**

| Pendekatan | Kelebihan | Kekurangan |
|---|---|---|
| Filter di backend | Mengurangi data yang ditransfer, performa baik jika data banyak | Lebih banyak query parameter |
| Filter di frontend | Semua data diambil, filter lokal cepat | Transfer data besar, lambat jika data banyak |

Gamon memilih filter di backend karena jumlah alert dapat bertambah signifikan seiring waktu. Mengambil semua alert lalu memfilter di frontend akan membebani jaringan dan memori browser.

### 6.2 Pencarian Teks di Client

Meskipun filter status, severity, dan device type dilakukan di backend, pencarian teks (search by name/title) dilakukan di client:

```ts
const filtered = useMemo(() => alerts.filter((alert) =>
  !search || `${alert.title} ${alert.device} ${alert.description}`
    .toLowerCase().includes(search.toLowerCase())
), [alerts, search]);
```

**Mengapa pencarian di client?**

Pencarian teks bersifat eksploratif dan sering berubah-ubah. Setiap kali user mengetik satu huruf, tidak efisien jika harus mengirim request ke backend. Dengan pencarian client, respons langsung tanpa latensi jaringan.

### 6.3 Resolve Alert dan Refresh

Ketika user mengklik "Resolve" pada alert:

```ts
const markResolved = async (id: number) => {
  await resolveAlert(id);        // PUT /api/alerts/{id}/resolve
  await load();                  // Refresh daftar alert dari backend
  setSelected((current) =>       // Update panel detail jika alert yang dipilih
    current?.id === id
      ? { ...current, status: 'resolved', resolvedTime: new Date().toLocaleString('id-ID') }
      : current
  );
};
```

**Mengapa perlu load() ulang?**

Setelah resolve, alert berubah dari `ongoing` ke `refresh dari backend. Refresh diperlukan karena:
1. Alert yang di-resolve mungkin tidak lagi muncul di filter "ongoing".
2. Jumlah alert summary cards berubah.
3. Data harus sinkron dengan database.

---

## 7. Halaman Monitoring (Baru)

### 7.1 Ringkasan Status Real-Time

Halaman Monitoring merupakan halaman baru yang dibuat pada Fase 5. Halaman ini menjawab pertanyaan: **"Perangkat mana yang sedang mengalami kondisi tersebut?"**

Monitoring Summary menampilkan:

| Metrik | Sumber Data | Keterangan |
|---|---|---|
| Total Device | `monitoring.length` | Jumlah semua device yang terdaftar |
| Online | Filter status === 'online' | Device yang merespons dengan baik |
| Offline | Filter status === 'offline' | Device yang tidak merespons |
| Warning | Filter status === 'warning' | Device dengan latency >= 200ms |

### 7.2 Filter dan Pencarian

Halaman Monitoring menyediakan tiga jenis filter:

| Filter | Tipe | Opsi |
|---|---|---|
| Search | Text input | Cari berdasarkan nama device atau IP |
| Status | Select | All, Online, Offline, Warning, Unknown |
| Device Type | Select | All, Server, Router, Switch, Access Point, Website |

Semua filter dilakukan di client setelah data diambil dari REST API dan digabung dengan data WebSocket.

### 7.3 Panel Detail Device

Ketika user mengklik satu device di tabel, panel detail muncul di sisi kanan:

| Informasi | Sumber |
|---|---|
| Name, Type, IP | Dari REST API (`fetchDevices`) |
| Status, Latency, Last Check | Dari WebSocket (`monitorResults`) |
| Method, Interval | Dari REST API |
| Description | Dari REST API |
| History | Dari REST API (`fetchDeviceHistory`, maks 50 record) |

### 7.4 Riwayat Pengecekan

Riwayat pengecekan diambil dari endpoint `GET /api/monitoring/{id}/history` yang mengembalikan maksimum 50 record terakhir. Data ditampilkan dalam format teks:

```
online · 12.5 ms · 2026-07-30T14:30:00Z
online · 15.2 ms · 2026-07-30T14:30:03Z
warning · 250.8 ms · 2026-07-30T14:30:06Z
```

---

## 8. Kontrak Data Frontend-Backend

### 8.1 Generalisasi Tipe TypeScript

Fase 5 menggeneralisasi tipe TypeScript di `frontend/src/types/index.ts` agar mendukung data dari WebSocket dan REST API:

| Tipe | Fungsi |
|---|---|
| `MonitorResult` | Data real-time dari WebSocket (status, latency, method, details) |
| `StatusChange` | Event perubahan status (device_id, old_status, new_status, timestamp) |
| `Alert` | Data alert yang sudah di-mapper untuk tampilan UI |
| `DashboardData` | Struktur data dashboard (summary, breakdown, alerts, system status) |
| `DeviceBreakdown` | Jumlah device per tipe (type, count, online) |

### 8.2 Mapper snake_case → camelCase (presenters.ts)

Backend Go menghasilkan JSON dengan format `snake_case` (misal: `device_name`, `started_at`), sedingga frontend React menggunakan `camelCase` (misal: `deviceName`, `startTime`). `presenters.ts` bertanggung jawab untuk konversi ini.

```ts
export function presentAlert(alert: APIAlert | DashboardAlert): Alert {
  return {
    id: alert.id,
    title: alert.title,
    device: alert.device_name,              // snake_case → camelCase
    deviceType: alert.device_type,          // snake_case → camelCase
    status: alert.status,
    severity: alert.severity,
    startTime: formatDateTime(alert.started_at),  // formatted
    resolvedTime: formatDateTime(alert.resolved_at),
    description: alert.description,
    monitoringMethod: alert.method,
    timestamp: new Date(alert.started_at).valueOf(),
  };
}
```

**Mengapa tidak langsung menggunakan snake_case di frontend?**

Konvensi penamaan di TypeScript/React adalah camelCase. Menggunakan camelCase membuat kode lebih konsisten dengan ekosistem React dan mengurangi kebingungan saat pengembangan.

### 8.3 Format Waktu Indonesia

`presenters.ts` menyediakan fungsi `formatDateTime` untuk mengubah timestamp ISO 8601 menjadi format yang mudah dibaca pengguna Indonesia:

```ts
export function formatDateTime(value: string | null): string | null {
  if (!value) return null;
  const date = new Date(value);
  return date.toLocaleString('id-ID', {
    dateStyle: 'medium',
    timeStyle: 'short'
  });
}
// Output: "30 Jul 2026, 14.30"
```

---

## 9. Pengujian dan Validasi

### 9.1 Go Test (Handler Tests)

Fase 5 menambahkan file `handler/device_test.go` berisi empat pengujian:

| Test | Fungsi | Yang Diuji |
|---|---|---|
| `TestStartStop_MustBePOST` | Memastikan sub-route /start dan /stop hanya menerima POST | HTTP 405 untuk GET/DELETE |
| `TestCreateActiveDevice_StartsMonitoring` | Memastikan create device active langsung memulai monitoring | Engine.IsMonitoring() == true |
| `TestUpdateInactive_StopsMonitoring` | Memastikan update device ke inactive menghentikan monitoring | Engine.IsMonitoring() == false |
| `TestDelete_StopsMonitoring` | Memastikan delete device menghentikan monitoring | Engine.IsMonitoring() == false |

**Cara menjalankan:**

```bash
go test -v ./handler/
```

**Hasil:**

```
=== RUN   TestStartStop_MustBePOST
--- PASS: TestStartStop_MustBePASS (0.02s)
=== RUN   TestCreateActiveDevice_StartsMonitoring
--- PASS: TestCreateActiveDevice_StartsMonitoring (0.01s)
=== RUN   TestUpdateInactive_StopsMonitoring
--- PASS: TestUpdateInactive_StopsMonitoring (0.01s)
=== RUN   TestDelete_StopsMonitoring
--- PASS: TestDelete_StopsMonitoring (0.01s)
PASS
ok  	gamon/handler	0.138s
```

### 9.2 Frontend Build dan Lint

```bash
cd frontend
npm run build    # TypeScript check + Vite bundling
npm run lint     # oxlint quality check
```

Kedua perintah harus menghasilkan output tanpa error. Build memastikan tidak ada TypeScript error, sedangkan lint memastikan kualitas kode sesuai aturan yang dikonfigurasi.

### 9.3 Smoke Test Manual

| Langkah | Aksi | Hasil yang Diharapkan |
|---|---|---|
| 1 | Buka Dashboard | Tampilkan summary device dan latest alerts dari database |
| 2 | Buka Device Management | Tampilkan daftar device dari database |
| 3 | Tambah device active baru | Device muncul di daftar, engine mulai monitoring |
| 4 | Buka Monitoring | Lihat status live device baru (online/warning/offline) |
| 5 | Edit device (ubah IP atau interval) | Monitoring restart dengan konfigurasi baru |
| 6 | Set device inactive | Monitoring berhenti, status tidak lagi update |
| 7 | Resume device (set active) | Monitoring mulai lagi |
| 8 | Hapus device | Device hilang dari daftar, monitoring berhenti |
| 9 | Buat kondisi offline (matikan device) | Setelah 3 failure, alert critical muncul |
| 10 | Nyalakan device kembali | Alert resolved, status kembali online |

---

## 10. Kesimpulan

### 10.1 Prinsip Desain Fase 5

1. **Single Source of Truth** — Database SQLite menjadi satu-satunya sumber data. Frontend hanya menampilkan apa yang ada di database, bukan data dummy.

2. **Separation of Concerns** — REST API menangani operasi transaksional (CRUD), WebSocket menangani data real-time. Keduanya tidak tumpang tindih.

3. **Optimistic UI Update** — WebSocket memberikan update real-time sehingga pengguna tidak perlu refresh halaman untuk melihat perubahan status.

4. **Graceful Degradation** — Jika WebSocket terputus, halaman tetap menampilkan data terakhir yang diketahui. Reconnect otomatis dilakukan setelah 3 detik.

5. **Consistent Error Handling** — Seluruh error dari API ditangkap dan dikonversi menjadi pesan yang dapat dimengerti pengguna.

### 10.2 Capaian Fase 5

| Komponen | Status | Fungsi |
|---|---|---|
| WebSocket Multi-Device | ✅ Selesai | State Map<deviceId, MonitorResult>, reconnect, re-fetch |
| Dashboard Real Data | ✅ Selesai | Summary + breakdown + alerts dari REST API + WebSocket |
| Device Management CRUD | ✅ Selesai | Create, read, update, delete via backend API |
| Device Form | ✅ Selesai | ICMP Ping, interval, location optional |
| Alert Center | ✅ Selesai | Server-side filter, client search, resolve |
| Monitoring Page | ✅ Selesai | Status live, filter, detail panel, history |
| Backend WebSocket | ✅ Selesai | Method field di DeviceStatus |
| Handler Tests | ✅ Selesai | 4 test cases lulus |

### 10.3 Hubungan dengan Fase 6

Fase 6 (Wiring) akan melakukan:

1. Inisialisasi database di `main.go`.
2. Registrasi semua route handler.
3. Startup logic (auto-start monitoring).
4. Koneksi seluruh komponen menjadi aplikasi yang berjalan.

Dengan selesainya Fase 5, seluruh komponen sudah berfungsi secara individual. Fase 6 hanya perlu "menyambung kabel" agar semuanya berjalan sebagai satu kesatuan.

### 10.4 Poin Jawaban untuk Sidang

| Pertanyaan penguji | Jawaban ringkas |
|---|---|
| Apa yang dilakukan Fase 5? | Mengintegrasikan seluruh halaman React dengan REST API dan WebSocket backend, menggantikan data dummy menjadi data real dari SQLite. |
| Mengapa WebSocket diperlukan? | Agar status monitoring diperbarui secara real-time tanpa refresh halaman. REST tidak cocok untuk data yang berubah setiap beberapa detik. |
| Bagaimana jika WebSocket terputus? | Hook useWebSocket melakukan reconnect otomatis setelah 3 detik. Setelah reconnect, data REST di-fetch ulang agar tetap sinkron. |
| Mengapa state monitoring menggunakan Map? | Karena operasi paling sering adalah update status per device ID. Map memberikan pencarian O(1) berdasarkan key. |
| Apa perbedaan REST API dan WebSocket di Gamon? | REST untuk CRUD dan data awal; WebSocket untuk pembaruan status real-time dari engine. |
| Bagaimana data alert dihasilkan? | Engine secara otomatis membuat alert setelah 3 ping gagal berturut-turut (device offline) dan meng-resolve saat device kembali online. |
| Apakah semua halaman sudah menggunakan data real? | Ya. Dashboard, Device Management, Monitoring, dan Alert Center semuanya mengambil data dari backend API. |
| Bagaimana pengujian dilakukan? | Go test untuk handler (4 test cases), npm run build untuk TypeScript check, npm run lint untuk kode quality, dan smoke test manual. |

---

## Daftar Pustaka

[^1]: MDN Web Docs. *WebSocket API*. https://developer.mozilla.org/en-US/docs/Web/API/WebSocket_API (diakses 30 Juli 2026).

[^2]: React Docs. *Hooks Reference*. https://react.dev/reference/react/hooks (diakses 30 Juli 2026).

[^3]: Vite. *Env Variables and Modes*. https://vite.dev/guide/env-and-mode (diakses 30 Juli 2026).

[^4]: Fielding, R., et al. *RFC 9110: HTTP Semantics*. Internet Engineering Task Force, 2022. https://www.rfc-editor.org/rfc/rfc9110.pdf (diakses 30 Juli 2026).

[^5]: TypeScript. *Handbook: Generics*. https://www.typescriptlang.org/docs/handbook/2/generics.html (diakses 30 Juli 2026).

---

*Dokumentasi ini merupakan bagian dari laporan Praktek Kerja Lapangan (PKL) Gamon — Garda Monitoring.*
