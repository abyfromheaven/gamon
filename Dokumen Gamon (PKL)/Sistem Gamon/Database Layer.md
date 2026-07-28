# Database Layer — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian Database Layer

Database layer merupakan lapisan sistem yang bertanggung jawab untuk menyimpan, mengelola, dan mengambil data secara persisten. Dalam konteks Gamon, database layer berperan sebagai "otak penyimpanan" yang menjamin semua data perangkat, riwayat pengecekan, dan catatan gangguan tetap tersimpan meskipun aplikasi dimatikan atau server di-restart.

Sebelum adanya database layer, seluruh data Gamon hanya tersimpan di dalam memori komputer (RAM). Artinya, saat aplikasi ditutup atau server mati, semua data hilang dan harus dimulai dari awal. Dengan adanya database layer, data disimpan ke dalam file di penyimpanan permanen (disk) sehingga dapat diakses kembali kapan saja.

### 1.2 Peran Database Layer dalam Gamon

Database layer dalam Gamon memiliki beberapa peran penting:

| Peran | Penjelasan |
|---|---|
| **Penyimpanan Data Perangkat** | Menyimpan informasi lengkap tentang setiap perangkat yang didaftarkan (nama, tipe, IP, metode monitoring, lokasi, dll) |
| **Riwayat Pengecekan** | Mencatat setiap hasil pengecekan status perangkat beserta timestampnya |
| **Penyimpanan Alert** | Menyimpan catatan gangguan yang terjadi secara otomatis maupun manual |
| **Dukungan Query** | Menyediakan mekanisme pencarian dan filter data yang efisien |
| **Integritas Data** | Menjamin konsistensi data melalui mekanisme foreign key dan constraint |

### 1.3 Mengapa Database Diperlukan

Tanpa database, sistem monitoring Gamon hanya bersifat sementara (transient). Berikut perbandingan kondisi dengan dan tanpa database:

| Aspek | Tanpa Database | Dengan Database |
|---|---|---|
| Data perangkat | Hilang saat restart | Tersimpan permanen |
| Riwayat pengecekan | Hanya data terakhir | Riwayat lengkap |
| Alert history | Tidak tersimpan | Tersimpan dan dapat diakses |
| Pencarian data | Tidak memungkinkan | Query SQL yang powerful |
| Skalabilitas | Terbatas di memori | Dapat menangani data besar |

---

## 2. Pemilihan Teknologi

### 2.1 Mengapa SQLite

Untuk memilih sistem database, dilakukan evaluasi terhadap beberapa opsi:

| Database | Kelebihan | Kekurangan | Kesesuaian dengan Gamon |
|---|---|---|---|
| **SQLite** | Zero config, file-based, tidak perlu server terpisah | Write concurrency terbatas | Sangat sesuai untuk single-server NOC |
| PostgreSQL | Fitur lengkap, concurrent read/write tinggi | Butuh install server terpisah, overhead besar | Terlalu berat untuk prototype |
| MySQL | Populer, banyak dokumentasi | Butuh install server terpisah | Terlalu berat untuk prototype |
| MongoDB | Schema fleksibel | Butuh server terpisah, JSON-based | Tidak sesuai untuk data relasional |

**Alasan utama pemilihan SQLite:**

1. **Zero Configuration** — SQLite tidak memerlukan instalasi server database terpisah. Cukup dengan satu file `.db`, database langsung dapat digunakan. Ini sangat cocok untuk lingkungan NOC di mana server harus dapat di-deploy dengan cepat.

2. **File-Based** — Seluruh database tersimpan dalam satu file (`data/gamon.db`). File ini mudah dibackup, dipindahkan, atau dikloning hanya dengan menyalin file.

3. **Single-Server Deployment** — Gamon dirancang untuk berjalan di satu server NOC yang diakses oleh seluruh staf melalui browser. SQLite sangat optimal untuk arsitektur ini karena tidak memerlukan koneksi jaringan ke database server terpisah.

