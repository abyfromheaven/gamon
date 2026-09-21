# Device Management — Halaman Manajemen Perangkat

---

## 1. HALAMAN INI UNTUK APA? (Tujuan / Fungsi)

Device Management adalah halaman yang berfungsi sebagai **pusat pengelolaan seluruh perangkat jaringan** yang sedang dipantau oleh sistem GAMON. Di halaman ini, network admin bisa **menambah, mengubah, menghapus, dan mengaktifkan/menonaktifkan** perangkat jaringan yang akan dipantau.

**Tujuan utamanya:**

- **Mengelola daftar perangkat** — Admin bisa melihat semua perangkat yang terdaftar dalam satu tabel
- **Menambah perangkat baru** — Admin bisa mendaftarkan perangkat baru ke dalam sistem pemantauan
- **Mengubah informasi perangkat** — Admin bisa memperbarui data perangkat seperti nama, IP, lokasi, dan interval pengecekan
- **Menghapus perangkat** — Admin bisa menghapus perangkat yang sudah tidak perlu dipantau
- **Mengaktifkan/menonaktifkan perangkat** — Admin bisa menjeda atau melanjutkan pemantauan pada perangkat tertentu tanpa menghapusnya

**Kenapa Device Management penting?**

Device Management adalah fondasi dari seluruh sistem pemantauan. Tanpa perangkat yang terdaftar dan dikonfigurasi dengan benar, sistem monitoring tidak akan bisa bekerja. Halaman ini memastikan admin memiliki kendali penuh atas perangkat mana saja yang sedang dipantau dan bagaimana cara pemantauannya.

---

## 2. USER BISA MELAKUKAN APA? (Aksi yang Dilakukan Network Admin)

Di halaman Device Management, network admin bisa melakukan **7 hal utama**:

---

### A. Melihat Daftar Seluruh Perangkat

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilihat** | Tabel berisi semua perangkat jaringan yang sudah terdaftar di sistem |
| **Dimana letaknya** | Area utama di bawah search bar, menempati seluruh lebar halaman |
| **Informasi yang ditampilkan** | Status (Active/Inactive), Nama device, Tipe, Alamat IP, Metode monitoring, Port, Lokasi |
| **Format tampilan** | Desktop: tabel dengan kolom-kolom. Mobile: kartu-kartu vertikal |

**Keterangan kolom di tabel:**

| Kolom | Penjelasan |
|-------|------------|
| **Status** | Indikator apakah perangkat sedang aktif dipantau (hijau) atau nonaktif (abu-abu) |
| **Name** | Nama perangkat yang diberikan oleh admin (contoh: "Router Utama") |
| **Type** | Jenis perangkat: Server, Router, Switch, Access Point, atau Website |
| **IP Address** | Alamat IP perangkat yang dipantau (contoh: 192.168.1.1) |
| **Method** | Metode pemantauan yang digunakan (contoh: ICMP Ping) |
| **Port** | Nomor port perangkat (jika menggunakan HTTP Check atau TCP Port) |
| **Location** | Lokasi fisik perangkat (opsional) |

---

### B. Mencari Perangkat (Search)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Mencari perangkat berdasarkan nama atau alamat IP |
| **Dimana letaknya** | Search bar di bagian atas halaman |
| **Cara kerja** | Admin mengetik nama atau IP, sistem secara otomatis memfilter daftar perangkat |
| **Sifat pencarian** | Tidak case-sensitive (huruf besar/kecil tidak mempengaruhi hasil) |

---

### C. Memfilter Perangkat Berdasarkan Tipe

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Memfilter perangkat hanya menampilkan tipe tertentu |
| **Dimana letaknya** | Dropdown filter di sebelah search bar |
| **Pilihan filter** | All (semua), Server, Router, Switch, Access Point, Website |
| **Sifat filter** | Dinamis — hanya menampilkan tipe yang sudah ada di database |

---

