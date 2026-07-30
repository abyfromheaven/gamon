# PRD: UX Improvements — Toggle, Notifikasi & Latency Chart

## Pendahuluan

### Latar Belakang

Setelah 6 fase integrasi MVP Gamon selesai (Database, Backend API, Engine, Frontend API, Frontend Integration, Wiring), sistem sudah berfungsi penuh: CRUD device real, monitoring real-time via WebSocket, dan alert otomatis. Namun dari evaluasi menggunakan prototype MVP secara langsung, ditemukan 3 masalah UX kritis yang menghambat kegunaan sistem dalam skenario operasional NOC (Network Operations Center):

1. **Toggle Active/Inactive tidak intuitif** — Untuk mengaktifkan/menonaktifkan device, user harus klik Edit → ubah dropdown Status → Save. Padahal operasi ini sering dilakukan (misal: saat maintenance, device perlu dinonaktifkan sementara). Butuh 3 klik + navigasi form yang seharusnya 1 klik saja.

2. **Tidak ada notifikasi** — Sistem sudah mengirim event `status_change` via WebSocket, tapi frontend tidak menampilkannya. User tidak tahu device tiba-tiba offline kecuali secara aktif membuka halaman Monitoring. Di NOC, keterlambatan deteksi anomali bisa berakibat fatal.

3. **Tidak ada visualisasi latency** — MonitoringPage hanya menampilkan angka latency saat ini tanpa riwayat. Untuk menganalisis pola (misal: latency naik turun menjelang jam sibuk), operator perlu melihat grafik.

### Tujuan

Meningkatkan kegunaan prototype MVP dari kondisi "berfungsi secara teknis" menjadi "siap untuk demo operasional NOC" melalui 3 perbaikan UX:

- Toggle Aktif/Nonaktif 1 klik langsung dari tabel device
- Sistem notifikasi dua lapis: toast untuk info aplikasi, banner peringatan untuk anomali device
- Grafik latency real-time untuk analisis performa jaringan

---

## Kondisi Saat Ini

### Toggle Active/Inactive

| Aspek | Kondisi |
|---|---|
| Backend | Endpoint `PUT /api/devices/{id}/status` sudah ada dan berfungsi |
| Frontend — Tabel | Tidak ada tombol toggle; status hanya ditampilkan sebagai badge pasif |
| Frontend — Form | Field `Status` (Active/Inactive) ada di form tambah & edit |
| User flow | Edit → ubah dropdown → Save = 3 langkah untuk operasi sederhana |

**Masalah:** Operator NOC perlu menonaktifkan device sesaat (misal: firmware update) tanpa harus membuka form edit lengkap. Form edit memuat semua field (nama, tipe, IP, dll) yang tidak relevan untuk operasi toggle.

### Sistem Notifikasi

| Aspek | Kondisi |
|---|---|
| Backend | WebSocket sudah mengirim pesan `status_change` saat device online↔offline |
| Frontend — WebSocket hook | `lastStatusChange` di-track di `useWebSocket.ts` tapi tidak ditampilkan ke user |
| Frontend — Dashboard | Field `notifications` di SystemStatus selalu `'Paused'` (hardcoded) |
| Frontend — Alert | Alert hanya muncul di halaman Alert Center setelah user navigate ke sana |

**Masalah:** Di NOC, operator memantau banyak layar. Perubahan status harus menarik perhatian secara visual dan temporal — muncul sekarang, hilang setelah dibaca. Tanpa ini, anomali bisa terlewat berjam-jam.

### Latency Chart

| Aspek | Kondisi |
|---|---|
| Backend | `GET /api/monitoring/{id}/history` mengembalikan `PingHistoryRecord[]` (max 50 terakhir) |
| Frontend — MonitoringPage | Detail panel menampilkan teks: `status · latency_ms · timestamp` (5 baris terakhir) |
| Frontend — Dependency | Tidak ada chart library yang ter-install |

**Masalah:** Angka latency tanpa konteks temporal tidak bermakna. Latency 50ms bisa normal untuk satu device tapi anomali untuk device lain. Operator perlu melihat tren naik/turun untuk memutuskan apakah perlu tindakan.

---