4. **Ringan dan Cepat** — SQLite memiliki overhead yang sangat rendah dibandingkan database client-server. Untuk aplikasi monitoring dengan puluhan hingga ratusan perangkat, SQLite mampu menangani beban kerja dengan baik.

5. **Standar Industri** — SQLite digunakan oleh jutaan aplikasi di seluruh dunia, termasuk sistem operasi mobile (Android/iOS), browser web, dan aplikasi desktop. Reliabilitasnya sudah terbukti.

### 2.2 Mengapa modernc.org/sqlite

Di Go, terdapat beberapa driver SQLite yang tersedia. Evaluasi dilakukan terhadap dua opsi utama:

| Driver | Tipe | Ketergantungan | Kelebihan | Kekurangan |
|---|---|---|---|---|
| `github.com/mattn/go-sqlite3` | CGO-based | Membutuhkan gcc/libc | Lebih mature, performa tinggi | Harus install gcc saat build, rumit di Windows |
| `modernc.org/sqlite` | Pure Go | Tidak ada ketergantungan C | Zero-dependency, mudah build dan deploy | Performa sedikit di bawah CGO-based |

**Alasan pemilihan modernc.org/sqlite:**

1. **CGO-Free** — Driver ini ditulis murni dalam bahasa Go tanpa menggunakan library C (CGO). Artinya, tidak perlu menginstal gcc atau toolchain C saat build aplikasi. Proses build menjadi lebih simpel dan cepat.

2. **Cross-Platform Tanpa Masalah** — Karena tidak bergantung pada library C, driver ini berjalan dengan mulus di berbagai sistem operasi (Linux, Windows, macOS) tanpa perlu konfigurasi khusus.

3. **Single Binary** — Hasil kompilasi Go berupa satu file executable. Dengan menggunakan driver pure Go, tidak perlu menyertakan library C tambahan. Deploy aplikasi menjadi sesederhana menyalin satu file.

4. **Sangat Mencukupi untuk Gamon** — Untuk aplikasi monitoring dengan ratusan perangkat, performa driver ini sudah lebih dari mencukupi. Perbedaan performa dengan CGO-based hanya terasa pada workload yang sangat berat (ribuan concurrent writes).

---

## 3. Arsitektur Database Layer

### 3.1 Posisi dalam Arsitektur Gamon

Database layer menempati posisi sentral dalam arsitektur Gamon, berada di antara backend logic dan penyimpanan data:

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
└───────┼────────────────┼────────────────────┼───────────┘
        │                │                    │
        │ REST API + WebSocket                 │
        │                │                    │
┌───────┼────────────────┼────────────────────┼───────────┐
│       ▼                ▼                    ▼           │
│  ┌──────────┐   ┌──────────────┐   ┌────────────────┐  │
│  │ REST API │   │ WebSocket    │   │ Handler Layer  │  │
│  │ Endpoints│   │ Hub          │   │ (CRUD Logic)   │  │
│  └────┬─────┘   └──────┬───────┘   └───────┬────────┘  │
│       │                │                    │           │
│       └────────────────┼────────────────────┘           │
│                        │                                │
│                        ▼                                │
│              ┌──────────────────┐                       │
│              │   DATABASE LAYER │  ◄── Fokus Dokumentasi│
│              │   (SQLite)       │       Ini             │
│              │                  │                       │
│              │  ┌────────────┐  │                       │
│              │  │ db.go      │  │                       │
│              │  │ models.go  │  │                       │
│              │  └────────────┘  │                       │
│              └────────┬─────────┘                       │
│                       │                                │
│                       ▼                                │
│              ┌──────────────────┐                       │
│              │   data/gamon.db  │                       │
│              │   (File SQLite)  │                       │
│              └──────────────────┘                       │
│                                                         │
│                      GO BACKEND                         │
│                      Port :8080                         │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Komponen Database Layer

Database layer terdiri dari dua komponen utama:

**1. `database/db.go` — Manajemen Koneksi dan Migrasi**

