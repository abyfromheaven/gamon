# Teknologi: Golang (Go 1.26.5)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| Go (Golang) | Backend API, proses monitoring, deteksi & resolusi alert, komunikasi WebSocket, integrasi database, koneksi Telegram |

---

## 1. Peran Teknologi dalam Project GAMON

Golang merupakan satu-satunya bahasa pemrograman yang membangun seluruh sisi *backend* atau bagian server dari aplikasi GAMON (Garda Monitoring). Berbeda dengan *frontend* yang berjalan di peramban pengguna, seluruh pemrosesan, penyimpanan, dan komunikasi data di GAMON dikerjakan oleh program Go yang berjalan sebagai sebuah layanan server pada port 8080.

Secara umum, Go di dalam project ini tidak hanya berperan sebagai pembuat antarmuka API (Application Programming Interface), melainkan mengemban seluruh tanggung jawab teknis yang menjadi "jantung" sistem. Peran-peran tersebut dapat dirinci sebagai berikut.

### a. Layanan REST API (Application Programming Interface)

Go menyediakan seluruh *endpoint* API yang dikonsumsi oleh *frontend*. Melalui protokol REST dan format JSON, *frontend* React melakukan operasi CRUD (Create, Read, Update, Delete) terhadap data perangkat, membaca data ringkasan dashboard, melihat riwayat alert, mengubah pengaturan, hingga mengontrol proses monitoring (mulai atau berhenti). Seluruh *endpoint* ini didaftarkan dan dijalankan menggunakan paket bawaan Go, `net/http`.

Bukti: pendaftaran seluruh rute `/api/*` hingga aplikasi dinyatakan pada `main.go`.

### b. Mesin Monitoring Perangkat (Monitoring Engine)

Go menjalankan proses pemantauan perangkat secara otomatis dan berkelanjutan. Setiap perangkat aktif yang terdaftar akan disimak oleh sebuah proses independen yang menjalankan siklus pengecekan (check loop) secara berulang sesuai interval yang ditentukan. Go mengeksekusi perintah `ping` sistem untuk tiap alamat IP, menguraikan (parse) hasil keluarannya untuk memperoleh status, latensi, dan informasi lain, lalu mencatat hasil tersebut.

Mesin ini dijelaskan secara rinci pada dokumen `Dokumen Gamon (PKL)/Sistem Gamon/Engine Enhancement.md` yang menyebutnya sebagai "jantung" sistem yang terus berdetak agar setiap perangkat selalu dalam kondisi terpantau.

### c. Manajemen Concurrency dengan Goroutine

Salah satu peran paling penting Go adalah mengelola banyak pengecekan perangkat secara bersamaan. Go menggunakan mekanisme goroutine yang merupakan eksekusi ringan dan hemat memori dibandingkan *thread* biasa. Dengan *goroutine*, GAMON mampu memeriksa banyak alamat IP sekaligus tanpa membuat server menjadi lambat atau membeku.

Pada implementasi, satu *goroutine* dibuka untuk setiap perangkat yang dimonitor. Proses tersebut dapat dihentikan ketika diperlukan (misalnya saat perangkat dinonaktifkan) menggunakan *context.Context*, dan guard pengaman dengan `sync.Mutex` agar akses data bersama antar-*goroutine* tetap aman.

### e. Komunikasi Real-time (WebSocket)

Go bertanggung jawab mengoperasikan sebuah Hub WebSocket yang menghubungkan backend dengan semua client browser yang sedang terhubung. Saat sebuah perangkat berubah status (misalnya dari online menjadi offline), Go dengan cepat mendorong (broadcast) data terbaru ke seluruh client melalui WebSocket. Dengan cara ini dashboard dapat diperbarui secara instan tanpa perlu halaman di-refresh oleh pengguna.

### f. Deteksi & Pengelolaan Alert

