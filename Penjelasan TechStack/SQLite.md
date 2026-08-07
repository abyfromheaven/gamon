# Teknologi: SQLite (modernc.org/sqlite v1.54.0)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| SQLite | Penyimpanan data perangkat, hasil monitoring (riwayat ping), alert, data pasangan Telegram, dan pengaturan aplikasi |

---

## 1. Peran Teknologi dalam Project GAMON

SQLite adalah sistem manajemen basis data yang digunakan GAMON untuk menyimpan seluruh data yang dihasilkan dan dibutuhkan aplikasi secara permanen. Berbeda dengan *engine* monitoring yang bekerja di memori untuk menampilkan status secara real-time, SQLite berfungsi sebagai tempat penyimpanan tiap data agar tidak hilang ketika aplikasi dimulai ulang.

Dalam project GAMON, SQLite mengelola lima tabel utama sebagai berikut:

| Tabel | Peran |
|-------|-------|
| `devices` | Menyimpan data perangkat jaringan (nama, tipe, IP, metode, interval, status, deskripsi) |
| `ping_history` | Menyimpan riwayat hasil pengecekan tiap perangkat (status, latency, ttl, waktu) |
| `alerts` | Menyimpan data alert/peringatan beserta status dan tingkat keparahannya |
| `telegram_pairing` | Menyimpan data pasangan akun Telegram (token, chat_id, status, masa berlaku) |
| `settings` | Menyimpan pengaturan aplikasi berbentuk pasangan kunci-nilai (key-value) |

Salah satu peran utama SQLite adalah menjadi "penopang" riwayat monitoring. Sebelum penggunaan database, hasil pengecekan hanya dikirim melalui WebSocket dan hilang saat koneksi terputus. Dengan SQLite, setiap hasil pengecekan tersimpan secara permanen sehingga mendukung analisis historis, tampilan dashboard, dan data alert.

**Bukti penggunaan:**
- `database/db.go` — inisialisasi koneksi, konfigurasi, dan *migration* tabel
- `database/models.go` — struktur data Go yang dipetakan ke tabel database
- `go.mod` — dependensi `modernc.org/sqlite` v1.54.0
- File `data/gamon.db` — file database hasil pembuatan aplikasi

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan SQLite dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

### a. Penyimpanan berbasis file dan tanpa konfigurasi

SQLite menyimpan seluruh data dalam satu file lokal (`data/gamon.db`) tanpa memerlukan *server database* terpisah atau konfigurasi rumit. Ini sangat cocok untuk aplikasi yang dijalankan pada satu server (seperti GAMON yang berjalan pada server NOC) karena sederhana untuk di-*deploy* dan tidak membutuhkan infrastruktur tambahan.

**Bukti:** pada `database/db.go` baris 13, koneksi dibuka langsung ke file `data/gamon.db` tanpa layanan server database.

### b. Cocok untuk aplikasi berukuran menengah dengan satu server

GAMON merupakan aplikasi dengan karakteristik penyimpanan data yang diakses terutama oleh satu server dan satu proses. SQLite sangat cocok untuk skenario ini karena tidak perlu *server database* tersendiri, cukup file yang dikelola langsung oleh proses aplikasi.

### c. Pemanfaatan driver pure-Go tanpa CGO

GAMON memakai `modernc.org/sqlite`, sebuah driver SQLite yang ditulis murni dalam Go (pure-Go) tanpa memerlukan kompilasi C (CGO). Ini memudahkan pembangunan (build) lintas platform karena tidak membutuhkan library C khusus. Bukti pada `go.mod` yang mencantumkan `modernc.org/sqlite`.

### d. Mode WAL untuk performa baca-tulis

SQLite dijalankan dengan mode WAL (Write-Ahead Logging) untuk meningkatkan performa saat terjadi banyak pembacaan dan penulisan secara bersamaan, yang sesuai dengan sifat aplikasi monitoring yang terus menyimpan hasil ping. Bukti pada `database/db.go` baris 20: `?_journal_mode=WAL`.

### e. Tanpa memakai ORM (Object-Relational Mapping)

GAMON memakai SQL langsung (raw SQL) melalui `database/sql` tanpa framework ORM. Ini memberi kendali penuh atas perintah yang dijalankan dan menjaga kesederhanaan aplikasi. Bukti pada seluruh operasi query berupa string SQL pada `database/db.go` dan berbagai *handler*.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Cara kerja SQLite dalam GAMON dapat diuraikan sebagai berikut.

### a. Inisialisasi koneksi