Komponen ini bertanggung jawab untuk:
- Membuka koneksi ke database SQLite
- Membuat direktori penyimpanan jika belum ada
- Menjalankan migrasi (pembuatan tabel)
- Mengatur konfigurasi koneksi (WAL mode, timeout, dll)

**2. `database/models.go` — Definisi Data**

Komponen ini mendefinisikan struktur data (struct) dalam bahasa Go yang merepresentasikan setiap tabel di database. Struct ini digunakan sebagai jembatan antara database dan bagian lain dari aplikasi.

### 3.3 Hubungan dengan Komponen Lain

```
┌─────────────┐
│   Handler   │ ──► Memanggil fungsi database untuk CRUD
│  (API REST) │     (Create, Read, Update, Delete devices & alerts)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Engine    │ ──► Menulis hasil pengecekan ke ping_history
│  (Monitor)  │     Membuat alert saat status berubah
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Database   │ ──► Menyimpan dan mengambil data
│   Layer     │     Menjamin integritas data
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ SQLite File │ ──► Penyimpanan permanen di disk
│ (gamon.db)  │
└─────────────┘
```

---

## 4. Entity Relationship Diagram (ERD)

### 4.1 Diagram Relasi

```
┌──────────────────────────────┐
│          devices             │
├──────────────────────────────┤
│ id (PK, INTEGER, AUTO)      │ ◄── Primary Key
│ name (TEXT, NOT NULL)        │
│ type (TEXT, NOT NULL)        │ ◄── Server/Router/Switch/AP/Website
│ ip (TEXT, NOT NULL)          │
│ url (TEXT, DEFAULT '')       │ ◄── Untuk HTTP Check (opsional)
│ port (INTEGER, NULL)         │ ◄── Untuk TCP Port (opsional)
│ method (TEXT, DEFAULT ICMP)  │ ◄── Metode monitoring
│ location (TEXT, DEFAULT '')  │ ◄── Lokasi (opsional)
│ check_interval (INT, DEF 3)  │ ◄── Interval detik
│ description (TEXT, DEF '')   │
│ created_at (DATETIME)        │
│ updated_at (DATETIME)        │
└──────────┬───────────────────┘
           │
           │ 1
           │
           │ *
┌──────────▼───────────────────┐
│       ping_history           │
├──────────────────────────────┤
│ id (PK, INTEGER, AUTO)      │ ◄── Primary Key
│ device_id (FK, INTEGER)      │ ◄── Foreign Key → devices.id
│ status (TEXT, NOT NULL)      │ ◄── online/offline/warning
│ latency_ms (REAL, DEFAULT 0) │ ◄── Waktu respons (ms)
│ ttl (INTEGER, DEFAULT 0)     │ ◄── Time-to-live (ICMP)
│ seq (INTEGER, DEFAULT 0)     │ ◄── Sequence number
│ details (TEXT, DEF '{}')     │ ◄── JSON fleksibel
│ timestamp (DATETIME)         │ ◄── Waktu pengecekan
└──────────────────────────────┘

┌──────────────────────────────┐
│           alerts             │
├──────────────────────────────┤
│ id (PK, INTEGER, AUTO)      │ ◄── Primary Key
│ device_id (FK, INTEGER)      │ ◄── Foreign Key → devices.id
│ title (TEXT, NOT NULL)       │ ◄── Judul alert
│ status (TEXT, DEF 'ongoing') │ ◄── ongoing/resolved
│ severity (TEXT, DEF 'low')   │ ◄── low/medium/high/critical
│ started_at (DATETIME)        │ ◄── Waktu mulai gangguan
│ resolved_at (DATETIME, NULL) │ ◄── Waktu resolved (null jika ongoing)
│ description (TEXT, DEF '')   │ ◄── Deskripsi gangguan
└──────────────────────────────┘
```

### 4.2 Penjelasan Relasi

| Relasi | Tipe | Penjelasan |
|---|---|---|
| `devices` → `ping_history` | One-to-Many | Satu device memiliki banyak riwayat pengecekan |
| `devices` → `alerts` | One-to-Many | Satu device memiliki banyak alert (gangguan) |
| Foreign Key Constraint | CASCADE | Hapus device → otomatis hapus semua ping_history dan alerts terkait |