Go secara otomatis membuat (generate) alert ketika mendeteksi perubahan status perangkat yang mengindikasikan gangguan, dan secara otomatis menandai alert selesai (resolve) ketika perangkat kembali pulih. Alert tersebut disimpan ke dalam database dan disiarkan ke seluruh client melalui WebSocket.

### g. Koneksi Telegram (Notification)

Go juga bertugas menghubungkan dengan Bot Telegram untuk mengirimkan notifikasi. Ketika sebuah perangkat bermasalah, Go mengirimkan pesan alert; ketika perangkat pulih, Go mengirim pesan *recovery*. Selain itu, Go menjalankan *poller* untuk membaca perintah yang dikirim pengguna ke bot Telegram, seperti `/pair`, `/status`, dan `/help`.

### h. Akses dan Pengelolaan Database

Go bertanggung jawab membuka koneksi ke database SQLite, menjalankan *migration* (pembaruan struktur tabel), serta melakukan operasi baca dan tulis data perangkat, hasil pemantauan, alert, dan pengaturan aplikasi.

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan Go dapat disimpulkan dari dokumentasi rancangan dan bukti implementasi pada project GAMON, sebagai berikut.

### a. Concurrency yang ringan (Goroutines)

Alasan paling utama. Aplikasi monitoring harus melakukan pengecekan ke puluhan hingga ratusan alamat IP perangkat secara bersamaan dan *real-time*. Dalam bahasa lain, membuat banyak proses sekaligus dapat memakan memori besar dan membuat aplikasi melambat. Go dengan *goroutine* memungkinkan pemerik berbagai IP secara paralel dengan konsumsi memori yang jauh lebih hemat.

**Bukti implementasi:** pemakaian *goroutine* (sintaks `go ...`) dan `context.Context` pada `main.go` (misalnya fungsi `autoStartMonitoring`) serta `monitor/engine.go` yang membuktikan pemanfaatan konkuren untuk memonitor banyak perangkat.

### b. kecepatan eksekusi tinggi

Go termasuk bahasa yang **dikompilasi hansung ke bahasa mesin** (*compiled language*). Karena hasil kompilasinya adalah kode mesin langsung, eksekusi Go lebih cepat dibandingkan pendekatan berbentuk skrip atau *runtime* lain (seperti yang umum dalam bahasa Python atau Node.js) dalam memproses data jaringan mentah.

### c. Hasil berupa single binary

Produk akhir dari backend Go berupa satu file *executable* tunggal (`.exe` atau `.bin`). Di lingkungan yang membutuhkan kemudahan *deploy*, hal ini sangat disukai karena aplikasi cukup dijalankan tanpa perlu menginstal *runtime* per bahasa tambahan di server.

**Bukti pada implementasi:** pada `reset.sh`, script menangani file binary `gamon` dan memulainya kembali dengan perintah `go run .`.

### d. Mendukung WebSocket real-time

Go dipilih bersama library `gorilla/websocket` untuk memungkinkan backend mendorong data perubahan status ke frontend secara instan, tanpa harus memakai metode polling yang lebih boros sumber daya. Ini kesesuaian dengan kebutuhan GAMON yang menekankan pembaruan tampilan tanpa proses manual `refresh`.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Ketika aplikasi GAMON dijalankan perintah `go run .`, fungsi `main()` pada `main.go` dieksekusi sebagai titik masuk. Prosesnya dapat diuraikan menjadi beberapa tahap berikut.

1. **Inisialisasi Database.** Aplikasi menghubungkan Go dengan database SQLite melalui `database.NewDB()`. Database ini menyimpan data perangkat, data hasil, alert, state telegram, dan pengaturan. Apabila ada perubahan struktur tabel, Go menjalankan *migration* otomatis.

2. **Pembuatan Hub WebSocket.** Aplikasi membuat sebuah `Hub` yang bertugas mengelola semua koneksi WebSocket dari *frontend*. Hub dijalankan sebagai *goroutine* latar belakang agar mampu menangani banyak koneksi sekaligus.

