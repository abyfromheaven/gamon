# TESTING GAMON - Garda Monitoring

Dokumen testing E2E dari perspektif user daily untuk monitoring jaringan.

**Tanggal Testing:** 2 Agustus 2026  
**Environment:** Linux, Backend Go + Frontend React (Vite)

---

## Persiapan

### Langkah 1: Jalankan Backend ✅

```bash
cd /home/aby/gamon
go run main.go
```

Pastikan muncul output:
```
===========================================
   GAMON - Garda Monitoring v0.4
   Web Backend Server
   Running on http://localhost:8080
===========================================
```

### Langkah 2: Jalankan Frontend ✅

```bash
cd /home/aby/gamon/frontend
npm run dev
```

Pastikan muncul output dengan URL `http://localhost:5173`

### Langkah 3: Buka Browser ✅

Buka `http://localhost:5173` di browser.  
Pastikan indikator **WebSocket connected** muncul di bagian atas halaman.

---

## Hasil Testing

| No  | Kategori          | Status | Keterangan |
| --- | ----------------- | ------ | ---------- |
| 1   | Setup & Launch    | ✅      |            |
| 2   | Device Management | ✅      |            |
| 3   | Monitoring Live   | ✅      |            |
| 4   | Alert System      | ✅      |            |
| 5   | Dashboard         | ✅      |            |
| 6   | Edge Cases        | ✅      |            |

---

## 1. SETUP & LAUNCH

### TC-1.1: Backend berjalan 

|                      |                                                           |
| -------------------- | --------------------------------------------------------- |
| **Langkah**          | Jalankan `go run main.go` di terminal                     |
| **Hasil Diharapkan** | Server berjalan di port 8080, tidak ada error di terminal |
| **Status**           | ✅                                                         |
| **Keterangan**       |                                                           |

### TC-1.2: Frontend berjalan 

|                      |                                                              |
| -------------------- | ------------------------------------------------------------ |
| **Langkah**          | Jalankan `npm run dev` di folder frontend                    |
| **Hasil Diharapkan** | Dev server berjalan di port 5173, halaman terbuka di browser |
| **Status**           | ✅                                                            |
| **Keterangan**       |                                                              |

### TC-1.3: WebSocket connected 

|                      |                                                                      |
| -------------------- | -------------------------------------------------------------------- |
| **Langkah**          | Buka `http://localhost:5173`, perhatikan indikator koneksi di TopBar |
| **Hasil Diharapkan** | Indikator menunjukkan "connected" (hijau/aktif)                      |
| **Status**           | ✅                                                                    |
| **Keterangan**       |                                                                      |

### TC-1.4: Halaman Dashboard tampil 

|                      |                                                                         |
| -------------------- | ----------------------------------------------------------------------- |
| **Langkah**          | Setelah login, perhatikan halaman awal                                  |
| **Hasil Diharapkan** | Dashboard tampil dengan metrik (total device, online, offline, warning) |
| **Status**           | ✅                                                                       |
| **Keterangan**       |                                                                         |

---

## 2. DEVICE MANAGEMENT

	### TC-2.1: Tambah device Website (Public)

|                      |                                                                                                                                |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| **Langkah**          | Klik "Add Device" → Isi: Name="Google DNS", Type=Website, IP=8.8.8.8, Interval=3, Description="DNS Google" → Klik "Add Device" |
| **Hasil Diharapkan** | Device berhasil ditambah, muncul di tabel, status active                                                                       |
| **Status**           | ✅                                                                                                                              |
| **Keterangan**       |                                                                                                                                |

### TC-2.2: Tambah device Router (Local)

|                      |                                                                                                                |
| -------------------- | -------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | Klik "Add Device" → Isi: Name="Router Utama", Type=Router, IP=(IP router lokal, misal 192.168.1.1), Interval=3 |
| **Hasil Diharapkan** | Device berhasil ditambah, muncul di tabel                                                                      |
| **Status**           | ✅                                                                                                              |
| **Keterangan**       |                                                                                                                |

### TC-2.3: Tambah device Server (Local)