### 4.3 Mengapa Relasi CASCADE?

Relasi `ON DELETE CASCADE` dipilih karena:

1. **Konsistensi Data** — Saat device dihapus, semua data terkait (riwayat pengecekan dan alert) juga harus dihapus untuk menjaga konsistensi.

2. **Kemudahan Manajemen** — Pengguna tidak perlu menghapus data terkait satu per satu secara manual.

3. **Mencegah Orphan Data** — Tanpa cascade, data ping_history dan alerts bisa kehilangan referensi ke device yang sudah dihapus (orphan data).

---

## 5. Struktur Tabel

### 5.1 Tabel `devices` — Buku Induk Perangkat

Tabel ini menyimpan informasi lengkap tentang setiap perangkat yang didaftarkan dalam sistem Gamon.

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | Identifikasi unik perangkat |
| `name` | TEXT | NOT NULL | Nama perangkat (contoh: "Router-01") |
| `type` | TEXT | NOT NULL | Tipe perangkat (Server/Router/Switch/AP/Website) |
| `ip` | TEXT | NOT NULL | Alamat IP perangkat |
| `url` | TEXT | DEFAULT '' | URL untuk HTTP Check (opsional) |
| `port` | INTEGER | NULL | Port untuk TCP Port Check (opsional) |
| `method` | TEXT | NOT NULL, DEFAULT 'ICMP Ping' | Metode monitoring yang digunakan |
| `location` | TEXT | DEFAULT '' | Lokasi fisik perangkat (opsional) |
| `check_interval` | INTEGER | DEFAULT 3 | Interval pengecekan dalam detik |
| `description` | TEXT | DEFAULT '' | Deskripsi tambahan tentang perangkat |
| `created_at` | DATETIME | DEFAULT CURRENT_TIMESTAMP | Waktu pendaftaran perangkat |
| `updated_at` | DATETIME | DEFAULT CURRENT_TIMESTAMP | Waktu pembaruan terakhir |

**Catatan tentang field opsional:**

- **`url`** — Hanya diisi untuk metode HTTP Check. Berguna untuk menentukan path yang akan di-check (contoh: `/health`, `/api/status`).

- **`port`** — Hanya diisi untuk metode TCP Port Check. Berguna untuk menentukan port yang akan di-check (contoh: `22` untuk SSH, `80` untuk HTTP, `3306` untuk MySQL).

- **`location`** — Bersifat opsional. Pengguna bebas ingin menambahkan informasi lokasi perangkat atau tidak. Berguna untuk identifikasi fisik perangkat di lapangan.

### 5.2 Tabel `ping_history` — Catatan Hasil Pengecekan

Tabel ini mencatat setiap hasil pengecekan yang dilakukan terhadap perangkat. Data ini bersifat time-series (data yang tercatat berdasarkan waktu).

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | Identifikasi unik catatan |
| `device_id` | INTEGER | NOT NULL, FOREIGN KEY | Relasi ke tabel devices |
| `status` | TEXT | NOT NULL | Status hasil pengecekan (online/offline/warning) |
| `latency_ms` | REAL | DEFAULT 0 | Waktu respons dalam milidetik |
| `ttl` | INTEGER | DEFAULT 0 | Time-to-live (untuk ICMP Ping) |
| `seq` | INTEGER | DEFAULT 0 | Sequence number (urutan pengecekan) |
| `details` | TEXT | DEFAULT '{}' | Data tambahan dalam format JSON |
| `timestamp` | DATETIME | DEFAULT CURRENT_TIMESTAMP | Waktu pengecekan dilakukan |

**Catatan tentang field `details`:**

Field `details` dirancang sebagai JSON yang fleksibel untuk menyimpan data spesifik sesuai metode monitoring:

| Metode | Contoh Isi `details` |
|---|---|
| ICMP Ping | `{"ttl": 64}` |
| HTTP Check | `{"status_code": 200, "content_type": "text/html"}` |
| TCP Port | `{"connected": true, "banner": "SSH-2.0-OpenSSH_8.2"}` |
| SNMP | `{"uptime": 86400, "cpu_usage": 45.2}` |

Desain ini memungkinkan penambahan metode monitoring baru tanpa perlu mengubah struktur tabel.

### 5.3 Tabel `alerts` — Catatan Gangguan

Tabel ini menyimpan catatan setiap gangguan yang terjadi pada perangkat. Alert dibuat secara otomatis oleh sistem saat mendeteksi perubahan status.

| Kolom | Tipe | Constraint | Keterangan |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY AUTOINCREMENT | Identifikasi unik alert |
| `device_id` | INTEGER | NOT NULL, FOREIGN KEY | Relasi ke tabel devices |
| `title` | TEXT | NOT NULL | Judul alert (contoh: "Device unreachable") |
| `status` | TEXT | NOT NULL, DEFAULT 'ongoing' | Status alert (ongoing/resolved) |
| `severity` | TEXT | NOT NULL, DEFAULT 'low' | Tingkat keparahan (low/medium/high/critical) |
| `started_at` | DATETIME | DEFAULT CURRENT_TIMESTAMP | Waktu gangguan terjadi |
| `resolved_at` | DATETIME | NULL | Waktu gangguan terselesaikan (null jika masih ongoing) |
| `description` | TEXT | DEFAULT '' | Deskripsi detail gangguan |

**Penjelasan Status Alert:**

| Status | Penjelasan |
|---|---|
| `ongoing` | Gangguan masih berlangsung, perangkat belum normal |
| `resolved` | Gangguan telah terselesaikan, perangkat sudah normal kembali |

**Penjelasan Severity Level:**

| Level | Penjelasan | Contoh |
|---|---|---|
| `low` | Gangguan ringan, tidak mendesak | Backup selesai, info normal |
| `medium` | Gangguan sedang, perlu perhatian | DNS lambat, suhu tinggi |
| `high` | Gangguan serius, perlu tindakan segera | Packet loss tinggi, interface flapping |
| `critical` | Gangguan kritis, perlu tindakan segera | Device unreachable, service down |

---

## 6. Fungsi-Fungsi Database

### 6.1 `NewDB()` — Inisialisasi Koneksi

Fungsi `NewDB()` merupakan titik masuk utama untuk menginisialisasi database. Fungsi ini melakukan serangkaian langkah secara berurutan:

**Alur Eksekusi:**

```
NewDB() dipanggil
       │
       ▼
1. Buat direktori "data/" jika belum ada
   Tujuan: Memastikan lokasi penyimpanan database tersedia
       │
       ▼
2. Buka koneksi ke SQLite
   File: data/gamon.db
   Parameter: WAL mode, busy timeout 5 detik
       │
       ▼
3. Atur batas koneksi
   MaxOpenConns: 1 (SQLite hanya boleh 1 writer)
   MaxIdleConns: 1
       │
       ▼
4. Test koneksi (Ping)
   Tujuan: Memastikan database dapat diakses
       │
       ▼
5. Jalankan migrate()
   Tujuan: Membuat tabel jika belum ada
       │
       ▼
6. Kembalikan koneksi database
   Digunakan oleh komponen lain untuk operasi CRUD
```

**Mengapa MaxOpenConns(1)?**

SQLite memiliki keterbatasan teknis: hanya boleh ada satu operasi penulisan (write) dalam satu waktu. Jika ada dua goroutine mencoba menulis secara bersamaan, SQLite akan mengembalikan error atau deadlock. Dengan mengatur `MaxOpenConns(1)`, Go memastikan semua operasi database dijalankan secara berurutan (serialized), sehingga tidak terjadi konflik penulisan.

### 6.2 `migrate()` — Pembuatan Tabel

Fungsi `migrate()` bertanggung jawab untuk membuat struktur tabel di database. Fungsi ini menggunakan perintah SQL `CREATE TABLE IF NOT EXISTS` yang berarti:

- Jika tabel **belum ada** → buat tabel baru
- Jika tabel **sudah ada** → lewati (tidak error)

**Tabel yang dibuat:**

| Urutan | Tabel | Fungsi |
|---|---|---|
| 1 | `devices` | Menyimpan data perangkat |
| 2 | `ping_history` | Menyimpan riwayat pengecekan |
| 3 | `alerts` | Menyimpan catatan gangguan |

**Index yang dibuat:**

| Urutan | Index | Kolom | Tujuan |
|---|---|---|---|
| 4 | `idx_ping_history_device_id` | `ping_history.device_id` | Mempercepat query riwayat per device |
| 5 | `idx_ping_history_timestamp` | `ping_history.timestamp` | Mempercepat query berdasarkan waktu |
| 6 | `idx_alerts_device_id` | `alerts.device_id` | Mempercepat query alert per device |
| 7 | `idx_alerts_status` | `alerts.status` | Mempercepat filter berdasarkan status |

---

## 7. Konfigurasi Database

### 7.1 WAL Mode (Write-Ahead Logging)

WAL (Write-Ahead Logging) adalah mode penulisan SQLite yang memungkinkan operasi baca (read) dan tulis (write) dilakukan secara bersamaan.

**Cara Kerja WAL:**

```
Mode Tanpa WAL (Default):
┌─────────────────────────────────────────┐
│  Write 1 ──► Blokir semua Read          │
│  Write 2 ──► Tunggu Write 1 selesai     │
│  Read    ──► Tunggu semua Write selesai  │
└─────────────────────────────────────────┘

Mode Dengan WAL:
┌─────────────────────────────────────────┐
│  Write 1 ──► Tulis ke WAL file          │
│  Read    ──► Baca dari database utama    │ ◄── Bisa bersamaan!
│  Write 2 ──► Tulis ke WAL file          │
│  WAL Sync ──► Sync ke database utama     │
└─────────────────────────────────────────┘
```

**Manfaat WAL untuk Gamon:**

- **Concurrent Read/Write** — Sambil engine menulis hasil pengecekan ke database (setiap 3 detik per device), frontend dapat membaca data untuk ditampilkan di dashboard tanpa terganggu.
- **Performa Lebih Baik** — WAL menulis ke file log terlebih dahulu, lalu di-sync ke database utama secara berkala. Ini lebih efisien daripada langsung menulis ke database utama.
- **Crash Recovery** — Jika aplikasi crash mendadak, data di WAL file masih dapat dipulihkan.

### 7.2 Busy Timeout

Parameter `_busy_timeout=5000` memberitahu SQLite untuk menunggu hingga 5 detik jika ada kunci (lock) yang sedang dipegang oleh operasi lain.

```
Tanpa Busy Timeout:
┌─────────────────────────────────────────┐
│  Operasi A sedang menulis               │
│  Operasi B mencoba menulis              │
│  → Langsung ERROR: database locked      │
└─────────────────────────────────────────┘

Dengan Busy Timeout 5 detik:
┌─────────────────────────────────────────┐
│  Operasi A sedang menulis               │
│  Operasi B mencoba menulis              │
│  → Operasi B MENUNGGU (max 5 detik)     │
│  → Operasi A selesai                    │
│  → Operasi B menjalankan                │
└─────────────────────────────────────────┘
```

### 7.3 MaxOpenConns dan MaxIdleConns

| Parameter | Nilai | Penjelasan |
|---|---|---|
| `MaxOpenConns` | 1 | Batas maksimum koneksi database yang dibuka secara bersamaan |
| `MaxIdleConns` | 1 | Batas maksimum koneksi idle (tidak aktif) yang disimpan |

Pengaturan ini memastikan SQLite hanya melayani satu operasi penulisan dalam satu waktu, yang merupakan batasan fundamental SQLite.

---

## 8. Index dan Optimasi

### 8.1 Pengertian Index