## Desain Keputusan (Design Decisions)

### DD-1: Bahasa Indonesia

Label antarmuka menggunakan Bahasa Indonesia sesuai keputusan Phase 5:

| Istilah Lama | Istilah Baru |
|---|---|
| Active | Aktif |
| Inactive | Nonaktif |
| Notifications | Notifikasi |

### DD-2: Toggle — Satu Tombol di Tabel

**Pilihan:** Tombol toggle di kolom Actions pada tabel device.

**Alternatif yang ditolak:**
- Switch di samping badge status → terlalu kecil, sulit diklik di mobile
- Checkbox bulk toggle → over-engineering untuk use case ini
- Toggle di detail panel → tersembunyi, butuh navigasi

**Alasan:** Operator NOC ingin aksi cepat dari satu layar. Tombol di tabel = 1 klik langsung.

### DD-3: Notifikasi Dua Lapis

Sistem notifikasi dibagi menjadi 2 layer sesuai tingkat urgensinya:

| Layer | Target | Mekanisme | Contoh |
|---|---|---|---|
| **Layer 1: Toast** | Info aplikasi | Auto-dismiss 5 detik, top-right | Error API, WebSocket disconnect, sukses simpan |
| **Layer 2: Banner** | Anomali device | Auto-dismiss 10 detik atau manual close, full-width | Device offline, device warning |

**Alasan:**
- Toast terlalu "lambat" untuk anomali — operator bisa kehilangan notifikasi di pojok layar
- Banner full-width menempati area visual utama, sulit diabaikan
- Dua layer mencegah "notification fatigue" — tidak semua info sama pentingnya

### DD-4: Banner "Sakral" untuk Anomali

Banner anomali harus mendapatkan **perhatian segera** dari operator. Spesifikasi visual:

- **Posisi:** Fixed, di bawah TopBar (z-index tinggi), full-width
- **Warna:** Merah (`bg-danger`) untuk offline, Oranye (`bg-warning`) untuk warning
- **Animasi:** Pulse ring pada ikon + flash effect pada background
- **Ikon:** Warning triangle besar (SVG)
- **Teks:** `PERINGATAN: [Device Name] ([IP]) berubah dari [old_status] ke [new_status]`
- **Aksi:** Tombol "Lihat Detail" → navigasi ke MonitoringPage
- **Durasi:** Auto-dismiss 10 detik, atau manual close (X button)

**Alasan:** Di NOC, operator bisa terlena dengan banyak informasi. Alert harus "berteriak" secara visual tanpa memblokir seluruh layar.

### DD-5: Chart Library — Recharts

**Pilihan:** `recharts` (React charting library, wrapper D3)

**Alternatif yang ditolak:**
- `chart.js` / `react-chartjs-2` → lebih ringan tapi kurang fleksibel untuk custom styling
- `@nivo/terracotta` → terlalu berat, banyak dependency
- Custom SVG → terlalu banyak pekerjaan manual
- `lightweight-charts` (TradingView) → untuk financial charts, bukan monitoring

**Alasan:**
- Recharts sudah mature, well-documented
- Komponen declarative (cocok dengan React)
- Responsive out-of-the-box
- License MIT
- Bundle size ~60KB (acceptable)

### DD-6: Batas Data Chart

- Chart menampilkan **20 data point terakhir** dari history
- Data diambil dari `GET /api/monitoring/{id}/history` (max 50 dari backend, frontend ambil 20)
- X-axis: waktu (timestamp), Y-axis: latency (ms)
- Line warna: `#3b82f6` (blue-500) — konsisten dengan theme

**Alasan:** 20 data point cukup untuk melihat tren 1-2 menit terakhir (dengan interval 3-5 detik). Lebih dari itu membebani rendering tanpa manfaat signifikan.

---

## Spesifikasi Fitur

### Fitur 1: Toggle Aktif/Nonaktif

#### User Flow

1. User membuka halaman Device Management
2. Melihat tabel device dengan kolom Actions
3. Klik ikon power button di samping ikon Edit dan Delete
4. Jika device **Aktif** → badge berubah menjadi "Nonaktif", monitoring berhenti
5. Jika device **Nonaktif** → badge berubah menjadi "Aktif", monitoring dimulai
6. Toast notification muncul: "Device [name] berhasil diaktifkan/dinonaktifkan"