### D. Menambah Perangkat Baru (Add Device)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Mendaftarkan perangkat baru ke dalam sistem pemantauan |
| **Dimana tombolnya** | Tombol "Add Device" di pojok kanan atas halaman |
| **Form yang diisi** | Device Name, Device Type, IP Address, Check Interval, Location (opsional), Description (opsional) |

**Detail field di form Add Device:**

| Field | Wajib? | Penjelasan |
|-------|--------|------------|
| **Device Name** | Ya | Nama perangkat (contoh: "Router Utama") |
| **Device Type** | Ya | Jenis perangkat: Server, Router, Switch, Access Point, Website. Ada fitur autocomplete/suggestion |
| **IP Address** | Ya | Alamat IP perangkat (contoh: 192.168.1.1) |
| **Check Interval** | Ya | Jarak waktu pengecekan dalam detik (minimal 1 detik, default: 3 detik) |
| **Location** | Tidak | Lokasi fisik perangkat (contoh: "Gedung A Lantai 2") |
| **Description** | Tidak | Keterangan tambahan tentang perangkat |

---

### E. Mengubah Informasi Perangkat (Edit Device)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Memperbarui data perangkat yang sudah terdaftar |
| **Dimana tombolnya** | Tombol ikon pensil di kolom Actions pada tabel |
| **Form yang diisi** | Sama dengan form Add Device, tapi sudah terisi dengan data lama |
| **Yang bisa diubah** | Nama, tipe, IP, interval pengecekan, lokasi, deskripsi |

---

### F. Menghapus Perangkat (Delete Device)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Menghapus perangkat dari sistem secara permanen |
| **Dimana tombolnya** | Tombol ikon tempat sampah di kolom Actions pada tabel |
| **Proses** | Muncul dialog konfirmasi "Are you sure you want to delete [nama device]?" |
| **Yang terjadi** | Perangkat dihapus dari database DAN pemantauan dihentikan |

---

### G. Mengaktifkan/Menonaktifkan Perangkat (Toggle Status)

| Aspek | Penjelasan |
|-------|------------|
| **Apa yang dilakukan** | Menjeda atau melanjutkan pemantauan pada perangkat tanpa menghapusnya |
| **Dimana tombolnya** | Tombol ikon daya (power) di kolom Actions pada tabel |
| **Pilihan** | **Active** (hijau): Perangkat sedang dipantau. **Inactive** (abu-abu): Pemantauan dijeda |
| **Kapan digunakan** | Saat admin ingin menonaktifkan sementara perangkat yang sedang maintenance |

---

## 3. SISTEM MELAKUKAN APA? (Proses di Belakang Aksi)

Berikut penjelasan proses yang terjadi di belakang layar untuk setiap aksi:

---

### A. Saat Admin Membuka Halaman Device Management

**Proses yang terjadi:**

1. **Browser meminta daftar perangkat** — Browser mengirim permintaan ke server untuk mendapatkan semua data perangkat
2. **Server mengambil data dari database** — Server mengambil semua data perangkat dari tabel devices, diurutkan dari yang paling baru ditambahkan
3. **Data dikirim ke browser** — Seluruh data perangkat dikirim ke browser untuk ditampilkan dalam bentuk tabel
4. **Jumlah device dan active device ditampilkan** — Di bagian header halaman, ditampilkan total perangkat dan jumlah yang aktif

---

### B. Saat Admin Menambah Perangkat Baru

**Proses yang terjadi:**

