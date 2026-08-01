# UX Improvements — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian UX Improvements

UX (User Experience) Improvements merupakan serangkaian peningkatan antarmuka pengguna yang bertujuan untuk meningkatkan kegunaan, efisiensi, dan kepuasan pengguna dalam berinteraksi dengan sistem. Dalam konteks Gamon, fase ini berfokus pada transformasi dari kondisi "berfungsi secara teknis" menjadi "siap untuk demo operasional NOC" (Network Operations Center).

Istilah "improvement" di sini merujuk pada perbaikan yang berdampak langsung pada cara pengguna berinteraksi dengan sistem:

| Aspek | Sebelum Improvement | Sesudah Improvement |
|---|---|---|
| **Toggle Device** | 3 klik (Edit → ubah dropdown → Save) | 1 klik langsung dari tabel |
| **Notifikasi** | Tidak ada feedback visual | Dua layer: toast + banner |
| **Visualisasi Data** | Angka statis tanpa konteks | Grafik latency real-time |

### 1.2 Peran UX Improvements dalam Gamon

UX Improvements memiliki peran krusial dalam menjembatani kesenjangan antara fungsi teknis dan kebutuhan operasional:

| Peran | Penjelasan |
|---|---|
| **Efisiensi Operasional** | Mengurangi jumlah langkah untuk operasi yang sering dilakukan |
| **Deteksi Dini Anomali** | Memberikan notifikasi real-time saat status perangkat berubah |
| **Analisis Performa** | Menyediakan visualisasi data untuk pengambilan keputusan |
| **Pengalaman Pengguna** | Membuat sistem terasa responsif dan informatif |

### 1.3 Kondisi Sebelum Fase 2

Pada akhir Fase 1 (Integrasi MVP), seluruh fitur inti sudah berfungsi: CRUD device, monitoring real-time, dan alert otomatis. Namun, evaluasi menggunakan prototype MVP secara langsung menunjukkan 3 masalah UX kritis:

| Masalah | Dampak | Contoh Skenario |
|---|---|---|
| **Toggle tidak intuitif** | Butuh 3 klik untuk operasi sederhana | Operator ingin menonaktifkan device saat firmware update |
| **Tidak ada notifikasi** | Perubahan status tidak terdeteksi | Device offline tapi operator tidak tahu karena tidak membuka halaman Monitoring |
| **Tidak ada visualisasi latency** | Data tidak bermakna tanpa konteks temporal | Latency 50ms bisa normal atau anomali, tergantung tren |

### 1.4 Tujuan Fase 2

Tujuan utama Fase 2 adalah meningkatkan kegunaan prototype MVP melalui 3 perbaikan UX:

1. **Toggle Aktif/Nonaktif 1 klik** — operator dapat mengaktifkan/menonaktifkan device tanpa membuka form edit
2. **Sistem notifikasi dua lapis** — toast untuk info aplikasi, banner peringatan untuk anomali device
3. **Grafik latency real-time** — visualisasi tren performa jaringan untuk analisis

---

## 2. Fitur 1: Toggle Aktif/Nonaktif

### 2.1 Masalah

Sebelum implementasi fitur ini, untuk mengaktifkan atau menonaktifkan device, operator harus:

```
Klik Edit → Buka form edit lengkap → Ubah dropdown Status → Klik Save
```

Alur ini membutuhkan **3 klik** dan **navigasi form** yang memuat semua field (nama, tipe, IP, dll) — padahal yang diubah hanya status. Dalam skenario NOC di mana operator perlu menonaktifkan banyak device secara cepat (misal: saat maintenance), alur ini sangat tidak efisien.

### 2.2 Solusi

Solusinya adalah menambahkan **tombol toggle** di kolom Actions pada tabel device. Operator cukup **1 klik** untuk mengubah status tanpa membuka form edit.

| Aspek | Sebelum | Sesudah |
|---|---|---|
| Jumlah klik | 3 klik | 1 klik |
| Navigasi | Buka form edit | Tetap di tabel |
| Field yang diakses | Semua field | Hanya status |
| Feedback | Tidak ada (harus reload) | Toast + badge update |

### 2.3 Arsitektur & Alur Data

