# Monitoring — Halaman Pemantauan Real-Time

---

## 1. HALAMAN INI UNTUK APA? (Tujuan / Fungsi)

Monitoring adalah halaman yang berfungsi sebagai **pusat pemantauan real-time** bagi network admin untuk melihat kondisi setiap perangkat jaringan secara detail. Berbeda dengan Dashboard yang menampilkan ringkasan, halaman ini memberikan **kedalaman informasi** lebih lengkap tentang status, latency, dan riwayat pemantauan setiap device.

**Tujuan utamanya:**

- **Memantau status real-time setiap device** — Admin bisa melihat device mana yang Online dan Offline secara langsung
- **Melihat latency (waktu respon) setiap device** — Admin bisa mengetahui seberapa cepat device merespon
- **Melihat riwayat pemantauan** — Admin bisa melihat grafik latency 50 data terakhir untuk menganalisis pola
- **Melihat detail informasi device** — IP, tipe, metode monitoring, interval, lokasi, dan deskripsi
- **Memfilter dan mencari device** — Admin bisa menemukan device tertentu dengan cepat

**Kenapa halaman Monitoring penting?**

Halaman ini adalah **jantung dari aplikasi GAMON**. Di sinilah admin bisa melihat langsung bagaimana kinerja setiap perangkat jaringan, mendeteksi masalah latency, dan menganalisis apakah device berjalan normal atau ada pola gangguan tertentu.

---

## 2. USER BISA MELAKUKAN APA? (Aksi yang Dilakukan Network Admin)

Di halaman Monitoring, network admin bisa melakukan **5 hal utama**:

---

### A. Melihat Daftar Device dengan Status Real-Time

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Tabel berisi semua perangkat beserta status pemantauan terkini |
| **Dimana letaknya** | Area kiri halaman (sekitar 2/3 lebar layar) |
| **Informasi yang ditampilkan** | Nama device, IP, tipe, status (Online/Offline/Unknown), latency (ms), waktu pengecekan terakhir |
| **Kapan berubah** | Secara otomatis dan real-time tanpa perlu refresh |

**Keterangan kolom di tabel:**

| Kolom | Penjelasan |
|-------|------------|
| **Device** | Nama perangkat, alamat IP, dan tipe device (contoh: "Router-A, 192.168.1.1, Router") |
| **Status** | Indikator Online (hijau berkedip), Offline (merah), atau Unknown (abu-abu) |
| **Latency** | Waktu respon device dalam milidetik (ms). Semakin kecil semakin bagus |
| **Last Check** | Waktu terakhir kali device dicek (format: jam:menit:detik) |

---

### B. Mencari Device (Search)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Mencari device berdasarkan nama atau alamat IP |
| **Dimana letaknya** | Search bar di bagian atas halaman |
| **Cara kerja** | Admin mengetik nama atau IP, sistem secara otomatis memfilter daftar device |
| **Sifat pencarian** | Tidak case-sensitive (huruf besar/kecil tidak mempengaruhi hasil) |

---

### C. Memfilter Device Berdasarkan Status

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Memfilter device hanya menampilkan yang memiliki status tertentu |
| **Dimana letaknya** | Dropdown filter di sebelah search bar |
| **Pilihan filter** | **All status** (semua), **Online** (hanya yang hidup), **Offline** (hanya yang mati), **Unknown** (tidak diketahui) |
| **Kapan digunakan** | Saat admin ingin fokus melihat device yang bermasalah (Offline) atau device yang berjalan normal (Online) |

---

### D. Memfilter Device Berdasarkan Tipe

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Memfilter device hanya menampilkan tipe tertentu |
| **Dimana letaknya** | Dropdown filter di sebelah filter status |
| **Pilihan filter** | All device types, Server, Router, Switch, Access Point, Website |
| **Kapan digunakan** | Saat admin ingin melihat hanya Router, atau hanya Server, atau tipe device lainnya |

---

### E. Melihat Detail Device dan Grafik Latency

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Memilih device untuk melihat informasi detail dan grafik latency |
| **Dimana letaknya** | Panel kanan (sekitar 1/3 lebar layar) |
| **Cara mengakses** | Klik pada baris device di tabel |
| **Informasi yang ditampilkan** | Nama device, tipe, metode, IP, status, latency, interval, lokasi, deskripsi |
| **Fitur unggulan** | Grafik latency 50 data terakhir (Area Chart) |

**Detail informasi di panel kanan:**