Index dalam database analogisnya seperti indeks di akhir buku. Tanpa indeks, untuk menemukan informasi tertentu, kita harus membuka semua halaman satu per satu. Dengan indeks, kita langsung tahu halaman mana yang harus dibuka.

### 8.2 Index dalam Gamon

| Index | Tabel | Kolom | Alasan Pembuatan |
|---|---|---|---|
| `idx_ping_history_device_id` | ping_history | device_id | Query "ambil riwayat pengecekan device X" menjadi sangat cepat |
| `idx_ping_history_timestamp` | ping_history | timestamp | Query "ambil riwayat pengecekan pada waktu tertentu" menjadi cepat |
| `idx_alerts_device_id` | alerts | device_id | Query "ambil semua alert untuk device X" menjadi cepat |
| `idx_alerts_status` | alerts | status | Filter "tampilkan semua alert ongoing" menjadi cepat |

### 8.3 Dampak Performa

Tanpa index, setiap query harus melakukan **full table scan** (membaca semua baris). Sebagai ilustrasi:

```
Tanpa Index (Full Table Scan):
┌─────────────────────────────────────────┐
│  Query: "Cari ping_history device_id=5" │
│                                         │
│  Baris 1: device_id=1 → Skip           │
│  Baris 2: device_id=3 → Skip           │
│  Baris 3: device_id=1 → Skip           │
│  ... (ribuan baris lainnya)             │
│  Baris N: device_id=5 → MATCH!          │
│                                         │
│  Total: harus baca SEMUA baris          │
└─────────────────────────────────────────┘

Dengan Index:
┌─────────────────────────────────────────┐
│  Query: "Cari ping_history device_id=5" │
│                                         │
│  Index langsung tunjuk: "Baris ke-N"    │
│  → Langsung baca baris yang dicari      │
│                                         │
│  Total: hanya baca 1 baris              │
└─────────────────────────────────────────┘
```

Untuk Gamon yang mencatat hasil pengecekan setiap 3 detik per device, jumlah data di `ping_history` akan bertambah sangat cepat. Tanpa index, query untuk menampilkan riwayat satu device akan melambat seiring bertambahnya data.

---

## 9. Cascade Rules

### 9.1 Pengertian Foreign Key Cascade

Foreign key constraint memastikan integritas relasi antar tabel. CASCADE adalah opsi yang menentukan apa yang terjadi pada data terkait saat data induk dihapus.

### 9.2 Implementasi dalam Gamon

```
Situasi: Device "Router-01" (id=1) dihapus

Tanpa CASCADE:
┌─────────────────────────────────────────┐
│  devices: Router-01 DIHAPUS             │
│  ping_history: device_id=1 masih ada    │ ◄── ORPHAN DATA!
│  alerts: device_id=1 masih ada          │ ◄── ORPHAN DATA!
└─────────────────────────────────────────┘

Dengan CASCADE:
┌─────────────────────────────────────────┐
│  devices: Router-01 DIHAPUS             │
│  ping_history: device_id=1 IKUT HAPUS   │ ◄── Konsisten!
│  alerts: device_id=1 IKUT HAPUS         │ ◄── Konsisten!
└─────────────────────────────────────────┘
```

### 9.3 Manfaat CASCADE untuk Gamon

1. **Konsistensi Data** — Tidak ada data yang kehilangan referensi ke device yang sudah dihapus.

2. **Kemudahan Penggunaan** — Pengguna cukup menghapus device, semua data terkait otomatis terhapus.

3. **Mencegah Error** — Query yang melibatkan join antara device dan data terkait tidak akan error karena foreign key violation.

---

## 10. Flow Data

### 10.1 Flow Penambahan Device

```
User mengisi form di Device Management
       │
       ▼
Frontend mengirim POST /api/devices
       │
       ▼
Handler menerima request
       │
       ▼
Validasi data (nama, IP, metode, dll)
       │
       ▼
Simpan ke tabel devices
       │
       ▼
Kirim response sukses ke frontend
       │
       ▼
Engine mulai monitoring device baru
```