```
┌─────────────────────────────────────────────────────────────┐
│                    DEVICE TABLE                              │
│                                                              │
│  ┌──────────┬──────────┬──────────┬──────────┬──────────┐   │
│  │ Status   │ Name     │ Type     │ IP       │ Actions  │   │
│  ├──────────┼──────────┼──────────┼──────────┼──────────┤   │
│  │ ● Aktif  │ Router A │ Router   │ 10.0.0.1 │ ✏️ 🗑️ ⏻ │   │
│  │ ○ Nonaktif│ Switch B│ Switch   │ 10.0.0.2 │ ✏️ 🗑️ ⏻ │   │
│  └──────────┴──────────┴──────────┴──────────┴──────────┘   │
│                                        │                     │
│                                        ▼                     │
│                              ┌──────────────┐               │
│                              │  Klik Tombol  │               │
│                              │    Power ⏻    │               │
│                              └──────┬───────┘               │
│                                     │                        │
└─────────────────────────────────────┼────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    API LAYER                                  │
│                                                              │
│  PUT /api/devices/{id}/status                                │
│  Body: {"status": "active"} atau {"status": "inactive"}     │
│                                                              │
│  Response:                                                   │
│  {                                                           │
│    "success": true,                                          │
│    "data": {                                                 │
│      "id": 1,                                                │
│      "status": "inactive",                                   │
│      "message": "Device deactivated"                         │
│    }                                                         │
│  }                                                           │
└─────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    BACKEND HANDLER                            │
│                                                              │
│  1. Update status di database (SQLite)                       │
│  2. Stop/Start engine monitoring                             │
│  3. Return response dengan status baru                       │
│                                                              │
│  Flow:                                                       │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐               │
│  │ Update DB│ →  │ Stop     │ →  │ Return   │               │
│  │ Status   │    │ Engine   │    │ Response │               │
│  └──────────┘    └──────────┘    └──────────┘               │
│                     │                                        │
│                     ▼                                        │
│               ┌──────────┐                                   │
│               │ Start    │ (jika status = "active")          │
│               │ Engine   │                                   │
│               └──────────┘                                   │
└─────────────────────────────────────────────────────────────┘
```

### 2.4 Implementasi Komponen

Perubahan dilakukan pada 4 komponen:

| Komponen | Perubahan |
|---|---|
| `StatusBadge.tsx` | Label "Active" → "Aktif", "Inactive" → "Nonaktif" |
| `DeviceTable.tsx` | Tambah prop `onToggle`, tombol power icon di kolom Actions |
| `DeviceManagementPage.tsx` | Handler `toggleDevice()`, optimistic update state lokal |
| `DeviceFormModal.tsx` | Hapus field Status dari form (create & edit) |

**Tombol Toggle di DeviceTable:**

```tsx
// Desktop table - Actions column
<button
  onClick={() => onToggle(device)}
  className={`p-1.5 rounded-md transition-colors duration-150 cursor-pointer ${
    device.status === 'active'
      ? 'text-success hover:bg-success-muted'
      : 'text-text-muted hover:bg-surface-elevated'
  }`}
  title={device.status === 'active' ? 'Nonaktifkan' : 'Aktifkan'}
>
  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
    <path strokeLinecap="round" strokeLinejoin="round" d="M5.636 5.636a9 9 0 1012.728 0M12 3v9" />
  </svg>
</button>
```

**Handler Toggle di DeviceManagementPage:**

```tsx
const toggleDevice = async (device: Device) => {
  const newStatus = device.status === 'active' ? 'inactive' : 'active';
  try {
    await toggleDeviceStatus(device.id, newStatus);
    setDevices((current) =>
      current.map((d) => d.id === device.id ? { ...d, status: newStatus } : d)
    );
  } catch (err) {
    setError(err instanceof Error ? err.message : 'Gagal mengubah status device.');
  }
};
```

### 2.5 Perubahan Form

Field Status dihapus dari form tambah dan edit. Device baru secara default dalam kondisi "Aktif" (dihandle oleh backend). Hal ini sesuai dengan prinsip desain Gamon: semua device yang didaftarkan harus aktif secara default, dan perubahan status hanya dilakukan secara eksplisit oleh operator.

---

## 3. Fitur 2: Sistem Notifikasi Dua Layer

### 3.1 Masalah

Sistem Gamon sudah mengirim event `status_change` via WebSocket saat device berubah status, tetapi frontend tidak menampilkannya. Operator tidak tahu device tiba-tiba offline kecuali secara aktif membuka halaman Monitoring. Di NOC di mana operator memantau banyak layar, keterlambatan deteksi anomali bisa berakibat fatal.

