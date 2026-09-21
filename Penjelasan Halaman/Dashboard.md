# Dashboard — Halaman Utama (Overview)

---

## 1. HALAMAN INI UNTUK APA? (Tujuan / Fungsi)

Dashboard adalah **halaman pertama yang ditampilkan** saat network admin membuka aplikasi GAMON. Halaman ini berfungsi sebagai **pusat informasi utama** yang memberikan gambaran menyeluruh tentang kondisi seluruh jaringan dalam satu tampilan.

**Tujuan utamanya:**

- Memberikan **ringkasan cepat** tentang berapa banyak perangkat jaringan yang sedang dipantau
- Memberikan **informasi real-time** tentang device mana yang sedang hidup (Online) dan mati (Offline)
- Menampilkan **log aktivitas monitoring** terakhir secara langsung tanpa perlu refresh
- Menunjukkan **status mesin monitoring** apakah sedang berjalan atau berhenti
- Menampilkan **peringatan (alert) terbaru** agar admin bisa cepat menanggapi masalah

**Kenapa Dashboard penting?**

Tanpa harus membuka halaman satu per satu, admin bisa langsung mengetahui kondisi jaringan secara keseluruhan dalam hitungan detik. Ini menghemat waktu dan mempercepat proses deteksi masalah.

---

## 2. USER BISA MELAKUKAN APA? (Aksi yang Dilakukan Network Admin)

Di halaman Dashboard, network admin bisa melakukan **6 hal utama**:

---

### A. Melihat Jumlah Total Device

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Angka total perangkat jaringan yang sedang dipantau oleh sistem |
| **Dimana letaknya** | Kartu pertama di bagian paling atas halaman Dashboard |
| **Contoh tampilan** | "Total Device — 12 sedang dimonitor" |
| **Kapan berubah** | Saat admin menambah atau menghapus device di halaman Device Management |

---

### B. Melihat Jumlah Device Online

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Angka perangkat yang sedang **hidup dan merespon** pemantauan |
| **Dimana letaknya** | Kartu kedua di bagian paling atas, berwarna hijau |
| **Contoh tampilan** | "Online — 10 perangkat aktif" |
| **Kapan berubah** | Secara real-time, setiap kali mesin monitoring selesai melakukan pengecekan |

---

### C. Melihat Jumlah Device Offline

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Angka perangkat yang **tidak merespon** pemantauan |
| **Dimana letaknya** | Kartu ketiga di bagian paling atas, berwarna merah |
| **Contoh tampilan** | "Offline — 2 tidak merespon" |
| **Kapan berubah** | Secara real-time, saat device tidak merespon pemantauan |

---

### D. Memantau Log Monitoring Secara Real-Time

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Daftar 50 aktivitas monitoring terakhir yang terjadi secara langsung |
| **Dimana letaknya** | Area besar di sisi kiri Dashboard (sekitar 2/3 layar) |
| **Informasi yang ditampilkan** | Waktu, status (Online/Offline), nama device, alamat IP, metode monitoring, nomor urut (sequence), TTL, dan latency (waktu respon) |
| **Fitur khusus** | Admin bisa menghapus log dengan tombol "Clear" jika ingin mengosongkan tampilan |
| **Kapan berubah** | Secara otomatis dan langsung, tanpa perlu refresh halaman |

**Keterangan informasi di setiap baris log:**

| Field | Penjelasan |
|-------|------------|
| **Timestamp** | Waktu terjadinya aktivitas pemantauan (format: jam:menit:detik) |
| **Status** | Hasil pemantauan: **ONLINE** (hijau) jika device merespon, atau **OFFLINE** (merah berkedip) jika device tidak merespon |
| **Device Name** | Nama perangkat yang sedang dipantau |
| **IP Address** | Alamat IP perangkat |
| **Method** | Metode pemantauan yang digunakan: ICMP Ping, HTTP Check, atau TCP Port |
| **Sequence (seq)** | Nomor urut paket pemantauan |
| **TTL (Time to Live)** | Nilai TTL dari paket ping (hanya muncul jika status Online) |
| **Latency** | Waktu respon device dalam milidetik (ms), hanya muncul jika status Online |
| **Request Timeout** | Pesan yang muncul jika device tidak merespon (status Offline) |

---

### E. Melihat Status Mesin Monitoring

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Informasi tentang kondisi mesin pemantau jaringan |
| **Dimana letaknya** | Panel di sisi kanan atas |
| **Informasi yang ditampilkan** | Status mesin (Running/Stopped), interval pengecekan, waktu scan terakhir, dan status notifikasi |

**Detail informasi Status Mesin Monitoring:**