| Field | Penjelasan |
|-------|------------|
| **Nama Device** | Nama perangkat yang dipilih |
| **Tipe & Metode** | Jenis perangkat dan metode pemantauan (ICMP Ping, HTTP Check, atau TCP Port) |
| **IP Address** | Alamat IP perangkat |
| **Status** | Status terkini: Online, Offline, atau Unknown |
| **Latency** | Waktu respon terkini dalam milidetik |
| **Interval** | Jarak waktu antara satu pengecekan ke pengecekan berikutnya (dalam detik) |
| **Lokasi** | Lokasi fisik perangkat (jika diisi) |
| **Deskripsi** | Keterangan tambahan tentang perangkat (jika diisi) |
| **Grafik Latency** | grafik garis (Area Chart) yang menampilkan 50 data latency terakhir, dari yang paling lama di kiri ke yang paling baru di kanan |

---

## 3. SISTEM MELAKUKAN APA? (Proses di Belakang Aksi)

Berikut penjelasan proses yang terjadi di belakang layar untuk setiap aksi:

---

### A. Saat Admin Membuka Halaman Monitoring

**Proses yang terjadi:**

1. **Browser meminta data device** — Browser mengirim permintaan ke server untuk mendapatkan daftar semua perangkat
2. **Browser meminta data monitoring** — Bersamaan, browser juga meminta data hasil pemantauan terakhir dari setiap device
3. **Server mengambil data dari database** — Server mengambil data device dari tabel devices, kemudian mengambil data pemantauan terakhir (ping_history) dari setiap device
4. **Data dikirim ke browser** — Data device beserta status pemantauan terakhir dikirim ke browser
5. **Data ditampilkan dalam tabel** — Browser menampilkan semua device dalam tabel dengan status, latency, dan waktu pengecekan terakhir

---

### B. Saat Data Monitoring Update Secara Real-Time

**Proses yang terjadi:**

1. **Mesin monitoring berjalan terus-menerus** — Di belakang server, mesin pemantau secara berkala melakukan pengecekan ke setiap device
2. **Hasil pengecekan dikirim ke browser** — Setiap kali ada hasil pengecekan baru, server langsung mengirim data ke browser melalui WebSocket
3. **Data di tabel diperbarui** — Status, latency, dan waktu pengecekan di tabel langsung berubah tanpa perlu refresh
4. **Data di panel detail diperbarui** — Jika device yang sedang dipilih mengalami perubahan, informasi di panel kanan juga langsung berubah

---

### C. Saat Admin Memilih Device (Klik Baris di Tabel)

**Proses yang terjadi:**

1. **Admin menekan baris device** — Browser menandai device yang dipilih sebagai device aktif
2. **Browser meminta riwayat pemantauan** — Browser mengirim permintaan riwayat ping untuk device yang dipilih
3. **Server mengambil data riwayat** — Server mengambil 50 data pemantauan terakhir dari tabel ping_history untuk device tersebut
4. **Data riwayat dikirim ke browser** — 50 data riwayat dikirim ke browser
5. **Panel detail ditampilkan** — Informasi detail device ditampilkan di panel kanan
6. **Grafik latency dibuat** — 50 data riwayat diolah menjadi grafik garis (Area Chart) yang menunjukkan pola latency dari waktu ke waktu

---

### D. Saat Admin Menggunakan Filter Status atau Tipe

**Proses yang terjadi:**

1. **Admin memilih filter** — Admin memilih filter status (Online/Offline) atau filter tipe (Server/Router/Switch/etc.)
2. **Pencocokan dilakukan di browser** — Filter dilakukan di sisi browser (client-side), bukan ke server
3. **Tabel diperbarui** — Hanya device yang cocok dengan filter yang ditampilkan di tabel
4. **Hitungan di header diperbarui** — Jumlah total, online, dan offline di header halaman juga berubah mengikuti filter

---

### E. Saat Admin Mencari Device (Search)

**Proses yang terjadi:**

1. **Admin mengetik di search bar** — Admin mengetik nama device atau IP
2. **Pencocokan dilakukan di browser** — Sistem mencocokkan teks yang diketik dengan nama dan IP semua device
3. **Tabel diperbarui** — Hanya device yang cocok dengan pencarian yang ditampilkan
4. **Pencarian instan** — Hasil pencarian muncul secara langsung tanpa perlu menekan tombol apapun

---

### F. Saat Ada Perubahan Status Device (Online ↔ Offline)

**Proses yang terjadi:**

1. **Mesin monitoring mendeteksi perubahan** — Saat device berubah status (misal: dari Online menjadi Offline)
2. **Perubahan dikirim ke browser** — Melalui WebSocket, perubahan status langsung dikirim ke browser
3. **Indikator status berubah** — Titik warna di tabel berubah (hijau → merah atau sebaliknya)
4. **Latency diperbarui** — Jika device Offline, latency akan menunjukkan nilai terakhir atau "Request Timeout"
5. **Grafik diperbarui** — Jika device yang dipilih berubah status, grafik latency juga akan menambah titik baru

---

## 4. HASILNYA APA? (Output dan Manfaat)