### 3.2 Arsitektur Notifikasi

Gamon menerapkan sistem notifikasi dua layer sesuai tingkat urgensinya:

```
┌─────────────────────────────────────────────────────────────┐
│                    NOTIFICATION SYSTEM                       │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  Layer 1: TOAST                                         ││
│  │  Target: Info aplikasi                                  ││
│  │  Mekanisme: Auto-dismiss 5 detik, top-right            ││
│  │  Contoh: Error API, CRUD success, WebSocket status      ││
│  │  Komponen: Toast.tsx, ToastContainer.tsx                ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐│
│  │  Layer 2: BANNER "SAKRAL"                               ││
│  │  Target: Anomali device                                 ││
│  │  Mekanisme: Auto-dismiss 10 detik, full-width           ││
│  │  Contoh: Device offline, device warning                 ││
│  │  Komponen: AlertBanner.tsx, AlertBannerContainer.tsx    ││
│  └─────────────────────────────────────────────────────────┘│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

**Mengapa Dua Layer?**

| Pertimbangan | Penjelasan |
|---|---|
| **Toast terlalu "lambat" untuk anomali** | Operator bisa kehilangan notifikasi di pojok layar |
| **Banner full-width menarik perhatian** | Menempati area visual utama, sulit diabaikan |
| **Mencegah notification fatigue** | Tidak semua info sama pentingnya |

### 3.3 Layer 1: Toast Notification

Toast notification digunakan untuk feedback operasional yang tidak memerlukan perhatian segera.

**Kapan muncul:**
- Error saat fetch data API
- WebSocket terputus (disconnect)
- WebSocket reconnect berhasil
- CRUD operation sukses/gagal
- Toggle status sukses/gagal

**Spesifikasi Visual:**

| Aspek | Spesifikasi |
|---|---|
| Posisi | Fixed, top-right layar (z-50) |
| Auto-dismiss | 5 detik |
| Manual close | Tombol X |
| Animasi masuk | Slide-in dari kanan (0.4s) |
| Animasi keluar | Fade-out (0.3s) |
| Warna | Success: hijau, Error: merah, Info: oranye |

**Komponen Toast:**

```tsx
// Toast.tsx - Komponen individual toast
interface ToastItem {
  id: string;
  type: 'success' | 'error' | 'info';
  message: string;
}

// Auto-dismiss via useEffect + setTimeout
useEffect(() => {
  requestAnimationFrame(() => setVisible(true));
  const timer = setTimeout(() => {
    setVisible(false);
    setTimeout(() => onDismiss(toast.id), 300);
  }, 5000);
  return () => clearTimeout(timer);
}, [toast.id, onDismiss]);
```

### 3.4 Layer 2: Banner Peringatan "Sakral"

Banner peringatan dirancang untuk mendapatkan **perhatian segera** dari operator NOC.

**Kapan muncul:**
- WebSocket menerima pesan `status_change` dengan `new_status` = `offline` atau `warning`

**Spesifikasi Visual:**

| Aspek | Spesifikasi |
|---|---|
| Posisi | Fixed, di bawah TopBar (z-40), full-width |
| Warna | Merah (`bg-danger`) untuk offline, Oranye (`bg-warning`) untuk warning |
| Animasi | Pulse ring pada ikon + flash effect pada background |
| Ikon | Warning triangle besar (SVG 24x24) |
| Teks | "PERINGATAN: [Device Name] ([IP]) berubah dari [old_status] ke [new_status]" |
| Aksi | Tombol "Lihat Detail" → navigasi ke MonitoringPage |
| Auto-dismiss | 10 detik, atau manual close (X button) |
| Queue | Maksimum 3 banner; yang lebih lama di-drop |

**Animasi Flash:**

```css
/* index.css */
@keyframes banner-flash {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}