#### Perubahan Form

- Form tambah device: Status field **dihapus**. Device baru selalu dalam kondisi "Aktif" (sesuai PRD Phase 5)
- Form edit device: Status field **dihapus**. Perubahan status hanya melalui tombol toggle

#### Komponen yang Dimodifikasi

| Komponen | Perubahan |
|---|---|
| `DeviceTable.tsx` | Tambah prop `onToggle`, tombol power di kolom Actions |
| `DeviceManagementPage.tsx` | Handler `toggleDevice()`, update state lokal |
| `DeviceFormModal.tsx` | Hapus field Status dari form (create & edit) |
| `StatusBadge.tsx` | Label "Active" → "Aktif", "Inactive" → "Nonaktif" |

#### API Call

```
PUT /api/devices/{id}/status
Body: {"status": "active"} atau {"status": "inactive"}
Response: {"success": true, "data": {"id": 1, "status": "active", "message": "Device activated"}}
```

Backend sudah menangani:
- Update status di database
- Start/stop engine monitoring
- Return response dengan status baru

---

### Fitur 2: Sistem Notifikasi Dua Layer

#### Layer 1: Toast Notification

**Kapan muncul:**
- Error saat fetch data API
- WebSocket terputus (disconnect)
- WebSocket reconnect berhasil
- CRUD operation sukses/gagal
- Toggle status sukses/gagal

**Desain:**
- Posisi: top-right layar
- Auto-dismiss: 5 detik
- Manual close: tombol X
- Animasi: slide-in dari kanan, fade-out
- Warna: `bg-danger` (error), `bg-success` (sukses), `bg-accent` (info)

**Komponen baru:**
- `Toast.tsx` — individual toast item
- `ToastContainer.tsx` — container yang mengelola multiple toasts

#### Layer 2: Banner Peringatan "Sakral"

**Kapan muncul:**
- WebSocket menerima pesan `status_change` dengan `new_status` = `offline` atau `warning`

**Desain:**
- Posisi: fixed, di bawah TopBar (full-width)
- Tinggi: auto (bergantung konten)
- Animasi:
  - Pulse ring pada ikon warning
  - Background flash effect (opacity 100% → 80% → 100%, repeat 3x)
- Ikon: warning triangle (SVG 24x24)
- Teks: nama device, IP, status lama → status baru, timestamp
- Tombol: "Lihat Detail" (navigasi ke MonitoringPage) + X (close)
- Auto-dismiss: 10 detik
- Priority: Jika ada beberapa status change, tampilkan yang terakhir. Yang sebelumnya masuk queue dan muncul setelah yang sebelumnya di-dismiss

**Komponen baru:**
- `AlertBanner.tsx` — banner peringatan individual
- `AlertBannerContainer.tsx` — container yang mengelola queue banner

**Integrasi dengan WebSocket:**
- `useWebSocket.ts` sudah meng-track `lastStatusChange`
- `App.tsx` listen `lastStatusChange` → render `<AlertBannerContainer>`
- `AlertBanner` menerima data: `device_name`, `device_ip`, `old_status`, `new_status`, `timestamp`

---

### Fitur 3: Latency Chart

#### User Flow

1. User membuka halaman Monitoring
2. Klik device di tabel → detail panel muncul di kanan
3. Detail panel menampilkan: nama, tipe, IP, status, method, interval, deskripsi
4. Di bawah deskripsi: grafik latency 20 data point terakhir
5. Grafik update otomatis setiap kali data baru masuk (real-time via WebSocket)

#### Data Flow

```
WebSocket check_result
  → useWebSocket: update monitorResults Map
  → MonitoringPage: detect selected device got new data
  → fetchDeviceHistory(device_id) → GET /api/monitoring/{id}/history
  → Map ke format chart: [{time: timestamp, latency: latency_ms}, ...]
  → Recharts LineChart render
```

#### Komponen Baru

- `LatencyChart.tsx` — komponen chart yang menggunakan Recharts

**Props:**
```typescript
interface LatencyChartProps {
  data: Array<{ time: string; latency: number }>;
}
```

