# Penjelasan Ilmiah: WebSocket pada Sistem GAMON

> Materi untuk dituliskan pada **Bab 2 (Landasan Teori)** dan **Bab 4 (Analisis & Perancangan)** laporan PKL / Karya Tulis Ilmiah. Salin sesuai kebutuhan, lalu sesuaikan nomor bab dan subbab dengan outline laporan.

---

## BAB 2. LANDASAN TEORI

### 2.X Komunikasi Data pada Aplikasi Web

Komunikasi data antara *frontend* dan *backend* pada aplikasi web dapat dilakukan dengan beberapa teknologi. Pada sistem monitoring, dua teknologi yang umum digunakan adalah **REST API** dan **WebSocket**. Keduanya berjalan di atas protokol HTTP dan TCP, namun memiliki karakteristik yang berbeda.

#### 2.X.1 REST API (Representational State Transfer)

REST API merupakan gaya arsitektur komunikasi yang bekerja berdasarkan pola **permintaan-jawaban** (*request-response*). Klien mengirimkan permintaan HTTP kepada server, server memproses dan mengembalikan hasil, lalu koneksi ditutup. Setiap permintaan bersifat independen. REST API cocok untuk operasi yang bersifat *one-shot*, seperti CRUD, membaca data historis, dan pengaturan sistem.

Kelemahan utama REST API adalah **tidak dapat mengirim data secara spontan** dari server ke klien. Server hanya membalas ketika klien meminta. Untuk memperoleh data yang selalu terbaru, klien harus melakukan **polling**, yaitu meminta data secara berulang dalam interval tertentu. Teknik polling memiliki kekurangan:

1. **Latensi** — data yang ditampilkan dapat terlambat hingga satu interval polling penuh dari kondisi sebenarnya;
2. **Overhead** — setiap permintaan HTTP membawa *header* tambahan sehingga lalu lintas jaringan boros bila dilakukan terus menerus;
3. **Beban server** — server membuka, memproses, dan menutup koneksi TCP baru setiap kali permintaan datang.

#### 2.X.2 WebSocket

WebSocket adalah protokol komunikasi jaringan yang menyediakan koneksi **dua arah penuh** (*full-duplex*) dan **persisten** antara klien dan server melalui satu koneksi TCP. Koneksi diawali dengan *handshake* HTTP menggunakan *header* `Upgrade: websocket`. Setelah berhasil, kedua pihak dapat mengirim pesan kapan saja tanpa menunggu permintaan.

Keunggulan utama WebSocket adalah kemampuan **server push**, yaitu server mengirim data kepada klien secara spontan tanpa diminta. Karakteristik ini menjadikan WebSocket sangat cocok untuk aplikasi **real-time** seperti *monitoring*, *chat*, dan *dashboard* yang membutuhkan pembaruan terus-menerus.

---

## BAB 4. ANALISIS & PERANCANGAN SISTEM

### 4.X Analisis Kebutuhan Komunikasi Data GAMON

GAMON merupakan aplikasi *monitoring* jaringan yang berjalan **24 jam**. Mesin monitoring (*monitoring engine*) memeriksa setiap perangkat secara berkala sesuai *check interval* dan menghasilkan data berupa:

- Hasil pemeriksaan perangkat (*check result*): status *online/offline* dan *latency*;
- Perubahan status (*status change*): saat perangkat berpindah dari *online* menjadi *offline* atau sebaliknya;
- *Alert* baru ketika perangkat terdeteksi gagal secara berkelanjutan.

Data tersebut bersifat **sangat dinamis** dan **terus-menerus diperbarui**. Apabila disalurkan melalui REST API dengan polling, *frontend* harus meminta data berulang setiap detik agar tampilan selalu terbaru. Hal ini menimbulkan tiga persoalan berikut.

#### 4.X.1 Analisis Masalah

1. **Masalah real-time** — pemantauan 24 jam menuntut *tim NOC* (Network Operation Center) mengetahui kondisi perangkat secepat mungkin. Dengan polling berskala detik, terdapat jeda antara saat perangkat sebenarnya mati dengan saat data tampil di layar. Semakin besar interval polling, semakin lama keterlambatan informasi.

2. **Masalah efisiensi jaringan** — setiap permintaan HTTP membawa *header* berukuran sekitar 1–2 KB. Jika sistem memantau banyak perangkat dan data diminta setiap detik, mayoritas lalu lintas jaringan hanya berisi *header*, bukan data yang bermanfaat.

3. **Masalah beban server** — setiap permintaan REST membuka koneksi TCP baru, memproses permintaan, mengakses database, mengirim respons, lalu menutup koneksi. Pola bongkar-pasang koneksi ini membebani CPU dan RAM server secara signifikan.

#### 4.X.2 Pemecahan Masalah: Hybrid REST API + WebSocket

Berdasarkan analisis di atas, GAMON menerapkan arsitektur komunikasi **hybrid** yang membagi peran kedua teknologi secara jelas.