.animate-banner-flash {
  animation: banner-flash 0.4s ease-in-out 3;
}
```

Banner berkedip 3 kali (1.2 detik total) untuk menarik perhatian visual operator.

### 3.5 Integrasi WebSocket

Sistem notifikasi terintegrasi dengan WebSocket yang sudah ada:

```
┌─────────────────────────────────────────────────────────────┐
│                    WEBSOCKET FLOW                             │
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────────┐   │
│  │ Engine   │ →  │ Hub      │ →  │ Client (Browser)     │   │
│  │ (status  │    │ (broadcast)    │                      │   │
│  │  change) │    │          │    │  ┌────────────────┐  │   │
│  └──────────┘    └──────────┘    │  │ useWebSocket   │  │   │
│                                   │  │ (hook)         │  │   │
│                                   │  └───────┬────────┘  │   │
│                                   │          │           │   │
│                                   │          ▼           │   │
│                                   │  ┌────────────────┐  │   │
│                                   │  │ lastStatus     │  │   │
│                                   │  │ Change state   │  │   │
│                                   │  └───────┬────────┘  │   │
│                                   │          │           │   │
│                                   │          ▼           │   │
│                                   │  ┌────────────────┐  │   │
│                                   │  │ AlertBanner    │  │   │
│                                   │  │ Container      │  │   │
│                                   │  └────────────────┘  │   │
│                                   └──────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Pesan `status_change`:**

```json
{
  "type": "status_change",
  "data": {
    "device_id": 1,
    "device_name": "Router A",
    "old_status": "online",
    "new_status": "offline",
    "timestamp": "2026-07-31T12:00:00Z"
  }
}
```

---

## 4. Fitur 3: Latency Chart

### 4.1 Masalah

MonitoringPage sebelumnya hanya menampilkan angka latency saat ini tanpa riwayat. Angka latency tanpa konteks temporal tidak bermakna. Latency 50ms bisa normal untuk satu device tapi anomali untuk device lain. Operator perlu melihat tren naik/turun untuk memutuskan apakah perlu tindakan.

### 4.2 Pilihan Library

Dalam pengembangan Gamon, dipilih **Recharts** sebagai chart library berdasarkan pertimbangan:

| Library | Kelebihan | Kekurangan | Keputusan |
|---|---|---|---|
| **Recharts** | Mature, well-documented, declarative, responsive, MIT | ~60KB | ✅ Dipilih |
| Chart.js | Ringan | Kurang fleksibel untuk custom styling | ❌ Ditolak |
| @nivo | Fitur lengkap | Terlalu berat, banyak dependency | ❌ Ditolak |
| Custom SVG | Full kontrol | Terlalu banyak pekerjaan manual | ❌ Ditolak |
| Lightweight Charts | Ringan | Untuk financial charts, bukan monitoring | ❌ Ditolak |

**Alasan Recharts:**
- Komponen declarative (cocok dengan React)
- Responsive out-of-the-box
- License MIT
- Bundle size ~60KB (acceptable untuk aplikasi NOC)

### 4.3 Spesifikasi Chart

| Aspek | Spesifikasi |
|---|---|
| Tipe | AreaChart (line + gradient fill) |
| Data | 20 data point terakhir dari history |
| X-axis | Waktu (formatted "HH:mm:ss") |
| Y-axis | Latency (ms), auto-scale |
| Line | Stroke `#3b82f6`, strokeWidth 2 |
| Fill | Gradient dari `#3b82f6` (atas) ke transparent (bawah) |
| Grid | Horizontal only, stroke `#44403C` |
| Tooltip | Waktu + latency saat hover |
| Responsive | `ResponsiveContainer` width="100%" height={200} |
| Animasi | `isAnimationActive={false}` untuk real-time |

**Mengapa 20 Data Point?**

20 data point cukup untuk melihat tren 1-2 menit terakhir (dengan interval 3-5 detik). Lebih dari itu membebani rendering tanpa manfaat signifikan.

### 4.4 Integrasi Real-time

Chart update otomatis setiap kali data baru masuk dari WebSocket:

```
┌─────────────────────────────────────────────────────────────┐
│                    DATA FLOW CHART                            │
│                                                              │
│  ┌──────────┐    ┌──────────┐    ┌──────────────────────┐   │
│  │ WebSocket│ →  │ useWeb   │ →  │ MonitoringPage       │   │
│  │ check_   │    │ Socket   │    │                      │   │
│  │ result   │    │ (hook)   │    │  ┌────────────────┐  │   │
│  └──────────┘    └──────────┘    │  │ fetchDevice    │  │   │
│                                   │  │ History()      │  │   │
│                                   │  └───────┬────────┘  │   │
│                                   │          │           │   │
│                                   │          ▼           │   │
│                                   │  ┌────────────────┐  │   │
│                                   │  │ GET /api/      │  │   │
│                                   │  │ monitoring/    │  │   │
│                                   │  │ {id}/history   │  │   │
│                                   │  └───────┬────────┘  │   │
│                                   │          │           │   │
│                                   │          ▼           │   │
│                                   │  ┌────────────────┐  │   │
│                                   │  │ LatencyChart   │  │   │
│                                   │  │ (Recharts)     │  │   │
│                                   │  └────────────────┘  │   │
│                                   └──────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Format Data untuk Chart:**

```typescript
// Transformasi dari PingHistoryRecord ke format chart
const chartData = history
  .slice(0, 20)
  .reverse()
  .map((item) => {
    const date = new Date(item.timestamp);
    const time = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
    return { time, latency: item.latency_ms };
  });