**Spesifikasi Chart:**
- Tipe: LineChart (single line)
- X-axis: waktu (formatted "HH:mm:ss")
- Y-axis: latency (ms), auto-scale
- Line: stroke `#3b82f6`, strokeWidth 2, dot: false
- Grid: horizontal only, stroke `#e5e7eb`
- Tooltip: tampilkan waktu + latency saat hover
- Responsive: `ResponsiveContainer` width="100%" height={200}
- Area fill: gradient dari `#3b82f6` (atas) ke transparent (bawah)

#### Integrasi

- Install `recharts` dependency
- `MonitoringPage.tsx`: render `<LatencyChart>` di detail panel
- Fetch history saat device dipilih + saat ada update dari WebSocket

---

## Dependensi

| Package | Versi | Alasan |
|---|---|---|
| `recharts` | ^2.x | Chart library untuk latency visualization |

Install command:
```bash
cd frontend && npm install recharts
```

---

## File yang Dimodifikasi

| File | Perubahan | Fitur |
|---|---|---|
| `frontend/src/components/DeviceTable.tsx` | Tambah `onToggle` prop + tombol power | Toggle |
| `frontend/src/pages/DeviceManagementPage.tsx` | Tambah `toggleDevice()` handler | Toggle |
| `frontend/src/components/DeviceFormModal.tsx` | Hapus field Status | Toggle |
| `frontend/src/components/StatusBadge.tsx` | Label "Aktif"/"Nonaktif" | Toggle |
| `frontend/src/components/Toast.tsx` | **Baru** — komponen toast | Notifikasi |
| `frontend/src/components/ToastContainer.tsx` | **Baru** — container toast | Notifikasi |
| `frontend/src/App.tsx` | State `toasts[]` + render `<ToastContainer>` | Notifikasi |
| `frontend/src/components/AlertBanner.tsx` | **Baru** — banner peringatan | Notifikasi |
| `frontend/src/components/AlertBannerContainer.tsx` | **Baru** — container banner | Notifikasi |
| `frontend/src/App.tsx` | Render `<AlertBannerContainer>` + listen `lastStatusChange` | Notifikasi |
| `frontend/src/components/LatencyChart.tsx` | **Baru** — Recharts LineChart | Chart |
| `frontend/src/pages/MonitoringPage.tsx` | Render `<LatencyChart>` di detail panel | Chart |
| `frontend/package.json` | Tambah dependency `recharts` | Chart |
| `frontend/src/pages/DashboardPage.tsx` | Hapus hardcoded `notifications: 'Paused'` | Cleanup |

---

## File yang Tidak Dimodifikasi

- Backend (`handler/`, `monitor/`, `database/`) — tidak ada perubahan
- `main.go` — tidak ada perubahan
- `useWebSocket.ts` — tidak ada perubahan (sudah handle `status_change`)
- `api.ts` — tidak ada perubahan (sudah ada `toggleDeviceStatus()`)

---

## Urutan Implementasi

### Fase A: Toggle Aktif/Nonaktif

```
1. StatusBadge.tsx          — ubah label ke Bahasa Indonesia
2. DeviceTable.tsx          — tambah onToggle prop + tombol power
3. DeviceManagementPage.tsx — tambah toggleDevice handler
4. DeviceFormModal.tsx      — hapus field Status
```

### Fase B: Notifikasi

```
5. Toast.tsx                — buat komponen toast
6. ToastContainer.tsx       — buat container toast
7. App.tsx                  — integrasi toast (state + render)
8. AlertBanner.tsx          — buat komponen banner
9. AlertBannerContainer.tsx — buat container banner
10. App.tsx                 — integrasi banner (listen lastStatusChange)
```

### Fase C: Latency Chart

```
11. Install recharts        — npm install recharts
12. LatencyChart.tsx        — buat komponen chart
13. MonitoringPage.tsx      — integrasi chart di detail panel
```

### Fase D: Cleanup

```
14. DashboardPage.tsx       — hapus/hapus field notifications
15. Frontend build test     — npm run build
16. Lint test               — npm run lint
```

---

## Pengujian

### Toggle Active/Inactive