|                      |                                                                                           |
| -------------------- | ----------------------------------------------------------------------------------------- |
| **Langkah**          | Klik "Add Device" → Isi: Name="NAS Server", Type=Server, IP=(IP server lokal), Interval=5 |
| **Hasil Diharapkan** | Device berhasil ditambah                                                                  |
| **Status**           | Tidak ada NAS                                                                             |
| **Keterangan**       |                                                                                           |

### TC-2.4: Tambah device Switch/AP

|                      |                                                                                          |
| -------------------- | ---------------------------------------------------------------------------------------- |
| **Langkah**          | Klik "Add Device" → Isi: Name="AP Lantai 2", Type="Access Point", IP=(IP AP), Interval=3 |
| **Hasil Diharapkan** | Device berhasil ditambah                                                                 |
| **Status**           | Tidak ada AP                                                                             |
| **Keterangan**       |                                                                                          |

### TC-2.5: Edit device

|                      |                                                                                                              |
| -------------------- | ------------------------------------------------------------------------------------------------------------ |
| **Langkah**          | Klik ikon edit pada device "Google DNS" → Ubah Description menjadi "DNS Server Google" → Klik "Save Changes" |
| **Hasil Diharapkan** | Perubahan tersimpan, deskripsi updated di tabel                                                              |
| **Status**           | ✅                                                                                                            |
| **Keterangan**       |                                                                                                              |

### TC-2.6: Toggle device inactive

|                      |                                                                       |
| -------------------- | --------------------------------------------------------------------- |
| **Langkah**          | Klik toggle button pada device "NAS Server" (dari active ke inactive) |
| **Hasil Diharapkan** | Status berubah ke inactive, device tidak di-monitor                   |
| **Status**           | ✅                                                                     |
| **Keterangan**       |                                                                       |

### TC-2.7: Toggle device active kembali

|                      |                                                                       |
| -------------------- | --------------------------------------------------------------------- |
| **Langkah**          | Klik toggle button pada device "NAS Server" (dari inactive ke active) |
| **Hasil Diharapkan** | Status berubah ke active, monitoring dimulai ulang                    |
| **Status**           | ✅                                                                     |
| **Keterangan**       |                                                                       |

### TC-2.8: Hapus device

|                      |                                                                          |
| -------------------- | ------------------------------------------------------------------------ |
| **Langkah**          | Klik ikon hapus pada salah satu device → Konfirmasi dengan klik "Delete" |
| **Hasil Diharapkan** | Device hilang dari tabel, jumlah device berkurang                        |
| **Status**           | ✅                                                                        |
| **Keterangan**       |                                                                          |

### TC-2.9: Batal hapus device

|                      |                                                           |
| -------------------- | --------------------------------------------------------- |
| **Langkah**          | Klik ikon hapus pada device → Klik "Cancel" atau tombol X |
| **Hasil Diharapkan** | Device tetap ada, tidak terhapus                          |
| **Status**           | ✅                                                         |
| **Keterangan**       |                                                           |

### TC-2.10: Validasi form kosong

|                      |                                                                     |
| -------------------- | ------------------------------------------------------------------- |
| **Langkah**          | Klik "Add Device" → Langsung klik "Add Device" tanpa mengisi apapun |
| **Hasil Diharapkan** | Muncul error "Nama device dan alamat IP wajib diisi"                |
| **Status**           | ✅                                                                   |
| **Keterangan**       |                                                                     |

### TC-2.11: Validasi interval invalid

|                      |                                                                           |
| -------------------- | ------------------------------------------------------------------------- |
| **Langkah**          | Klik "Add Device" → Isi Name dan IP, set Interval = 0 → Klik "Add Device" |
| **Hasil Diharapkan** | Muncul error "Interval harus berupa angka minimal 1 detik"                |
| **Status**           | ✅                                                                         |
| **Keterangan**       |                                                                           |

### TC-2.12: Search device

|                      |                                                        |
| -------------------- | ------------------------------------------------------ |
| **Langkah**          | Ketik "Google" di search bar                           |
| **Hasil Diharapkan** | Hanya device dengan nama "Google" yang muncul di tabel |
| **Status**           | ✅                                                      |
| **Keterangan**       |                                                        |