```

---

## 5. Pemilihan Teknologi

### 5.1 Recharts vs Alternatif

Recharts dipilih karena beberapa alasan teknis:

| Kriteria | Recharts | Chart.js | @nivo | Custom SVG |
|---|---|---|---|---|
| **React Integration** | Native (JSX components) | Wrapper library | Native | Manual |
| **Bundle Size** | ~60KB | ~60KB | ~200KB | ~0KB |
| **Learning Curve** | Rendah | Sedang | Tinggi | Tinggi |
| **Customisasi** | Tinggi | Sedang | Tinggi | Full |
| **Responsiveness** | Built-in | Plugin | Built-in | Manual |
| **Dokumentasi** | Sangat baik | Baik | Baik | Tidak ada |
| **License** | MIT | MIT | MIT | - |

### 5.2 CSS Animations vs Library

Untuk animasi notifikasi, Gamon menggunakan CSS animations murni tanpa library tambahan:

| Pendekatan | Kelebihan | Kekurangan |
|---|---|---|
| **CSS Animations** | Zero dependency, performa baik | Kurang fleksibel untuk complex sequences |
| Framer Motion | API lebih mudah, features lengkap | Bundle size besar (~30KB) |
| React Spring | Performa baik | Learning curve tinggi |

CSS animations dipilih karena:
- Tidak menambah bundle size
- Sudah cukup untuk kebutuhan animasi sederhana (fade, slide, pulse)
- Mudah di-maintain
- Bisa di-disable dengan `prefers-reduced-motion`

---

## 6. Struktur File

### 6.1 File yang Dimodifikasi

| File | Perubahan | Fitur |
|---|---|---|
| `components/StatusBadge.tsx` | Label "Aktif"/"Nonaktif" | Toggle |
| `components/DeviceTable.tsx` | Tambah `onToggle` prop + tombol power | Toggle |
| `pages/DeviceManagementPage.tsx` | Tambah `toggleDevice()` handler | Toggle |
| `components/DeviceFormModal.tsx` | Hapus field Status | Toggle |
| `lib/api.ts` | Tambah `toggleDeviceStatus()` | Toggle |
| `components/StatusIndicator.tsx` | Tambah support 'unknown' | Chart |
| `pages/DashboardPage.tsx` | Hapus hardcoded notifications | Cleanup |
| `pages/MonitoringPage.tsx` | Integrasi LatencyChart | Chart |
| `src/index.css` | Animasi `banner-flash` | Notifikasi |

### 6.2 File yang Dibuat Baru

| File | Fungsi | Fitur |
|---|---|---|
| `components/Toast.tsx` | Komponen toast individual | Notifikasi |
| `components/ToastContainer.tsx` | Container untuk multiple toasts | Notifikasi |
| `components/AlertBanner.tsx` | Banner peringatan "sakral" | Notifikasi |
| `components/AlertBannerContainer.tsx` | Queue manager untuk banner | Notifikasi |
| `components/LatencyChart.tsx` | Recharts AreaChart wrapper | Chart |

### 6.3 File yang Tidak Dimodifikasi

- Backend (`handler/`, `monitor/`, `database/`) — tidak ada perubahan
- `main.go` — tidak ada perubahan
- `useWebSocket.ts` — tidak ada perubahan (sudah handle `status_change`)

---

## 7. Flow Data

### 7.1 Flow Toggle Device

```
User klik tombol Power
        │
        ▼
┌───────────────────┐
│ toggleDevice()    │
│ di DeviceMgmtPage │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ toggleDeviceStatus│
│ (API call)        │
│ PUT /api/devices/ │
│ {id}/status       │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ Backend Handler   │
│ 1. Update DB      │
│ 2. Stop/Start     │
│    Engine          │
│ 3. Return Response│
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ Update Local State│
│ (optimistic)      │
│ Badge berubah     │
└───────────────────┘
```

### 7.2 Flow Notifikasi Banner

```
Device status berubah (offline/warning)
        │
        ▼
