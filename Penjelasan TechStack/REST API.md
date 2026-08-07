# Teknologi: REST API (Representational State Transfer)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| REST API | Komunikasi data antara frontend dan backend: pengelolaan perangkat, alert, dashboard, monitoring, pengaturan, dan koneksi Telegram |

---

## 1. Peran Teknologi dalam Project GAMON

REST API (Representational State Transfer) merupakan gaya arsitektur komunikasi antar perangkat lunak yang digunakan GAMON sebagai jalur utama pertukaran data antara *frontend* React dan *backend* Go. Melalui REST API, seluruh operasi pengelolaan data yang dilakukan pengguna di antarmuka web diteruskan ke server, diproses, lalu hasilnya dikembalikan dalam format JSON.

Dalam konteks GAMON, REST API berperan sebagai "jembatan" untuk operasi yang bersifat *request-response* (permintaan-jawaban), yaitu operasi yang tidak membutuhkan pembaruan langsung secara terus-menerus. REST API menangani seluruh operasi CRUD (Create, Read, Update, Delete) dan query data, yang dirinci sebagai berikut.

| Kelompok Endpoint | Fungsi |
|-------------------|--------|
| `/api/devices` | Membaca daftar perangkat, menambah perangkat baru, mengubah, menghapus, serta mengendalikan proses monitoring (start/stop/status) |
| `/api/alerts` | Membaca daftar alert, melihat detail, menyelesaikan (resolve), dan menandai alert sebagai diketahui (acknowledge) |
| `/api/dashboard` | Mengambil data ringkasan dashboard (total, online, offline, alert terbaru) |
| `/api/monitoring` | Mengambil status monitoring terkini dan riwayat hasil ping per perangkat |
| `/api/settings` | Membaca dan mengubah pengaturan aplikasi (threshold, interval, notifikasi) |
| `/api/telegram/*` | Menghubungkan dan memutuskan koneksi bot Telegram |
| `/api/health` | Memeriksa ketersediaan server (health check) |

Seluruh endpoint REST disediakan oleh *backend* Go dengan memanfaatkan paket bawaan `net/http`, sedangkan *frontend* React memanggilnya melalui sebuah lapisan *API client* pada file `frontend/src/lib/api.ts`.

**Bukti penggunaan:**
- `main.go` — pendaftaran seluruh rute `/api/*` pada `http.NewServeMux`
- Folder `handler/` — `device.go`, `alert.go`, `dashboard.go`, `monitoring.go`, `settings.go`, `telegram.go`, `helpers.go`
- `frontend/src/lib/api.ts` — fungsi-fungsi pemanggil API dari sisi frontend

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan REST API sebagai jalur komunikasi data dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

### a. Sesuai kebutuhan operasi data yang bersifat request-response

Sebagian besar interaksi pengguna terhadap GAMON merupakan operasi data yang tidak perlu didorong server secara terus-menerus, seperti menambah perangkat, mengubah data, atau membaca daftar alert. REST API adalah pendekatan yang tepat dan sederhana untuk operasi-operasi seperti ini, sehingga dipisahkan dari WebSocket yang dikhususkan untuk pembaruan *real-time*.

### b. Format JSON yang ringan dan mudah diproses

REST API pada GAMON menggunakan format JSON untuk pertukaran data. JSON bersifat ringan, mudah dibaca manusia, dan mudah diproses oleh kedua sisi (Go dan JavaScript/TypeScript). Setiap respons memiliki struktur yang konsisten (`success` dan `data` atau `message`), sehingga memudahkan *frontend* memproses hasilnya.

### c. Standar dan mudah dipahami

REST API merupakan pendekatan standar industri dengan konsep sumber daya (resource) dan metode HTTP (GET, POST, PUT, DELETE) yang sudah dikenal luas. Ini membuat desain API GAMON konsisten, mudah didokumentasikan, dan mudah diuji.

### d. Mendukung arsitektur client-server yang terpisah

GAMON memisahkan *frontend* dan *backend* sebagai dua bagian yang berjalan terpisah. REST API menjadi kontrak komunikasi yang jelas antara keduanya, sehingga *frontend* tidak perlu mengetahui detail implementasi *backend*, dan sebaliknya.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Cara kerja REST API dalam GAMON dapat diuraikan sebagai berikut.

### a. Sisi Backend: pendaftaran dan penanganan rute