### TC-2.13: Filter by type

|                      |                                      |
| -------------------- | ------------------------------------ |
| **Langkah**          | Pilih filter type "Router"           |
| **Hasil Diharapkan** | Hanya device type Router yang muncul |
| **Status**           | ✅                                    |
| **Keterangan**       |                                      |

### TC-2.14: Filter All

|                      |                             |
| -------------------- | --------------------------- |
| **Langkah**          | Pilih filter "All"          |
| **Hasil Diharapkan** | Semua device muncul kembali |
| **Status**           | ✅                           |
| **Keterangan**       |                             |

---

## 3. MONITORING LIVE

### TC-3.1: Status real-time update

|                      |                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | Buka halaman Monitoring, biarkan selama 15-30 detik                                                                 |
| **Hasil Diharapkan** | Status (online/offline/warning) dan latency berubah real-time tanpa refresh                                         |
| **Status**           | ✅                                                                                                                   |
| **Keterangan**       | Sudah diperbaiki - chart sekarang auto-update real-time tanpa harus klik ulang |

### TC-3.2: Klik device untuk detail

|                      |                                                                             |
| -------------------- | --------------------------------------------------------------------------- |
| **Langkah**          | Klik salah satu row device di tabel monitoring                              |
| **Hasil Diharapkan** | Panel kanan muncul dengan detail device (nama, IP, status, latency, lokasi) |
| **Status**           | ✅                                                                           |
| **Keterangan**       |                                                                             |

### TC-3.3: Latency chart tampil

|                      |                                                                  |
| -------------------- | ---------------------------------------------------------------- |
| **Langkah**          | Klik device yang sudah di-monitor beberapa saat                  |
| **Hasil Diharapkan** | Grafik latency history (20 data terakhir) tampil di panel detail |
| **Status**           | ✅                                                                |
| **Keterangan**       |                                                                  |

### TC-3.4: Filter status Online

|                      |                                               |
| -------------------- | --------------------------------------------- |
| **Langkah**          | Pilih filter status "Online"                  |
| **Hasil Diharapkan** | Hanya device dengan status online yang muncul |
| **Status**           | ✅                                             |
| **Keterangan**       |                                               |

### TC-3.5: Filter status Offline

|                      |                                                                  |
| -------------------- | ---------------------------------------------------------------- |
| **Langkah**          | Pilih filter status "Offline"                                    |
| **Hasil Diharapkan** | Hanya device offline yang muncul (atau kosong jika semua online) |
| **Status**           | ✅                                                                |
| **Keterangan**       |                                                                  |

### TC-3.6: Filter device type

|                      |                                       |
| -------------------- | ------------------------------------- |
| **Langkah**          | Pilih filter type "Website"           |
| **Hasil Diharapkan** | Hanya device type Website yang muncul |
| **Status**           | ✅                                     |
| **Keterangan**       |                                       |

### TC-3.7: Search device di monitoring

|                      |                                             |
| -------------------- | ------------------------------------------- |
| **Langkah**          | Ketik IP device tertentu di search bar      |
| **Hasil Diharapkan** | Hanya device dengan IP tersebut yang muncul |
| **Status**           | ✅                                           |
| **Keterangan**       |                                             |

### TC-3.8: Subtitle info akurat

|                      |                                                                   |
| -------------------- | ----------------------------------------------------------------- |
| **Langkah**          | Perhatikan subtitle di bawah judul "Monitoring"                   |
| **Hasil Diharapkan** | Jumlah device, online, offline, warning sesuai dengan data aktual |
| **Status**           | ✅                                                                 |
| **Keterangan**       |                                                                   |

### TC-3.9: Pilih device, ganti ke device lain

|                      |                                               |
| -------------------- | --------------------------------------------- |
| **Langkah**          | Klik device pertama, lalu klik device kedua   |
| **Hasil Diharapkan** | Panel detail berubah menampilkan device kedua |
| **Status**           | ✅                                             |
| **Keterangan**       |                                               |