| Skenario | Langkah | Ekspektasi |
|---|---|---|
| Toggle aktif → nonaktif | Klik tombol power device aktif | Badge berubah ke "Nonaktif", monitoring berhenti, toast sukses |
| Toggle nonaktif → aktif | Klik tombol power device nonaktif | Badge berubah ke "Aktif", monitoring mulai, toast sukses |
| Toggle dari form edit | Buka form edit | Field Status tidak ada |
| Create device baru | Buka form tambah | Field Status tidak ada, device langsung "Aktif" |
| Toggle device offline | Klik tombol power device aktif yang statusnya offline | Status berubah ke "Nonaktif", engine stop |

### Notifikasi

| Skenario | Langkah | Ekspektasi |
|---|---|---|
| Error API | Putus backend, buka Dashboard | Toast error muncul top-right, auto-dismiss 5 detik |
| WebSocket disconnect | Putus koneksi WebSocket | Toast "Koneksi terputus" muncul |
| WebSocket reconnect | Koneksi ulang WebSocket | Toast "Koneksi berhasil" muncul |
| Device offline | Matikan device yang dimonitor | Banner "PERINGATAN: [device] berubah dari online ke offline" muncul, auto-dismiss 10 detik |
| Device recover | Hidupkan kembali device | Banner "PERINGATAN: [device] berubah dari offline ke online" muncul |
| Banner click detail | Klik "Lihat Detail" di banner | Navigasi ke MonitoringPage, pilih device yang bersangkutan |
| Banner dismiss | Klik X di banner | Banner hilang |
| Multiple alerts | Matikan 2 device sekaligus | Banner pertama muncul, setelah di-dismiss/disappear, banner kedua muncul |

### Latency Chart

| Skenario | Langkah | Ekspektasi |
|---|---|---|
| Lihat chart | Pilih device di MonitoringPage | Chart muncul di detail panel, 20 data point terakhir |
| Real-time update | Tunggu data baru dari WebSocket | Chart update otomatis, titik baru ditambahkan |
| Tooltip | Hover di chart | Tooltip tampilkan waktu + latency |
| Device tanpa history | Pilih device baru (belum pernah di-check) | Chart kosong dengan placeholder "Belum ada data" |
| Responsive | Resize browser | Chart menyesuaikan lebar container |

---

## Risiko & Mitigasi

| Risiko | Mitigasi |
|---|---|
| Recharts bundle size besar (~60KB) | Acceptable untuk aplikasi NOC (tidak mobile-first) |
| Banner menutupi konten | Banner memiliki z-index tinggi tapi auto-dismiss, tidak blocking |
| Toast hilang terlalu cepat | 5 detik cukup untuk membaca; manual close tersedia |
| Banner queue panjang | Maksimum 3 banner di queue; yang lebih lama di-drop |
| Chart flicker saat update | Gunakan `isAnimationActive={false}` pada update real-time |
| Toggle status race condition | Backend sudah handle: stop engine sebelum update DB |

---

## Fitur yang Tidak Termasuk dalam PRD Ini

- Push notification (browser notification API)
- Email/SMS alerting
- Sound alert
- Notification history / log
- Chart interaktif (zoom, pan, select range)
- Export chart ke PNG/PDF
- Filter chart berdasarkan rentang waktu
- Multi-device comparison chart

---

## Catatan untuk Pengembangan Selanjutnya

### Notifikasi Lanjutan

Setelah PRD ini diimplementasi, notifikasi bisa ditingkatkan dengan:

1. **Browser Notification API** — izinkan notifikasi muncul di luar browser tab
2. **Sound alert** — bunyi notifikasi saat anomali (bisa di-toggle)
3. **Notification history** — log semua notifikasi yang pernah muncul
4. **Filter notifikasi** — hanya tampilkan notifikasi untuk device tertentu

### Chart Lanjutan

1. **Time range selector** — pilih 5 menit, 15 menit, 1 jam terakhir
2. **Multiple metrics** — latency + packet loss dalam satu chart
3. **Threshold line** — garis horizontal untuk batas latency normal
4. **Export** — download chart sebagai PNG
5. **Comparison** — overlay latency 2+ device dalam satu chart