┌───────────────────┐
│ Engine detects    │
│ status change     │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ Hub broadcasts    │
│ status_change msg │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ useWebSocket hook │
│ updates           │
│ lastStatusChange  │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ AlertBanner       │
│ Container         │
│ shows banner      │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ Auto-dismiss      │
│ after 10 seconds  │
│ or manual close   │
└───────────────────┘
```

### 7.3 Flow Latency Chart

```
User pilih device di MonitoringPage
        │
        ▼
┌───────────────────┐
│ choose(record)    │
│ setSelected()     │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ fetchDeviceHistory│
│ GET /api/monitoring│
│ /{id}/history     │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ Transform data    │
│ to chart format   │
│ [{time, latency}] │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ LatencyChart      │
│ renders AreaChart │
│ with 20 points    │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│ WebSocket update  │
│ → refetch history │
│ → chart updates   │
└───────────────────┘
```

---

## 8. Pengujian

### 8.1 Toggle Active/Inactive

| Skenario | Langkah | Ekspektasi | Status |
|---|---|---|---|
| Toggle aktif → nonaktif | Klik tombol power device aktif | Badge berubah ke "Nonaktif", monitoring berhenti | ✅ |
| Toggle nonaktif → aktif | Klik tombol power device nonaktif | Badge berubah ke "Aktif", monitoring mulai | ✅ |
| Toggle dari form edit | Buka form edit | Field Status tidak ada | ✅ |
| Create device baru | Buka form tambah | Field Status tidak ada, device langsung "Aktif" | ✅ |
| Toggle device offline | Klik tombol power device aktif yang statusnya offline | Status berubah ke "Nonaktif", engine stop | ✅ |

### 8.2 Notifikasi

| Skenario | Langkah | Ekspektasi | Status |
|---|---|---|---|
| WebSocket disconnect | Matikan backend | Banner "PERINGATAN" muncul | ✅ |
| Device offline | Matikan device yang dimonitor | Banner muncul, auto-dismiss 10 detik | ✅ |
| Banner click detail | Klik "Lihat Detail" di banner | Navigasi ke MonitoringPage | ✅ |
| Banner dismiss | Klik X di banner | Banner hilang | ✅ |
| Multiple alerts | Matikan 2 device sekaligus | Banner queue, muncul bergantian | ✅ |

### 8.3 Latency Chart

| Skenario | Langkah | Ekspektasi | Status |
|---|---|---|---|
| Lihat chart | Pilih device di MonitoringPage | Chart muncul, 20 data point terakhir | ✅ |
| Real-time update | Tunggu data baru dari WebSocket | Chart update otomatis | ✅ |
| Tooltip | Hover di chart | Tooltip tampilkan waktu + latency | ✅ |
| Device tanpa history | Pilih device baru | Chart kosong dengan placeholder "Belum ada data" | ✅ |
| Responsive | Resize browser | Chart menyesuaikan lebar container | ✅ |

---

## 9. Bug Fix: API Response Validation

### 9.1 Root Cause Analysis

Selama pengembangan Fase 2, ditemukan bug kritis yang menghambat seluruh fungsionalitas frontend:

**Gejala:**
- WebSocket status: "Golang: Connected"
- Error di frontend: "Request gagal (200)."
- Backend log normal, tidak ada error

**Analisis:**

Error "Request gagal (200)" berasal dari `frontend/src/lib/api.ts:152`:

```typescript
if (!response.ok || !isMessageResponse(payload) || !payload.success) {
  const message = isMessageResponse(payload) ? payload.message : `Request gagal (${response.status}).`;
  throw new APIError(message, response.status);
}
```

`isMessageResponse` mengecek apakah payload punya field `message`:

```typescript
function isMessageResponse(payload: unknown): payload is MessageResponse {
  return typeof payload === 'object' && payload !== null && 'message' in payload;
}
```

**Masalah:** Backend mengirim response dalam format `DataResponse`:

```json
{ "success": true, "data": [...] }
```

`DataResponse` **tidak punya field `message`**, hanya `success` + `data`. Jadi `isMessageResponse(payload)` selalu `false` → kondisi `!isMessageResponse(payload)` jadi `true` → error selalu dilempar meskipun response valid.

### 9.2 Fix Implementation

**Solusi:** Ganti `isMessageResponse` dengan `isSuccessResponse` yang mengecek field `success`:

```typescript
// Sebelum
function isMessageResponse(payload: unknown): payload is MessageResponse {
  return typeof payload === 'object' && payload !== null && 'message' in payload;
}