| Field | Penjelasan |
|-------|------------|
| **Monitoring** | Status mesin monitoring: **Running** (sedang berjalan) atau **Stopped** (berhenti). Ditampilkan dengan indikator hijau (berkedip) atau merah |
| **Interval** | Jarak waktu antara satu pengecekan ke pengecekan berikutnya. Diambil dari interval tercepat di antara semua device yang dipantau (contoh: 30 seconds) |
| **Last Scan** | Waktu terakhir kali mesin monitoring melakukan pengecekan ke device |
| **Notif** | Status sistem notifikasi Telegram: **Active** (aktif mengirim peringatan) atau **Paused** (jeda/notifikasi dinonaktifkan) |

---

### F. Melihat 5 Alert (Peringatan) Terbaru

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Daftar 5 peringatan terakhir dari device yang bermasalah |
| **Dimana letaknya** | Panel di sisi kanan bawah |
| **Informasi yang ditampilkan** | Nama device, judul peringatan, dan waktu kejadian |
| **Fitur khusus** | Tombol "View all alerts" untuk melihat semua peringatan di halaman Alert Center |
| **Kapan berubah** | Saat ada device baru yang mengalami masalah |

---

## 3. SISTEM MELAKUKAN APA? (Proses di Belakang Aksi)

Berikut penjelasan proses yang terjadi di belakang layar untuk setiap aksi:

---

### A. Saat Admin Membuka Halaman Dashboard

**Proses yang terjadi:**

1. **Browser meminta data ke server** — Sistem mengirim permintaan ke server untuk mendapatkan data dashboard (jumlah device, status, dan alert)
2. **Server menghitung jumlah device** — Server menghitung total semua perangkat yang terdaftar di database
3. **Server menentukan status setiap device** — Server memeriksa data pemantauan terakhir dari setiap device untuk menentukan apakah device tersebut Online atau Offline. Caranya: server mengambil data ping terakhir dari setiap device, lalu mengecek statusnya
4. **Server mengambil 5 alert terbaru** — Server mengambil data peringatan terakhir yang paling baru dari database, diurutkan dari yang paling baru
5. **Server mengambil data monitoring** — Server mengambil data hasil pemantauan terakhir dari semua device untuk ditampilkan di log
6. **Data dikirim ke browser** — Semua data dikirim ke browser untuk ditampilkan ke admin

---

### B. Saat Log Monitoring Muncul Secara Real-Time

**Proses yang terjadi:**

1. **Mesin monitoring berjalan terus-menerus** — Di belakang server, ada mesin pemantau yang secara berkala melakukan pengecekan ke setiap device
2. **Setiap device dicek secara bergantian** — Mesin mengirim sinyal pemantauan (ICMP Ping, HTTP Check, atau TCP Port) ke setiap device sesuai interval yang ditentukan
3. **Hasil pemantauan dikirim ke browser** — Setiap kali ada hasil pengecekan baru, server langsung mengirim data ke browser menggunakan koneksi WebSocket (koneksi real-time tanpa perlu refresh)
4. **Log ditambahkan ke tampilan** — Browser menerima data baru dan langsung menambahkannya ke daftar log. Log hanya menampilkan 50 entri terbaru, entri paling lama akan terhapus otomatis
5. **Log disimpan di browser** — Agar data tidak hilang saat admin merefresh halaman, sistem menyimpan 50 log terakhir di penyimpanan lokal browser (localStorage)

---

### C. Saat Status Device Berubah (Online ke Offline atau sebaliknya)

**Proses yang terjadi:**

1. **Mesin monitoring mendeteksi perubahan** — Saat device yang sebelumnya Online tiba-tiba tidak merespon, atau sebaliknya device yang sebelumnya Offline kembali merespon
2. **Data diperbarui di server** — Status device diperbarui di database berdasarkan hasil pemantauan terbaru
3. **Perubahan dikirim ke browser** — Melalui WebSocket, perubahan status langsung dikirim ke browser secara real-time
4. **Dashboard otomatis update** — Beberapa bagian Dashboard otomatis berubah:
   - Angka Online/Offline di kartu metrics berubah
   - Log baru muncul dengan status yang sesuai
   - Jika perlu, alert baru akan dibuat secara otomatis oleh sistem

---

### D. Saat Admin Menekan Tombol "Clear" di Log Monitoring

**Proses yang terjadi:**

1. **Log dihapus dari tampilan** — Semua log yang sedang ditampilkan di layar dihapus
2. **Log dihapus dari penyimpanan browser** — Data log yang tersimpan di localStorage juga dihapus agar tidak muncul kembali saat refresh
3. **Log baru tetap bisa masuk** — Meskipun log sudah dihapus, log baru dari pemantauan berikutnya tetap akan muncul kembali

---

### E. Saat Admin Menekan Tombol "View all alerts"

**Proses yang terjadi:**

1. **Browser berpindah halaman** — Sistem berpindah dari halaman Dashboard ke halaman Alert Center
2. **Data alert dimuat** — Browser meminta data seluruh alert ke server
3. **Alert ditampilkan** — Semua peringatan ditampilkan di halaman Alert Center dengan informasi yang lebih lengkap dibandingkan di Dashboard