3. **Inisialisasi Notifikasi Telegram.** Aplikasi membuat *notifier* dan *poller* Telegram. Apabila token bot tersedia (variabel lingkungan `TELEGRAM_BOT_TOKEN`), *poller* dijalankan sebagai *goroutine* yang memeriksa perintah tele Telegram secara berkala.

4. **Pembuatan Monitoring Engine.** **Engine** monitoring dibentuk dengan menerima Hub WebSocket, database, dan notifier. Engine inilah yang nanti menggerakkan proses pemantauan perangkat.

5. **Pendaftaran Route API.** Go mendaftarkan seluruh *section* endpoint REST pada sebuah mux (`http.NewServeMux`) untuk `/api/devices`, `/api/alerts`, `/api/dashboard`, `/api/monitoring`, `/api/telegram/*`, `/api/settings`, dan `/api/health`. Semua dikelola oleh *handler handler* yang ditulis dalam paket Go.

6. **Auto-Start Monitoring.** Sebuah *goroutine* menjalankan `autoStartMonitoring`. Setelah menunggu dua detik, aplikasi membaca seluruh perangkat berstatus `active` dari database dan memulai satu *goroutine* monitoring untuk masing-masing perangkat tersebut. Dengan begitu, saat server dinyalakan, semua perangkat aktif otomatis terpantau tanpa perlu di-*start* manual.

7. **Server HTTP Dijalankan.** Akhirnya server HTTP dijalankan pada `:8080` dengan `http.ListenAndServe`. Pada tahap inilah aplikasi siap menerima permintaan dari *frontend* dan meneruskan proses monitoring.

Setelah server berjalan, setiap *goroutine* monitoring menjalankan siklus pengecekan per device. Satu siklus terdiri dari eksekusi perintah ping, penguraian hasil, pencatatan hasil ke database, pengecekan perubahan status, pembuatan / resolusi alert jika diperlukan, siaran real-time ke WebSocket, dan pengiriman notifikasi ke Telegram. Siklus ini berulang pada interval yang telah ditentukan.

---

## 4. Bukti Penggunaan

1. **File `go.mod`** — Deklarasi module bernama `gamon` dengan versi Go `1.26.5`, serta dependensi `gorilla/websocket` dan `modernc.org/sqlite`. Bukti bahwa project merupakan project Go.

2. **File `main.go`** — Berisi `package main` dan `func main()`; menginisiasi database, hub WebSocket, notifier, polling, memenregistrasi seluruh route `/api/*`, memetakkan CORS dahulu, menjalankan *auto-start monitoring*, dan akhirnya memuat server `http.ListenAndServe(":8080", ...)`.

3. **Folder backend** — `database/`, `handler/`, `monitor/`, `notification/` berisi file berkategori `package` dalam bahasa Go yang menjadi lapisan database, API, monitoring, dan notifikasi.

4. **File `monitor/engine.go`** — Manajemen goroutine, status tracking, auto-alert, auto-resolve; diterangkan juga dalam `Dokumen Gamon (PKL)/Sistem Gamon/Engine Enhancement.md`.

5. **File `monitor/ping.go`** — Memanfaatkan `os/exec` dan `runtime.GOOS` untuk mengeksekusi perintah `ping` sistem sesuai platform (Windows/Linux).

6. **File `handler/websocket.go`, `handler/telegram.go`** — Implementasi WebSocket Hub dan integrasi Telegram di sisi server menggunakan Go.

7. **File `reset.sh`** — Skrip yang mengelola binary Go `gamon` dan menjalankan `go run .`, menunjukkan alur *development* dan eksekusi berbasis Go.

---

## Catatan

- Dependensi eksternal langsung yang dipakai kode Go hanya `gorilla/websocket` dan `modernc.org/sqlite` (untuk akses SQLite tanpa CGO). Lainnya adalah ketergantungan transitif otomatis.
- Seluruh *framework* web eksternal tidak digunakan; routing dan server memanfaatkan Go *standard library* (`net/http`).