1. **Admin mengisi form** — Admin mengisi nama, tipe, IP, interval, dan data lainnya di form Add Device
2. **Browser memvalidasi input** — Sistem memastikan nama dan IP sudah diisi, interval minimal 1 detik
3. **Browser mengirim data ke server** — Data perangkat baru dikirim ke server
4. **Server memvalidasi data** — Server memastikan nama tidak kosong, tipe tidak kosong, IP tidak kosong, dan status harus active atau inactive
5. **Server menyimpan ke database** — Data perangkat baru disimpan ke tabel devices
6. **Server memulai pemantauan** — Jika status perangkat adalah active, server secara otomatis memulai mesin monitoring untuk perangkat tersebut
7. **Browser menerima konfirmasi** — Server mengirim data perangkat yang baru disimpan ke browser
8. **Daftar perangkat diperbarui** — Perangkat baru langsung muncul di tabel

---

### C. Saat Admin Mengubah Informasi Perangkat

**Proses yang terjadi:**

1. **Admin membuka form edit** — Admin menekan tombol edit, form muncul dengan data perangkat yang sudah ada
2. **Admin mengubah data** — Admin memodifikasi field yang ingin diubah
3. **Browser mengirim perubahan ke server** — Data perubahan dikirim ke server
4. **Server mengambil data lama** — Server mengambil data perangkat saat ini dari database untuk perbandingan
5. **Server memperbarui database** — Field yang diubah diperbarui di database, field yang tidak diubah tetap mempertahankan nilai lama
6. **Server restart pemantauan** — Server menghentikan pemantauan lama dan memulai ulang dengan konfigurasi baru (IP, metode, interval yang sudah diperbarui)
7. **Browser menerima konfirmasi** — Data perangkat yang sudah diperbarui dikirim ke browser
8. **Tabel diperbarui** — Perubahan langsung terlihat di tabel

---

### D. Saat Admin Menghapus Perangkat

**Proses yang terjadi:**

1. **Admin menekan tombol hapus** — Muncul dialog konfirmasi untuk memastikan admin yakin ingin menghapus
2. **Admin mengonfirmasi** — Admin menekan tombol "Delete" di dialog konfirmasi
3. **Browser mengirim permintaan hapus** — Browser mengirim permintaan hapus ke server
4. **Server menghapus dari database** — Data perangkat dihapus permanen dari tabel devices
5. **Server menghentikan pemantauan** — Mesin monitoring untuk perangkat tersebut dihentikan
6. **Browser menerima konfirmasi** — Server mengirim pesan sukses ke browser
7. **Tabel diperbarui** — Perangkat yang dihapus langsung hilang dari tabel

---

### E. Saat Admin Mengaktifkan/Menonaktifkan Perangkat

**Proses yang terjadi:**

1. **Admin menekan tombol toggle** — Status perangkat berubah (active → inactive atau sebaliknya)
2. **Browser mengirim perubahan status** — Status baru dikirim ke server
3. **Server memperbarui database** — Status perangkat diperbarui di database
4. **Server mengelola pemantauan:**
   - Jika berubah ke **active**: Server mengambil data perangkat dari database lalu memulai mesin monitoring
   - Jika berubah ke **inactive**: Server menghentikan mesin monitoring untuk perangkat tersebut
5. **Browser menerima konfirmasi** — Status baru dikonfirmasi ke browser
6. **Tabel diperbarui** — Indikator status berubah warna (hijau untuk active, abu-abu untuk inactive)

---

### F. Saat Admin Mencari atau Memfilter Perangkat

**Proses yang terjadi:**

1. **Admin mengetik di search bar atau memilih filter** — Sistem secara langsung memfilter daftar perangkat
2. **Pencocokan dilakukan di browser** — Filter dan pencarian dilakukan di sisi browser (client-side), bukan ke server
3. **Hasil ditampilkan secara instan** — Tabel langsung berubah menampilkan hanya perangkat yang cocok dengan pencarian atau filter

---

## 4. HASILNYA APA? (Output dan Manfaat)