// Sesudah
function isSuccessResponse(payload: unknown): payload is { success: boolean } {
  return typeof payload === 'object' && payload !== null && 'success' in payload;
}
```

```typescript
// Kondisi error diperbaiki
if (!response.ok || !isSuccessResponse(payload) || !payload.success) {
  const message = 'message' in (payload as object) && typeof (payload as MessageResponse).message === 'string'
    ? (payload as MessageResponse).message
    : `Request gagal (${response.status}).`;
  throw new APIError(message, response.status);
}
```

**Mengapa Fix Ini Benar:**

- `isSuccessResponse` mengecek field `success` yang DIMILIKI oleh **kedua** format response (DataResponse DAN APIResponse/MessageResponse)
- TypeScript type guard membuat `payload.success` aman diakses
- Error message extraction tetap cek `message` field untuk kasus error response

---

## 10. Kesimpulan

### 10.1 Prinsip Desain

Fase 2 UX Improvements menerapkan beberapa prinsip desain penting:

| Prinsip | Penerapan |
|---|---|
| **Efficiency** | Toggle 1 klik menggantikan 3 klik form edit |
| **Visibility** | Notifikasi dua layer memastikan informasi kritis terlihat |
| **Feedback** | Setiap aksi user mendapat feedback visual (toast/banner) |
| **Consistency** | Warna dan animasi konsisten dengan theme Gamon |
| **Accessibility** | Animasi bisa di-disable dengan `prefers-reduced-motion` |

### 10.2 Capaian Fase 2

| Fitur | Status | Capaian |
|---|---|---|
| Toggle Aktif/Nonaktif | ✅ Selesai | 1 klik langsung dari tabel, form cleaned up |
| Notifikasi Dua Layer | ✅ Selesai | Toast + banner terintegrasi dengan WebSocket |
| Latency Chart | ✅ Selesai | Recharts AreaChart, real-time update |
| Bug Fix API | ✅ Selesai | `isSuccessResponse` menggantikan `isMessageResponse` |

### 10.3 Poin Jawaban untuk Sidang (Q&A)

**Q: Mengapa toggle dilakukan dari tabel, bukan dari form edit?**
A: Operator NOC membutuhkan aksi cepat dari satu layar. Toggle dari tabel = 1 klik, sedangkan dari form edit = 3 klik + navigasi. Efisiensi waktu sangat krusial dalam operasional NOC.

**Q: Mengapa notifikasi dibagi menjadi dua layer?**
A: Karena tidak semua informasi sama pentingnya. Toast untuk info operasional (CRUD, WebSocket), banner untuk anomali device (offline/warning). Dua layer mencegah "notification fatigue" — operator tidak kehilangan peringatan penting di antara notifikasi biasa.

**Q: Mengapa banner disebut "sakral"?**
A: Di NOC, operator bisa terlena dengan banyak informasi. Banner harus "berteriak" secara visual tanpa memblokir seluruh layar. Animasi pulse + flash + full-width memastikan anomali tidak terlewat.

**Q: Mengapa Recharts dipilih sebagai chart library?**
A: Recharts sudah mature, well-documented, komponen declarative (cocok dengan React), responsive out-of-the-box, dan license MIT. Bundle size ~60KB acceptable untuk aplikasi NOC.

**Q: Bagaimana cara menangani race condition di WebSocket?**
A: Menggunakan `client.done` channel sebagai lifecycle signal. Hub hanya menutup `client.done` saat client disconnect, tidak menutup `client.send`. Semua goroutine mengecek `<-client.done` sebelum mengirim data.

**Q: Apa yang terjadi jika client disconnect saat sleep 100ms di sendInitialState?**
A: Goroutine akan mengecek `<-client.done` setelah sleep. Jika client sudah disconnect, goroutine langsung return tanpa mengirim data ke channel yang sudah ditutup.

---

*Dokumentasi ini merupakan bagian dari laporan Praktek Kerja Lapangan (PKL) Gamon — Garda Monitoring.*