Saat aplikasi dimulai, fungsi `NewDB()` pada `database/db.go` membuka koneksi ke file database `data/gamon.db`. Fungsi ini:
1. Membuat direktori `data` apabila belum ada.
2. Membuka koneksi SQLite dengan mode `_journal_mode=WAL` dan `_busy_timeout=5000`.
3. Membatasi jumlah koneksi terbuka menjadi 1 (`SetMaxOpenConns(1)`) karena SQLite berbasis file dengan penguncian level file.
4. Memeriksa koneksi dan menjalankan *migration*.

### b. Migration (pembaruan struktur tabel)

Fungsi `migrate()` menjalankan serangkaian perintah SQL untuk membuat tabel `devices`, `ping_history`, `alerts`, `telegram_pairing`, dan `settings` apabila belum ada. Selain itu, dijalankan perintah `ALTER TABLE` untuk menambahkan kolom-kolom baru pada tabel yang sudah ada (misalnya kolom `status` pada `devices`, serta kolom `alert_type`, `acknowledged`, `acknowledged_at` pada `alerts`) secara tidak merusak data.

### c. Penyimpanan hasil monitoring

Setiap kali sebuah perangkat diping oleh monitoring *engine*, hasil tersebut disimpan ke tabel `ping_history` melalui perintah SQL `INSERT`. Data yang disimpan meliputi `device_id`, `status`, `latency_ms`, `ttl`, `seq`, `details`, dan `timestamp`. Menge note bahwa seluruh riwayat pengecekan terakumulasi di database.

### d. Penyimpanan dan pembaruan alert

Ketika engine mendeteksi perangkat bermasalah, sebuah baris alert ditambahkan ke tabel `alerts`. Saat perangkat kembali normal, data status alert diperbarui menjadi `resolved` dan waktu `resolved_at` diisi. Alert juga dapat ditandai `acknowledged` oleh pengguna melalui aplikasi.

### e. Penyimpanan pengaturan dan koneksi Telegram

Tabel `settings` menyimpan pengaturan seperti *failure threshold*, *check interval*, dan status aktifnya notifikasi dalam bentuk pasangan kunci-nilai. Sementara tabel `telegram_pairing` menyimpan informasi koneksi bot Telegram, termasuk token pairing dan `chat_id`, yang digunakan oleh *handler* dan *poller* Telegram.

### f. Integrasi relasi antar tabel

Tabel `ping_history` dan `alerts` memiliki relasi ke tabel `devices` melalui *foreign key* (`device_id`) dengan aturan `ON DELETE CASCADE`. Artinya, ketika sebuah perangkat dihapus, seluruh riwayat ping dan alert yang terkait dengan perangkat tersebut ikut terhapus secara otomatis, menjaga integritas dan kebersihan data.

### g. Pengaturan melalui kunci-nilai

Fungsi `GetSetting` dan `SetSetting` pada `database/db.go` digunakan untuk membaca dan menyimpan pengaturan aplikasi. `SetSetting` menggunakan `INSERT ... ON CONFLICT(key) DO UPDATE` sehingga menambahkan pengaturan baru atau memperbarui nilai jika kuncinya sudah ada.

---

## 4. Bukti Penggunaan

1. **File `go.mod`** — Dependensi `modernc.org/sqlite v1.54.0` sebagai driver database SQLite.

2. **File `database/db.go`** — Inisialisasi koneksi ke `data/gamon.db`, pengaturan WAL, batas koneksi, fungsi `migrate`, `GetSetting`, dan `SetSetting`.

3. **File `database/models.go`** — Struct Go untuk data `Device`, `PingHistory`, `Alert`, `MonitoringStatus`, dan lainnya yang dipetakan ke tabel.

4. **File `data/gamon.db`** — File hasil database nyata yang dihasilkan oleh aplikasi.

5. **Folder `handler/`** — Para *handler* memakai `*sql.DB` untuk melakukan baca-tulis data (contoh pada `device.go`, `alert.go`, `dashboard.go`, `monitoring.go`).

6. **File `main.go`** — Inisialisasi database lewat `database.NewDB()` dan diteruskan ke seluruh handler serta engine.

7. **File `monitor/engine.go`** — penulisan hasil ping ke `ping_history` dan alert.

---

## Catatan

- Versi SQLite yang di-cantumkan adalah versi driver `modernc.org/sqlite` (v1.54.0); versi inti SQLite-nya ikut dibundekan oleh driver tersebut.
- Database tidak memakai *server* terpisah; seluruh data berada dalam satu file, sehingga tepat untuk penyimpanan satu server GAMON.
- GAMON tidak memakai ORM; semua interaksi database menggunakan SQL mentah (`database/sql`).