Pada saat server dijalankan, `main.go` mendaftarkan seluruh rute API pada `http.NewServeMux()`. Setiap rute dihubungkan ke sebuah fungsi *handler* yang berada pada folder `handler/`. Contohnya, permintaan ke `/api/devices` ditangani oleh `deviceHandler.HandleDevices`, sedangkan `/api/devices/{id}` ditangani oleh `deviceHandler.HandleDevice`.

Setiap *handler* memeriksa metode HTTP (GET, POST, PUT, DELETE) yang digunakan. Apabila metode tidak sesuai, server mengembalikan respons kesalahan *Method Not Allowed*. Untuk rute dengan parameter dinamis seperti `/api/devices/`, *handler* menguraikan *id* dari bagian akhir URL untuk menentukan perangkat mana yang diproses.

### b. Pemrosesan permintaan

Ketika *frontend* mengirim permintaan, *handler* Go melakukan beberapa langkah berikut secara umum:
1. Menerima dan memvalidasi data yang dikirim (untuk operasi seperti POST dan PUT).
2. Melakukan logika bisnis, misalnya berinteraksi dengan monitoring *engine* atau melakukan operasi pada database.
3. Menyusun respons dalam bentuk JSON melalui fungsi bantuan (`respondData`, `respondSuccess`, `respondError`) pada `handler/helpers.go`.

Struktur respons yang dikembalikan selalu konsisten, yaitu `{ "success": true, "data": ... }` untuk keberhasilan dengan data, atau `{ "success": false, "message": "..." }` untuk kegagalan.

### c. Sisi Frontend: pemanggilan API

Di sisi *frontend*, seluruh pemanggilan API dipusatkan pada `frontend/src/lib/api.ts`. File ini mendefinisikan fungsi-fungsi seperti `fetchDevices()`, `createDevice()`, `updateDevice()`, `deleteDevice()`, `fetchAlerts()`, `fetchDashboard()`, dan lain-lain. Fungsi-fungsi ini menggunakan API bawaan peramban (fetch) untuk mengirim permintaan HTTP ke alamat `http://localhost:8080` (dapat diubah melalui variabel lingkungan `VITE_API_BASE_URL`).

Alur pemanggilan dari sisi frontend:
1. Komponen React memanggil fungsi pada `lib/api.ts`.
2. Fungsi mengirim permintaan HTTP (GET, POST, PUT, DELETE) ke endpoint yang sesuai.
3. Server memproses dan mengembalikan respons JSON.
4. Fungsi memeriksa hasilnya; jika `success`, data dikembalikan untuk diolah komponen; jika gagal, *error* diangkat untuk ditampilkan ke pengguna.

### d. Integrasi dengan WebSocket

REST API dan WebSocket memiliki pembagian peran yang saling melengkapi:
- REST API menangani operasi data yang bersifat *request-response* (pengelolaan, query, konfigurasi).
- WebSocket menangani pembaruan *real-time* (status perangkat yang berubah tanpa diminta *frontend*).

Contohnya, saat pengguna menambah perangkat baru, *frontend* memakai REST API (`POST /api/devices`). Namun setelah itu, perubahan status perangkat tersebut ditampilkan melalui WebSocket. Kedua jalur ini bekerja berdampingan dalam satu aplikasi.

---

## 4. Bukti Penggunaan

1. **File `main.go`** — Pendaftaran seluruh rute API pada `http.NewServeMux()`: `/api/devices`, `/api/alerts`, `/api/dashboard`, `/api/monitoring`, `/api/telegram/*`, `/api/settings`, `/api/health`, serta pengaturan CORS.

2. **Folder `handler/`** — Implementasi *handler* REST API pada `device.go`, `alert.go`, `dashboard.go`, `monitoring.go`, `settings.go`, `telegram.go`, dan `api.go`.

3. **File `handler/helpers.go`** — Fungsi bantuan penyusunan respons JSON (`respondData`, `respondSuccess`, `respondError`).

4. **File `frontend/src/lib/api.ts`** — Lapisan *API client* berisi fungsi pemanggil seluruh endpoint REST dari sisi frontend.

5. **File `Dokumen Gamon (PKL)/Rest API & WebSockets.md`** — Dokumentasi perancangan REST API dan WebSocket pada project GAMON.

6. **File `Dokumen Gamon (PKL)/PRD Integrasi MVP.md`** — Spesifikasi endpoint API dan pembagian peran REST API dengan WebSocket.

---

## Catatan

- REST API pada GAMON diimplementasikan menggunakan *standard library* Go (`net/http`) tanpa *framework* web eksternal.
- Pembagian peran REST API (request-response) dan WebSocket (real-time) dijelaskan lebih lanjut pada bagian WebSocket.