| No | Output | Manfaat untuk Network Admin |
|----|--------|-----------------------------|
| 1 | **Tabel status real-time** | Admin bisa melihat kondisi semua device secara langsung dalam satu tampilan |
| 2 | **Indikator status (Online/Offline/Unknown)** | Admin bisa langsung mengenali device yang bermasalah tanpa perlu membuka detail |
| 3 | **Latency (waktu respon)** | Admin bisa menilai kualitas koneksi setiap device — latency rendah = koneksi bagus |
| 4 | **Waktu pengecekan terakhir** | Admin tahu kapan terakhir kali device dicek, memastikan pemantauan berjalan |
| 5 | **Detail informasi device** | Admin mendapat informasi lengkap tentang device yang dipilih (IP, tipe, lokasi, dll) |
| 6 | **Grafik latency** | Admin bisa menganalisis pola latency — apakah stabil, naik-turun, atau ada lonjakan tiba-tiba |
| 7 | **Filter status** | Admin bisa fokus melihat hanya device yang bermasalah (Offline) untuk penanganan cepat |
| 8 | **Filter tipe** | Admin bisa melihat device berdasarkan kategori (misal: hanya Router atau hanya Server) |
| 9 | **Pencarian cepat** | Admin bisa menemukan device tertentu dari ratusan device dalam hitungan detik |
| 10 | **Update real-time** | Admin tidak perlu refresh halaman, semua data berubah otomatis saat ada perubahan |

---

## 🧠 Rangkuman untuk Sidang

> *"Halaman Monitoring berfungsi sebagai pusat pemantauan real-time bagi network admin. Di halaman ini, admin bisa melihat daftar semua perangkat beserta status pemantauan terkini dalam satu tabel. Setiap device menampilkan informasi status (Online/Offline/Unknown), latency (waktu respon), dan waktu pengecekan terakhir.*
>
> *Admin bisa mencari device berdasarkan nama atau IP, memfilter berdasarkan status (Online/Offline), dan memfilter berdasarkan tipe device (Server, Router, Switch, Access Point, Website). Semua filter dan pencarian dilakukan di browser sehingga hasilnya instan.*
>
> *Ketika admin memilih device dengan mengklik baris di tabel, panel kanan akan menampilkan informasi detail tentang device tersebut meliputi IP, status, latency, interval pengecekan, lokasi, dan deskripsi. Selain itu, ditampilkan juga grafik latency 50 data terakhir yang menunjukkan pola waktu respon device dari waktu ke waktu.*
>
> *Semua data di halaman ini diperbarui secara real-time melalui koneksi WebSocket. Ketika mesin monitoring mendeteksi perubahan status device atau latency baru, data langsung berubah di tabel dan grafik tanpa perlu refresh halaman. Proses di belakang: server mengambil data dari database, menggabungkannya dengan data real-time dari mesin monitoring, lalu mengirimkannya ke browser untuk ditampilkan."*

---

## 📌 Tampilan Monitoring (Layout)

```
┌─────────────────────────────────────────────────────────┐
│  Monitoring                                             │
│  12 devices · 10 online · 2 offline                     │
├─────────────────────────────────────────────────────────┤
│  [🔍 Search device or IP...    ] [Status ▼] [Type ▼]   │
├────────────────────────────────┬────────────────────────┤
│  ┌──────────────────────────┐  │  ┌──────────────────┐  │
│  │ Device    │Stat│Lat│Time │  │  │  Router-A        │  │
│  ├───────────┼────┼────┼────┤  │  │  Router · ICMP   │  │
│  │Router-A   │ 🟢 │1.2│10:30│  │  │                  │  │
│  │192.168.1.1│    │ ms│:01  │  │  │ IP: 192.168.1.1  │  │
│  ├───────────┼────┼────┼────┤  │  │ Status: 🟢 Online│  │
│  │Server-B   │ 🔴 │ — │10:30│  │  │ Latency: 1.2 ms  │  │
│  │10.0.0.5   │    │   │:02  │  │  │ Interval: 30s    │  │
│  ├───────────┼────┼────┼────┤  │  │ Lokasi: Gd. A    │  │
│  │Switch-C   │ 🟢 │0.8│10:30│  │  │                  │  │
│  │192.168.2.1│    │ ms│:03  │  │  │ ─────────────── │  │
│  ├───────────┼────┼────┼────┤  │  │ Latency History  │  │
│  │Web-D      │ 🟢 │45 │10:30│  │  │  ┌──────────┐   │  │
│  │google.com │    │ ms│:04  │  │  │  │ /\  /\    │   │  │
│  └───────────┴────┴────┴────┘  │  │  │/  \/  \___│   │  │
│                                │  │  └──────────┘   │  │
│  Klik baris untuk lihat detail │  │  50 data terakhir│  │
│                                │  └──────────────────┘  │
└────────────────────────────────┴────────────────────────┘
```