---

## 4. ALERT SYSTEM

### TC-4.1: Alert offline muncul (manual)

|                      |                                                                                                                                                                                                               |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | 1. Pastikan ada device aktif yang bisa di-block (misal: matikan router/AP sementara, atau tambah device dengan IP palsu seperti 10.255.255.1) 2. Tunggu minimal 3 cycle check (9 detik jika interval 3 detik) |
| **Hasil Diharapkan** | Alert "Device Offline" dengan severity "critical" muncul di Alert Center                                                                                                                                      |
| **Status**           | ✅                                                                                                                                                                                                             |
| **Keterangan**       |                                                                                                                                                                                                               |

### TC-4.2: Alert banner muncul saat status change

|                      |                                                                                        |
| -------------------- | -------------------------------------------------------------------------------------- |
| **Langkah**          | Perhatikan bagian atas layar saat device berubah status (misal dari online ke offline) |
| **Hasil Diharapkan** | Banner notifikasi muncul di bagian atas dengan info device dan status change           |
| **Status**           | ✅                                                                                      |
| **Keterangan**       |                                                                                        |

### TC-4.3: Alert auto-resolve saat device recovery

|                      |                                                                                                                              |
| -------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | Setelah alert offline muncul, hidupkan kembali device tersebut (atau hapus device IP palsu, tambah device yang bisa di-ping) |
| **Hasil Diharapkan** | Alert otomatis berubah status ke "resolved"                                                                                  |
| **Status**           | ✅                                                                                                                            |
| **Keterangan**       |                                                                                                                              |

### TC-4.4: Manual resolve alert

|                      |                                                               |
| -------------------- | ------------------------------------------------------------- |
| **Langkah**          | Buka Alert Center → Klik alert ongoing → Klik "Mark Resolved" |
| **Hasil Diharapkan** | Alert status berubah ke resolved, waktu resolved tercatat     |
| **Status**           | ✅                                                             |
| **Keterangan**       |                                                               |

### TC-4.5: Filter alert by status

|                      |                                             |
| -------------------- | ------------------------------------------- |
| **Langkah**          | Pilih filter status "Ongoing"               |
| **Hasil Diharapkan** | Hanya alert yang belum resolved yang muncul |
| **Status**           | ✅                                           |
| **Keterangan**       |                                             |

### TC-4.6: Filter alert by severity

|                      |                                           |
| -------------------- | ----------------------------------------- |
| **Langkah**          | Pilih filter severity "Critical"          |
| **Hasil Diharapkan** | Hanya alert severity critical yang muncul |
| **Status**           | ✅                                         |
| **Keterangan**       |                                           |

### TC-4.7: Filter alert by device type

|                      |                                                 |
| -------------------- | ----------------------------------------------- |
| **Langkah**          | Pilih filter device type (misal "Router")       |
| **Hasil Diharapkan** | Hanya alert dari device type Router yang muncul |
| **Status**           | ✅                                               |
| **Keterangan**       |                                                 |

### TC-4.8: Search alert

|                      |                                                  |
| -------------------- | ------------------------------------------------ |
| **Langkah**          | Ketik nama device atau judul alert di search bar |
| **Hasil Diharapkan** | Alert yang cocok dengan search query muncul      |
| **Status**           | ✅                                                |
| **Keterangan**       |                                                  |

### TC-4.9: Alert summary cards

|                      |                                                                                    |
| -------------------- | ---------------------------------------------------------------------------------- |
| **Langkah**          | Perhatikan kartu ringkasan di atas halaman Alert Center (Total, Ongoing, Resolved) |
| **Hasil Diharapkan** | Jumlah pada kartu sesuai dengan jumlah alert aktual                                |
| **Status**           | ✅                                                                                  |
| **Keterangan**       |                                                                                    |

### TC-4.10: Klik alert untuk detail

|                      |                                                                           |
| -------------------- | ------------------------------------------------------------------------- |
| **Langkah**          | Klik salah satu alert di daftar                                           |
| **Hasil Diharapkan** | Panel detail alert muncul (judul, deskripsi, device, waktu mulai, status) |
| **Status**           | ✅                                                                         |
| **Keterangan**       |                                                                           |