| No | Output | Manfaat untuk Network Admin |
|----|--------|-----------------------------|
| 1 | **Tabel daftar perangkat** | Admin bisa melihat semua perangkat dalam satu tampilan tanpa perlu query database |
| 2 | **Informasi lengkap perangkat** | Admin tahu detail setiap perangkat (nama, tipe, IP, lokasi) dalam satu baris |
| 3 | **Indikator status** | Admin langsung tahu perangkat mana yang aktif dipantau dan mana yang tidak |
| 4 | **Form tambah perangkat** | Admin bisa mendaftarkan perangkat baru dalam hitungan detik |
| 5 | **Form edit perangkat** | Admin bisa memperbarui data perangkat tanpa perlu menghapus dan membuat ulang |
| 6 | **Dialog konfirmasi hapus** | Admin tidak akan salah menghapus perangkat karena ada konfirmasi |
| 7 | **Toggle status** | Admin bisa menjeda pemantauan sementara (misal saat maintenance) tanpa menghapus perangkat |
| 8 | **Pencarian dan filter** | Admin bisa cepat menemukan perangkat tertentu dari ratusan perangkat |
| 9 | **Otomatis mulai/stop monitoring** | Admin tidak perlu setting monitoring secara terpisah, cukup tambah perangkat dan monitoring langsung berjalan |
| 10 | **Responsive design** | Admin bisa mengelola perangkat dari desktop maupun mobile |

---

## 🧠 Rangkuman untuk Sidang

> *"Halaman Device Management berfungsi sebagai pusat pengelolaan seluruh perangkat jaringan yang dipantau oleh sistem. Di halaman ini, admin bisa melihat daftar lengkap semua perangkat dalam satu tabel, mencari perangkat berdasarkan nama atau IP, dan memfilter berdasarkan jenis perangkat.*
>
> *Admin bisa menambah perangkat baru dengan mengisi form yang berisi nama, tipe, alamat IP, interval pengecekan, lokasi, dan deskripsi. Saat perangkat baru ditambahkan dengan status active, sistem secara otomatis memulai mesin monitoring untuk perangkat tersebut tanpa perlu pengaturan tambahan.*
>
> *Selain itu, admin bisa mengubah informasi perangkat yang sudah ada, menghapus perangkat yang sudah tidak diperlukan, serta mengaktifkan atau menonaktifkan pemantauan pada perangkat tertentu. Saat status perangkat diubah, sistem otomatis menyesuaikan — jika diaktifkan, pemantauan dimulai; jika dinonaktifkan, pemantauan dihentikan.*
>
> *Proses di belakang layar: setiap aksi yang dilakukan admin dikirim ke server, divalidasi, disimpan ke database, dan diterapkan ke mesin monitoring. Semua perubahan langsung terlihat di tabel tanpa perlu refresh halaman. Pencarian dan filter dilakukan di browser sehingga hasilnya instan."*

---

## 📌 Tampilan Device Management (Layout)

```
┌─────────────────────────────────────────────────────────┐
│  Device Management                                      │
│  12 devices registered · 10 active        [+ Add Device]│
├─────────────────────────────────────────────────────────┤
│  [🔍 Search devices...        ] [Filter: All ▼]        │
├─────────────────────────────────────────────────────────┤
│  ┌────────┬──────────┬────────┬───────────┬──────┬───┐  │
│  │ Status │ Name     │ Type   │ IP        │Method│Act│  │
│  ├────────┼──────────┼────────┼───────────┼──────┼───┤  │
│  │ ● Act  │ Router-A │ Router │192.168.1.1│ ICMP │⚙✏🗑│  │
│  │ ● Act  │ Server-B │ Server │10.0.0.5   │ ICMP │⚙✏🗑│  │
│  │ ○ Ina  │ Switch-C │ Switch │192.168.2.1│ ICMP │⚙✏🗑│  │
│  │ ● Act  │ Web-D    │Website │google.com │ HTTP │⚙✏🗑│  │
│  └────────┴──────────┴────────┴───────────┴──────┴───┘  │
│                                                         │
│  ⚙ = Toggle Status  ✏ = Edit  🗑 = Delete              │
└─────────────────────────────────────────────────────────┘
```
