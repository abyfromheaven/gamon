# Teknologi: WebSocket (Gorilla WebSocket v1.5.3)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| WebSocket (Gorilla WebSocket) | Komunikasi data real-time antara backend dan frontend: status perangkat, hasil ping, dan perubahan status |

---

## 1. Peran Teknologi dalam Project GAMON

WebSocket merupakan protokol komunikasi yang memungkinkan pertukaran data secara terus-menerus antara dua sisi dalam satu koneksi. Berbeda dengan model permintaan-jawaban pada REST API, WebSocket membuka sebuah koneksi yang tetap terbuka sehingga server dapat mengirim data kepada client kapan saja tanpa harus menunggu permintaan terlebih dahulu.

Dalam GAMON, WebSocket dipakai untuk menyampaikan pembaruan **real-time** yang menjadikan aplikasi dapat menampilkan kondisi perangkat yang berubah tanpa diminta. Secara khusus, aplikasi menggunakan pustaka **Gorilla WebSocket** pada sisi backend Go untuk menangani seluruh koneksi ini.

Peran WebSocket dalam GAMON antara lain:

| Peran | Penjelasan |
|-------|-----------|
| Mengirim status awal (initial state) | Saat peramban baru membuka aplikasi, server langsung mengirim status terkini seluruh perangkat tanpa harus menunggu ping pertama |
| Mendorong hasil pemeriksaan (check result) | Setiap kali perangkat diping, hasilnya dikirim langsung ke peramban secara real-time |
| Memberi tahu perubahan status | Ketika sebuah perangkat berubah (misal online → offline), peramban segera menerima informasi perubahannya |
| Menghemat pengguna sumber daya | Tidak perlu peramban terus-menerus meminta data (polling) sebab server yang mengirim secara otomatis |

**Bukti penggunaan:**
- `go.mod` — dependensi `github.com/gorilla/websocket` v1.5.3
- `handler/websocket.go` — implementasi Hub WebSocket dan koneksi client menggunakan Gorilla WebSocket
- `frontend/src/hooks/useWebSocket.ts` — koneksi WebSocket di sisi frontend

---

## 2. Alasan Pemilihan Teknologi (Gorilla WebSocket)

Alasan pemilihan WebSocket dengan pustaka **Gorilla WebSocket** dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

**a. Mendukung pembaruan real-time tanpa polling.** GAMON menuntut tampilan yang responsif dan selalu terbaru. Dengan WebSocket, server langsung mendorong data ketika ada perubahan, sehingga frontend tidak perlu terus menerus meminta data (yang lebih boros sumber daya).

**b. Pustaka Go yang mapan dan mudah diintegrasikan.** Gorilla WebSocket merupakan pustaka pengelolaan koneksi WebSocket pada bahasa Go yang telah dikenal dan banyak dipakai. Kode ini mudah diintegrasikan dengan backend GAMON dan menangani proses peningkatan (upgrade) koneksi HTTP ke WebSocket.

**c. Mendukung komunikasi real-time yang dibutuhkan sistem monitoring.** Karena GAMON merupakan aplikasi monitoring, kebutuhan untuk menyampaikan kondisi perangkat yang terus berubah menjadi sangat penting. WebSocket memenuhi kebutuhan ini dengan jalur komunikasi yang berkelanjutan.

**d. Hemat penggunaan sumber daya dibanding polling.** Dengan menjaga satu koneksi tetap terbuka, WebSocket tidak menghabiskan sumber daya untuk membuka banyak permintaan seperti pada metode polling, yang selaras dengan tujuan aplikasi monitoring.

---

## 3. Bagaimana Teknologi Bekerja dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja WebSocket menggunakan Gorilla WebSocket dapat diuraikan sebagai berikut.

**a. Inisialisasi Hub.** Saat aplikasi dimulai, backend membuat sebuah objek *Hub* yang bertugas mengelola seluruh koneksi client. Hub memiliki daftar client yang terhubung, disertai sejumlah *channel* yang digunakan untuk mendaftarkan, menghapus, dan menyiarkan pesan ke semua client. Komponen ini menjadi pusat pengatur lalu lintas koneksi WebSocket dalam aplikasi.

**b. Pembukaan koneksi.** Ketika frontend membuka koneksi ke server (di `/ws`), backend menerima permintaan HTTP dan melakukan *upgrade* menjadi koneksi WebSocket melalui *upgrader* dari Gorilla WebSocket. Setelah berhasil, sebuah `Client` baru dibuat dan didaftarkan ke *Hub*. Koneksi tersebut kemudian dijalankan dengan dua proses terpisah (*pump*): satu untuk menerima pesan masuk (*readPump*) dan satu untuk mengirim pesan keluar (*writePump*).

**c. Pengiriman status awal (initial state).** Begitu sebuah client terdaftar, *Hub* langsung mengirimkan **initial state**, yaitu status terkini dari semua perangkat yang aktif. Data ini diambil dari database dengan menggabungkan data perangkat dan hasil pemeriksaan terakhir. Tujuannya adalah agar setiap kali pengguna membuka dashboard, halaman tidak kosong, tetapi langsung menampilkan kondisi terkini.

**d. Penyiaran perubahan (broadcast).** Ketika mesin monitoring (engine) menyelesaikan sebuah pemeriksaan atau mendeteksi perubahan status, server mengemas hasil tersebut menjadi sebuah pesan, lalu menyiarkannya ke seluruh client yang terhubung melalui *Hub*. Dengan cara ini, semua peramban yang terbuka menerima hasil yang sama secara real-time.

**e. Penerimaan di sisi peramban.** Di sisi frontend, koneksi WebSocket dibuat dengan API bawaan peramban. Hook yang digunakan memperbarui data berdasarkan pesan yang diterima, sehingga dashboard dan halaman monitoring selalu menampilkan nilai terkini. Apabila koneksi terputus, frontend secara otomatis mencoba terhubung kembali dalam beberapa detik agar pembaruan tetap berjalan.

Dengan alur tersebut, WebSocket menjadi tulang punggung penyajian data **real-time** GAMON, sementara REST API menangani operasi data yang bersifat permintaan-jawaban. Keduanya saling melengkapi dalam komunikasi antara frontend dan backend.

---

## 4. Bukti Penggunaan

1. **File `go.mod`** — mencantumkan dependensi `github.com/gorilla/websocket` v1.5.3 sebagai pustaka WebSocket pada sisi server.

2. **File `handler/websocket.go`** — berisi implementasi `Hub`, `Client`, serta proses `readPump` dan `writePump` yang menggunakan Gorilla WebSocket, termasuk fungsi penyiaran (`Broadcast`).

3. **File `main.go`** — membuat `handler.NewHub` dan mendaftarkan rute `/ws` untuk menangani koneksi WebSocket, serta menjalankan `go hub.Run()`.

4. **File `monitor/engine.go`** — menggunakan `Hub.Broadcast` untuk mengirim hasil pemeriksaan dan perubahan status.

5. **File `frontend/src/hooks/useWebSocket.ts`** — pembuatan koneksi WebSocket di sisi frontend, penanganan `initial_state`, `check_result`, `status_change`, serta koneksi ulang otomatis.

---

## Catatan

- Pada sisi server, WebSocket dikelola oleh pustaka **Gorilla WebSocket**; pada sisi frontend, memanfaatkan API WebSocket bawaan peramban tanpa pustaka tambahan.
- WebSocket dan REST API berbagi peran: REST menangani data permintaan-jawaban, sedangkan WebSocket menghadirkan pembaruan real-time yang berkesinambungan.