### 10.2 Flow Pengecekan Device

```
Engine menjalankan checkLoop untuk device
       │
       ▼
Panggil ICMPCheck(deviceID, ip, seq)
       │
       ▼
Jalankan perintah: ping -c 1 -W 3 <ip>
       │
       ▼
Parse output, dapatkan status + latency
       │
       ▼
Simpan hasil ke tabel ping_history
       │
       ▼
Broadcast via WebSocket ke semua client
       │
       ▼
Frontend update tampilan secara real-time
```

### 10.3 Flow Pembuatan Alert

```
Engine mendeteksi status berubah (contoh: online → offline)
       │
       ▼
Increment consecutiveFailures
       │
       ▼
Jika consecutiveFailures >= 3:
       │
       ▼
Buat record baru di tabel alerts
(title: "Device unreachable", status: ongoing, severity: critical)
       │
       ▼
Broadcast "status_change" via WebSocket
       │
       ▼
Dashboard dan Alert Center update
```

### 10.4 Flow Resolusi Alert

```
Device kembali online setelah sebelumnya offline
       │
       ▼
Engine mendeteksi status: online (sebelumnya offline)
       │
       ▼
Reset consecutiveFailures ke 0
       │
       ▼
Cari alert ongoing untuk device ini
       │
       ▼
Update status alert: ongoing → resolved
Isi resolved_at dengan waktu sekarang
       │
       ▼
Broadcast "status_change" via WebSocket
       │
       ▼
Dashboard dan Alert Center update
```

---

## 11. Struktur File

### 11.1 Letak File Database

```
/home/aby/gamon/
├── database/
│   ├── db.go          # Koneksi + migrasi
│   └── models.go      # Definisi struct
├── data/
│   └── gamon.db       # File database SQLite
├── handler/           # REST API endpoints
├── monitor/           # Engine monitoring
├── main.go            # Entry point
└── frontend/          # React application
```

### 11.2 Penjelasan File

| File | Fungsi | Analogi |
|---|---|---|
| `database/db.go` | Membuka koneksi, menjalankan migrasi, mengatur konfigurasi | Tukang bangun yang membangun fondasi |
| `database/models.go` | Mendefinisikan struktur data dalam bahasa Go | Cetak biru / arsitektur |
| `data/gamon.db` | File database SQLite (otomatis dibuat) | Gedung yang sudah jadi |

---

## 12. Referensi Kode

### 12.1 Import Statement

```go
package database

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "path/filepath"

    _ "modernc.org/sqlite"
)
```

- `database/sql` — Package standar Go untuk bekerja dengan database SQL
- `modernc.org/sqlite` — Driver SQLite pure Go (di-import dengan `_` karena hanya perlu driver-nya, bukan fungsi langsung)

### 12.2 Struct Definition (models.go)

```go
type Device struct {
    ID             int       `json:"id"`
    Name           string    `json:"name"`
    Type           string    `json:"type"`
    IP             string    `json:"ip"`
    // ... field lainnya
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
```

- Struct ini merepresentasikan satu baris data di tabel `devices`
- Tag `json:"id"` menentukan nama field saat di-serialize ke JSON untuk API response

---

## 13. Kesimpulan

Database layer Gamon dirancang dengan prinsip:

1. **Kesederhanaan** — Menggunakan SQLite yang zero-config dan file-based
2. **Reliabilitas** — WAL mode dan busy timeout untuk menjamin konsistensi data
3. **Fleksibilitas** — Field `details` JSON yang bisa menampung data metode monitoring apapun
4. **Integritas** — Foreign key dengan cascade untuk menjaga konsistensi relasi
5. **Performa** index untuk mempercepat query pada data yang sering diakses

Dengan fondasi database yang solid, seluruh komponen Gamon (handler, engine, frontend) dapat bekerja dengan data yang konsisten dan tersimpan secara permanen.

---

*Dokumentasi ini merupakan bagian dari laporan Praktek Kerja Lapangan (PKL) Gamon — Garda Monitoring.*