---

## 4. HASILNYA APA? (Output dan Manfaat)

| No | Output | Manfaat untuk Network Admin |
|----|--------|-----------------------------|
| 1 | **Ringkasan jumlah device** | Admin langsung tahu total aset jaringan yang sedang dipantau tanpa harus menghitung manual |
| 2 | **Jumlah device Online** | Admin yakin device mana yang berfungsi normal dan bisa diandalkan |
| 3 | **Jumlah device Offline** | Admin segera tahu ada berapa device yang bermasalah dan perlu ditindaklanjuti |
| 4 | **Log monitoring real-time** | Admin bisa memantau aktivitas jaringan secara langsung, mengetahui latency setiap device, dan mendeteksi gangguan seketika |
| 5 | **Status mesin monitoring** | Admin memastikan sistem monitoring berjalan normal, atau segera melakukan perbaikan jika mesin berhenti |
| 6 | **Interval pengecekan** | Admin tahu seberapa sering device dicek (misal: setiap 30 detik) |
| 7 | **Waktu scan terakhir** | Admin tahu kapan terakhir kali sistem melakukan pengecekan |
| 8 | **Status notifikasi** | Admin tahu apakah sistem notifikasi Telegram aktif atau tidak |
| 9 | **Alert terbaru** | Admin bisa cepat merespons masalah tanpa harus membuka halaman Alert Center terlebih dahulu |
| 10 | **Koneksi real-time** | Admin tidak perlu refresh halaman secara manual, semua data update sendiri |

---

## 🧠 Rangkuman untuk Sidang

> *"Dashboard adalah halaman utama aplikasi GAMON yang berfungsi sebagai pusat informasi kondisi jaringan. Di halaman ini, admin bisa melihat enam informasi utama: jumlah total device yang sedang dipantau, jumlah device yang Online (hidup), jumlah device yang Offline (tidak merespon), log aktivitas monitoring secara real-time, status mesin monitoring, serta 5 peringatan terbaru.*
>
> *Proses di belakang layar bekerja sebagai berikut: saat admin membuka Dashboard, browser meminta data ke server. Server menghitung jumlah device berdasarkan data di database, menentukan status setiap device berdasarkan data pemantauan terakhir, dan mengambil 5 alert terbaru. Semua data dikirim ke browser untuk ditampilkan.*
>
> *Untuk log real-time, di belakang server terdapat mesin monitoring yang secara berkala melakukan pengecekan ke setiap device menggunakan ICMP Ping, HTTP Check, atau TCP Port. Setiap kali ada hasil pengecekan baru, data langsung dikirim ke browser menggunakan WebSocket tanpa perlu refresh halaman. Log yang ditampilkan adalah 50 aktivitas terakhir, dan data disimpan di browser agar tidak hilang saat refresh.*
>
> *Manfaatnya, admin bisa memantau seluruh kondisi jaringan dalam satu tampilan, mendeteksi gangguan secara real-time, dan merespons masalah lebih cepat tanpa harus bolak-balik ke halaman lain."*

---

## 📌 Tampilan Dashboard (Layout)

```
┌─────────────────────────────────────────────────────────┐
│                    DASHBOARD                             │
├───────────────┬─────────────────┬───────────────────────┤
│  Total Device │     Online      │       Offline         │
│      12       │       10        │          2            │
│ sedang d...   │ perangkat a...  │   tidak merespon      │
├───────────────┴─────────────────┴───────────────────────┤
│                                                         │
│  ┌─────────────────────────────┐  ┌──────────────────┐  │
│  │   MONITORING ENGINE LOGS    │  │  System Status   │  │
│  │                             │  │                  │  │
│  │  [10:30:01] ONLINE Router-A │  │  Monitoring:     │  │
│  │   method: ICMP Ping seq:#1  │  │  Running ●       │  │
│  │   ttl: 64  +1.2ms           │  │                  │  │
│  │                             │  │  Interval:       │  │
│  │  [10:30:02] OFFLINE Server-B│  │  30 seconds ●    │  │
│  │   method: ICMP Ping seq:#5  │  │                  │  │
│  │   Request Timeout           │  │  Last Scan:      │  │
│  │                             │  │  10:30:15 ●      │  │
│  │  [10:30:03] ONLINE Switch-C │  │                  │  │
│  │   method: HTTP Check seq:#1 │  │  Notif:          │  │
│  │   ttl: 128  +45.7ms         │  │  Active ●        │  │
│  │                             │  │                  │  │
│  │  ...                        │  ├──────────────────┤  │
│  │                             │  │  Latest Alerts   │  │
│  │                             │  │  ● Server-B ...  │  │
│  │  [Clear]     50 events      │  │  ● Router-D ...  │  │
│  └─────────────────────────────┘  │  View all →      │  │
│                                   └──────────────────┘  │
└─────────────────────────────────────────────────────────┘
```
