Untuk kebutuhan proyek GAMON (Garda Monitoring), jawaban singkatnya adalah lu wajib menggunakan keduanya secara bersamaan (Hybrid), karena masing-masing memiliki peran yang sangat berbeda. Namun, untuk fitur pemantauan status perangkat secara _realtime_, WebSockets adalah juaranya.

Berikut adalah perbedaan, kecocokan, dan alasan ilmiah yang bisa langsung lu salin ke Bab 2 (Landasan Teori) atau Bab 4 (Analisis & Perancangan) di laporan PKL lu.

---

## 1. Perbedaan Mendasar (Analogi Sederhana)

- REST API (Sistem Ketuk Pintu): Frontend React harus terus-menerus mengetuk pintu backend Golang untuk bertanya, _"Apakah ada server yang down?"_. Jika tidak ada, Golang menjawab _"Tidak ada"_, lalu koneksi ditutup. Proses ini diulang terus setiap beberapa detik (_Polling_).
- WebSockets (Sistem Telepon/HT): Frontend React dan backend Golang membuka satu jalur telepon yang terus tersambung. Frontend tidak perlu bertanya. Begitu ada server NOC yang _down_, Golang akan langsung berteriak lewat telepon tersebut, _"Woi, Server B mati!"_.

---

## 2. Mana yang Cocok untuk GAMON? (Pembagian Tugas)

Dalam aplikasi GAMON, lu harus membagi tugas keduanya seperti ini:

- Gunakan REST API untuk fungsi CRUD (Create, Read, Update, Delete):
    
    - Menambahkan IP perangkat baru ke database.
    - Menghapus atau mengedit nama perangkat.
    - _Login_ dan _logout_ akun anak NOC. [1, 2]
    
- Gunakan WebSockets untuk fitur Monitoring Utama:
    
    - Menampilkan grafik _latency_ (ms) yang bergerak setiap detik.
    - Mengubah warna indikator lampu dari hijau (aman) ke merah (down) secara instan saat perangkat mati.
    

---

## 3. Alasan Ilmiah (Bahan Tulisan Laporan PKL Lu)

Jika penguji sidang bertanya, _"Kenapa kamu pakai WebSockets, kenapa tidak pakai REST API biasa untuk realtime-nya?"_, ini 3 alasan ilmiahnya:

## A. Efisiensi Bandwidth dan Protokol (Overhead Jaringan)

- REST API (HTTP): Setiap kali React meminta data (_HTTP Request_), browser harus mengirimkan data tambahan bernama _HTTP Header_ (berisi informasi browser, cookie, dll) berukuran sekitar 1 KB - 2 KB. Jika ada 50 perangkat dan dicek setiap 1 detik, lalu lintas jaringan akan penuh hanya untuk mengirim _header_ kosong. [3]
- WebSockets (TCP): WebSockets hanya melakukan _HTTP Handshake_ satu kali di awal. Setelah tersambung, pengiriman data status jaringan hanya memakan _overhead_ sebesar 2 Byte saja per pesan. Ini sangat menghemat bandwidth jaringan NOC. [4, 5]

## B. Mengurangi Beban Kerja Server (CPU & RAM)

- REST API: Jika lu memaksa menggunakan REST API secara berulang (_Short Polling_) setiap detik, server Golang harus terus-menerus membuka koneksi TCP baru, memproses _request_, membaca database, mengirim jawaban, lalu menutup koneksi TCP tersebut. Ini akan membuat CPU server melonjak drastis.
- WebSockets: Koneksi TCP dibuat sekali dan dijaga agar tetap terbuka (_persistent connection_). Server Golang tidak perlu membuang energi untuk bongkar-pasang koneksi baru setiap detik. [6, 7]

## C. Menurunkan Waktu Respons (SLA - Service Level Agreement)

- REST API: Ada jeda waktu (_latency gap_) antara waktu perangkat sebenarnya mati dengan waktu aplikasi melakukan _request_ berikutnya. Jika pengecekan diatur tiap 10 detik, ada potensi tim NOC terlambat sadar selama 10 detik.
- WebSockets: Menggunakan prinsip _Event-Driven_. Detik itu juga _Goroutine_ di Golang mendeteksi Ping gagal, detik itu juga data diledakkan ke Frontend React secara instan (realtime sejati) tanpa ada jeda tunggu.

---

## Kesimpulan untuk Laporan

Di dalam kesimpulan bab analisis, lu bisa menulis:

> _"Aplikasi GAMON menerapkan arsitektur Hybrid, menggunakan REST API untuk manajemen data perangkat (CRUD), dan memanfaatkan protokol WebSockets untuk menyalurkan data pemantauan secara realtime guna mencapai efisiensi bandwidth dan responsivitas sistem yang tinggi bagi tim NOC."_

Konsep ilmiahnya sudah sangat matang dan tidak bisa didebat lagi oleh penguji. Gimana? Mau gua buatkan contoh struktur JSON data monitoring yang ideal dikirim via WebSocket dari Go ke React?

---

Perbedaan Peran
┌─────────────────────────────────────────────────┐
│              DUA SALURAN KOMUNIKASI              │
├──────────────────────┬──────────────────────────┤
│    REST API          │      WebSocket           │
│    (Penunjang)       │      (Utama)             │
├──────────────────────┼──────────────────────────┤
│ HTTP Request/Response│ Persistent Connection    │
│ One-shot (tiap kali) │ Always-on (selalu hidup) │
│ Untuk KONTROL         │ Untuk DATA REALTIME      │
└──────────────────────┴──────────────────────────┘
Kapan REST API Dipakai?
Saat awal interaksi — sebelum WebSocket aktif:
1. User buka browser, ketik IP
2. Frontend kirim POST /api/monitor { "ip": "192.168.1.1" }
   → Ini REST API, cuma sekali waktu submit
3. Backend terima, mulai goroutine ping
4. SETELAH ITU, semua data jalan via WebSocket
   → REST API sudah tidak dipakai lagi
Flow-nya Jadi Begini
        REST API (sekali aja)          WebSocket (selanjutnya)
              │                              │
User Input ──▶│                              │
              │── POST /api/monitor ────────▶│
              │                              │
              │   Backend mulai goroutine    │
              │                              │
              │  ◀── WS: ping_result ────────│
              │  ◀── WS: ping_result ────────│
              │  ◀── WS: status_change ──────│
              │        (terus menerus)       │
              ▼                              ▼
         "Oke mulai"              Dashboard update realtime
Kenapa Tetap Pakai REST API?
Fungsi	Kenapa REST, bukan WS
Submit IP pertama kali	Simpel, POST sekali, dapat response sukses/gagal langsung
Stop monitoring	POST /api/stop — clean, tidak perlu kirim WS message
Cek status server	GET /api/health — untuk cek apakah backend hidup
CORS handling	Lebih natural di HTTP daripada WS
Jadi Intinya
- REST API = Tombol ON/OFF — untuk start/stop monitoring
- WebSocket = Layar LCD — untuk terima data status secara realtime terus-menerus
REST API itu seperti remote TV — kamu pencet sekali untuk nyalakan, tapi setelah itu TV nyala terus tanpa kamu pencet lagi.