---

## 5. DASHBOARD

### TC-5.1: Metrik devices akurat

|                      |                                                              |
| -------------------- | ------------------------------------------------------------ |
| **Langkah**          | Buka halaman Dashboard, perhatikan kartu metrik              |
| **Hasil Diharapkan** | Total devices, online, offline, warning sesuai jumlah aktual |
| **Status**           | ✅                                                            |
| **Keterangan**       |                                                              |

### TC-5.2: Device breakdown tampil

|                      |                                                                                     |
| -------------------- | ----------------------------------------------------------------------------------- |
| **Langkah**          | Perhatikan bagian "Device Summary" di Dashboard                                     |
| **Hasil Diharapkan** | Breakdown per tipe device (Server, Router, Switch, AP, Website) tampil dengan benar |
| **Status**           | ✅                                                                                   |
| **Keterangan**       |                                                                                     |

### TC-5.3: Latest alerts tampil

|                      |                                                  |
| -------------------- | ------------------------------------------------ |
| **Langkah**          | Perhatikan bagian "Latest Alerts" di Dashboard   |
| **Hasil Diharapkan** | Alert terbaru muncul (maksimal 5 alert terakhir) |
| **Status**           | ✅                                                |
| **Keterangan**       |                                                  |

### TC-5.4: System status info

|                      |                                                                       |
| -------------------- | --------------------------------------------------------------------- |
| **Langkah**          | Perhatikan bagian "System Status"                                     |
| **Hasil Diharapkan** | Monitoring: Running, Check Interval terdeteksi, Notifications: Active |
| **Status**           | ✅                                                                     |
| **Keterangan**       |                                                                       |

### TC-5.5: Quick action - Add Device

|                      |                                           |
| -------------------- | ----------------------------------------- |
| **Langkah**          | Klik tombol "Add Device" di Quick Actions |
| **Hasil Diharapkan** | Navigasi ke halaman Device Management     |
| **Status**           | ✅                                         |
| **Keterangan**       |                                           |

### TC-5.6: Quick action - View Monitoring

|                      |                                                |
| -------------------- | ---------------------------------------------- |
| **Langkah**          | Klik tombol "View Monitoring" di Quick Actions |
| **Hasil Diharapkan** | Navigasi ke halaman Monitoring                 |
| **Status**           | ✅                                              |
| **Keterangan**       |                                                |

### TC-5.7: Quick action - View Alerts

|                      |                                            |
| -------------------- | ------------------------------------------ |
| **Langkah**          | Klik tombol "View Alerts" di Quick Actions |
| **Hasil Diharapkan** | Navigasi ke halaman Alert Center           |
| **Status**           | ✅                                          |
| **Keterangan**       |                                            |

### TC-5.8: Refresh dashboard

|                      |                                                 |
| -------------------- | ----------------------------------------------- |
| **Langkah**          | Klik tombol refresh di Quick Actions            |
| **Hasil Diharapkan** | Data dashboard di-load ulang, metrik ter-update |
| **Status**           | ✅                                               |
| **Keterangan**       |                                                 |

---

## 6. EDGE CASES & STABILITY

### TC-6.1: Backend restart, auto-start monitoring

|                      |                                                                                                               |
| -------------------- | ------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | 1. Pastikan ada device active 2. Matikan backend (Ctrl+C) 3. Jalankan lagi `go run main.go` 4. Perhatikan log |
| **Hasil Diharapkan** | Log menunjukkan "Auto-started monitoring for X active devices", monitoring berjalan otomatis                  |
| **Status**           | ✅                                                                                                             |
| **Keterangan**       |                                                                                                               |

### TC-6.2: WebSocket reconnect otomatis

|                      |                                                                                              |
| -------------------- | -------------------------------------------------------------------------------------------- |
| **Langkah**          | 1. Buka Monitoring page 2. Matikan backend 3. Tunggu beberapa detik 4. Nyalakan backend lagi |
| **Hasil Diharapkan** | Frontend reconnect otomatis dalam ~3 detik, status kembali normal                            |
| **Status**           | ✅                                                                                            |
| **Keterangan**       |                                                                                              |