| Aspek | REST API | WebSocket |
|-------|----------|-----------|
| Pola komunikasi | *Request-response* (satu kali) | *Full-duplex* persisten |
| Koneksi | Dibuka dan ditutup setiap permintaan | Dibuka sekali, tetap terbuka |
| Inisiator pengiriman | Klien | Server (dan klien) |
| Kecocokan | Operasi CRUD, konfigurasi, data historis | Data real-time yang terus berubah |
| Kebutuhan refresh | Ya (polling manual) | Tidak |
| Overhead per pertukaran | ± 1–2 KB (header HTTP) | ± 2 byte (framing) |
| Beban server | Tinggi jika dipaksa polling | Rendah |

#### 4.X.3 Peran REST API pada GAMON

REST API menangani operasi yang bersifat *request-response* dan tidak membutuhkan pembaruan spontan:

- Pengelolaan perangkat: menambah, mengubah, menghapus, membaca daftar, serta mulai/menghentikan monitoring (`/api/devices`);
- Pengelolaan *alert*: membaca, menyelesaikan, dan menandai sebagai diketahui (`/api/alerts`);
- Data dashboard dan riwayat monitoring (`/api/dashboard`, `/api/monitoring`);
- Pengaturan aplikasi (`/api/settings`) dan koneksi Telegram (`/api/telegram/*`);
- Pemeriksaan kesehatan server (`/api/health`).

Operasi-operasi tersebut hanya terjadi saat pengguna berinteraksi, sehingga tidak memerlukan koneksi yang terus terbuka.

#### 4.X.4 Peran WebSocket pada GAMON

WebSocket menangani seluruh data yang harus sampai ke *frontend* secara spontan dan real-time:

- **Initial state** — saat pengguna membuka aplikasi, server langsung mengirim status terkini seluruh perangkat;
- **Check result** — setiap hasil pemeriksaan perangkat dikirim langsung ke semua peramban yang terhubung;
- **Status change** — ketika perangkat berubah status, peramban segera menerima informasi perubahan.

Dengan WebSocket, *frontend* **tidak perlu di-refresh** untuk memperbarui tampilan karena server mendorong data setiap kali ada perubahan.

#### 4.X.5 Implementasi pada GAMON

- **Sisi backend (Go)**: koneksi dikelola oleh komponen *Hub* pada `handler/websocket.go` menggunakan pustaka *Gorilla WebSocket*. Saat ada client terhubung, server mengirim `initial_state`; ketika *monitoring engine* (pada `monitor/engine.go`) menyelesaikan pemeriksaan, hasil disiarkan ke seluruh client melalui `Hub.Broadcast` dengan tipe pesan `check_result` dan `status_change`.
- **Sisi frontend (React)**: koneksi dibuat melalui *hook* `useWebSocket` pada `frontend/src/hooks/useWebSocket.ts`. Pesan yang diterima langsung memperbarui status pada dashboard tanpa reload. Apabila koneksi terputus, *hook* melakukan koneksi ulang otomatis setiap 3 detik.

#### 4.X.6 Alasan Tidak Menggunakan REST API Murni untuk Data Real-time

REST API murni tidak dipilih sebagai satu-satunya jalur komunikasi karena keterbatasan berikut:

1. **Tidak mendukung server push** — REST hanya berjalan dengan pola permintaan-jawaban, sehingga server tidak dapat memberitahu *frontend* saat perangkat mati secara instan;
2. **Harus polling** — untuk meniru real-time dengan REST, diperlukan polling yang boros sumber daya dan tetap menghasilkan keterlambatan informasi;
3. **Latensi deteksi tinggi** — pada sistem monitoring, kecepatan deteksi kegagalan sangat kritis. Jeda sebesar satu interval polling dapat menunda respons *tim NOC* terhadap insiden jaringan.

Oleh karena itu, REST API hanya digunakan untuk operasi data yang jarang dan tidak kritis waktunya, sedangkan data monitoring real-time disalurkan melalui WebSocket.

---

## 5.X Kesimpulan

GAMON menerapkan arsitektur komunikasi **hybrid**: **REST API** digunakan untuk manajemen data yang bersifat permintaan-jawaban (CRUD perangkat, alert, pengaturan, dan koneksi Telegram), sedangkan **WebSocket** digunakan untuk menyalurkan data pemantauan secara **real-time** (status perangkat, hasil pemeriksaan, dan perubahan status) tanpa perlu me-refresh halaman. Pembagian peran ini menjawab kebutuhan aplikasi monitoring yang berjalan 24 jam, sekaligus menjaga efisiensi jaringan dan beban server.

---

## Daftar Pustaka (contoh)

- Fette, I., & Melnikov, A. (2011). *The WebSocket Protocol*. RFC 6455, IETF.
- Fielding, R. T. (2000). *Architectural Styles and the Design of Network-based Software Architectures*. Doctoral dissertation, University of California, Irvine.
- Gorilla WebSocket. *gorilla/websocket* — library WebSocket untuk bahasa Go. https://github.com/gorilla/websocket