### TC-6.3: Device dengan IP tidak valid

|                      |                                                                    |
| -------------------- | ------------------------------------------------------------------ |
| **Langkah**          | Tambah device dengan IP="999.999.999.999"                          |
| **Hasil Diharapkan** | Device ditambahkan, tapi status langsung offline karena ping gagal |
| **Status**           | ✅                                                                  |
| **Keterangan**       |                                                                    |

### TC-6.4: Multiple browser tabs

|                      |                                                                     |
| -------------------- | ------------------------------------------------------------------- |
| **Langkah**          | Buka `http://localhost:5173` di 2 tab browser berbeda               |
| **Hasil Diharapkan** | Kedua tab menampilkan data yang sama, update real-time di kedua tab |
| **Status**           | ✅                                                                   |
| **Keterangan**       |                                                                     |

### TC-6.5: Check interval berbeda per device

|                      |                                                                                                                                        |
| -------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| **Langkah**          | 1. Tambah device A dengan interval 3 detik 2. Tambah device B dengan interval 10 detik 3. Buka Monitoring, bandingkan frekuensi update |
| **Hasil Diharapkan** | Device A di-check lebih sering dari device B                                                                                           |
| **Status**           | ✅                                                                                                                                      |
| **Keterangan**       |                                                                                                                                        |

### TC-6.6: Hapus device yang sedang di-monitor

|                      |                                                                      |
| -------------------- | -------------------------------------------------------------------- |
| **Langkah**          | Hapus device yang sedang aktif di-monitoring                         |
| **Hasil Diharapkan** | Device hilang dari daftar, monitoring berhenti untuk device tersebut |
| **Status**           | ✅                                                                    |
| **Keterangan**       |                                                                      |

### TC-6.7: Navigation sidebar

|                      |                                                                     |
| -------------------- | ------------------------------------------------------------------- |
| **Langkah**          | Klik setiap menu di sidebar: Dashboard, Devices, Monitoring, Alerts |
| **Hasil Diharapkan** | Setiap klik membuka halaman yang sesuai                             |
| **Status**           | ✅                                                                   |
| **Keterangan**       |                                                                     |

---

## Rangkuman Hasil

| Kategori | Total | Lulus | Gagal | Keterangan |
|----------|-------|-------|-------|------------|
| Setup & Launch | 4 | | | |
| Device Management | 14 | | | |
| Monitoring Live | 9 | | | |
| Alert System | 10 | | | |
| Dashboard | 8 | | | |
| Edge Cases | 7 | | | |
| **TOTAL** | **52** | | | |
gua males ngitung rangkumannya jir, lu aja yang rangkum.

---

## Catatan Tambahan

### Fix TC-3.1: Latency Chart Auto-Update (2 Agustus 2026)

**Masalah:** Latency chart tidak auto-update saat row device dibuka. User harus klik row device berulang kali untuk refresh chart.

**Root cause:** `history` (data chart) hanya di-fetch sekali saat user klik row (`choose()`). Tidak ada mekanisme untuk append data baru dari WebSocket ke history.

**Solusi:** In-Memory History Cache
- Cache history per device di `useRef<Map<number, PingHistoryRecord[]>>`
- Buka row pertama kali → fetch dari API (database punya semua historical data)
- Saat row terbuka → append new results dari WebSocket ke cache
- Tutup row → stop append, tapi cache tetap di memori
- Buka ulang → langsung pakai cache (tidak reset dari nol)

**File yang diubah:** `frontend/src/pages/MonitoringPage.tsx`
- Tambah `useRef` untuk history cache
- Rewrite `choose()` untuk gunakan cache
- Tambah `useEffect` untuk live update dari WebSocket
- Naikkan max data points dari 20 ke 50

**Hasil:** Chart sekarang real-time update tanpa harus klik ulang, dan history persisten saat row ditutup/dibuka ulang.

