# 🎓 PANDUAN LENGKAP PENJELASAN SIDANG PKL - PROJECT GAMON (Garda Monitoring)

Dokumen ini disusun secara **sederhana, jelas, dan non-teknis** agar mudah dipahami dan dijelaskan saat **Sidang PKL** di hadapan dosen/penguji.

---

## 1. Ringkasan Singkat Project (Elevator Pitch)

> **"GAMON (Garda Monitoring)** adalah sistem pemonitor jaringan dan perangkat berbasis web yang dibuat untuk membantu admin IT memantau status perangkat jaringan (seperti Router, Server, Switch, Access Point, dan Website) secara *real-time* dan otomatis."

### Masalah di Tempat PKL yang Diselesaikan:
- **Sebelumnya**: Pengecekan jaringan dilakukan secara manual satu per satu atau baru ketahuan jika ada gangguan setelah ada pengguna yang melapor.
- **Solusi GAMON**: Pengecekan dilakukan otomatis setiap beberapa detik via metode **ICMP Ping**. Jika ada perangkat mati/down, sistem langsung mendeteksi, menampilkan di web dashboard secara *live*, dan mengirim **notifikasi bahaya ke Telegram HP Admin**.

---

## 2. Arsitektur & Alur Kerja Sistem (Data Flow)

Sistem GAMON terdiri dari 4 lapisan utama yang saling terhubung:

```mermaid
flowchart LR
    subgraph 1. Pemantauan Latensi
        Pinger["ICMP Ping (monitor/ping.go)"]
    end

    subgraph 2. Pengolah Utama (Backend Go)
        Engine["Engine (monitor/engine.go)"]
        DB[("Database SQLite (gamon.db)")]
    end

    subgraph 3. Penyebar Pesan Real-time
        WS["WebSocket Hub (handler/websocket.go)"]
        Bot["Telegram Bot (notification/telegram.go)"]
    end

    subgraph 4. Antarmuka Pengguna
        Web["Web Frontend (React + Vite)"]
        HP["HP Admin (Telegram)"]
    end

    Pinger -->|"Status & Latensi (ms)"| Engine
    Engine -->|"Simpan Log History"| DB
    Engine -->|"Broadcast Data Live"| WS
    Engine -->|"Picu Alert Jika Offline 3x"| Bot
    WS -->|"Update Otomatis (Tanpa Refresh)"| Web
    Bot -->|"Pesan Peringatan"| HP
```

### Langkah demi Langkah Alur Kerja:
1. **Pengecekan (Pinger)**: Server Go mengirim paket ICMP Ping ke IP perangkat setiap $X$ detik (contoh: 3 detik).
2. **Evaluasi Status (Engine)**: Engine menghitung latensi respon ($ms$). Jika perangkat tidak merespon selama 3x pengecekan berturut-turut, status diubah menjadi **OFFLINE**.
3. **Penyimpanan (Database)**: Hasil pengecekan dicatat ke tabel `ping_history` SQLite untuk grafik dan riwayat.
4. **Web Broadcast (WebSocket)**: Perubahan status dan nilai latensi langsung dipancar ke browser. Grafik dan indikator web berubah secara *live* tanpa pengguna menekan tombol refresh.
5. **Notifikasi HP (Telegram Bot)**: Jika terjadi kondisi offline, bot Telegram otomatis mengirim notifikasi bahaya (ALERT) ke HP admin. Saat perangkat menyala kembali, bot mengirim notifikasi pemulihan (RECOVERY).

---

## 3. Struktur Berkas & Kode Utama (File Structure)

Saat penguji meminta membuka kode program, buka file-file berikut berdasarkan fungsinya:

| Direktori / File | Fungsi & Peran Utama | Penjelasan Sederhana untuk Penguji |
| :--- | :--- | :--- |
| `main.go` | **Titik Masuk Utama** | Mengoperasikan database, menjalankan mesin pemantau, dan membuka server web port 8080. |
| `database/db.go` | **Pengelola Database** | Membuat tabel SQLite (`devices`, `ping_history`, `alerts`, `telegram_pairing`, `settings`). |
| `database/models.go` | **Struktur Data** | Mendefinisikan bentuk data perangkat, riwayat ping, dan peringatan (alert). |
| `monitor/ping.go` | **Eksekutor ICMP Ping** | Menjalankan perintah `ping` OS untuk mengukur kecepatan respon IP perangkat dalam milidetik ($ms$). |
| `monitor/engine.go` | **Mesin Utama Pemantau** | Mengatur interval waktu pemantauan, menghitung kegagalan, dan memicu notifikasi. |
| `handler/websocket.go` | **Koneksi Live Web** | Menyebarkan data latensi dan status ke tampilan browser secara *real-time*. |
| `notification/telegram.go` | **Notifikasi HP** | Mengirimkan pesan peringatan otomatis ke aplikasi Telegram Admin. |
| `frontend/src/` | **Tampilan Web (UI)** | Dibuat dengan React & Tailwind CSS untuk menampilkan Dashboard, Monitoring, dan Alert Center. |

---

## 4. Pertanyaan yang Sering Ditanyakan Penguji (FAQ Sidang PKL)

### ❓ Q1: "Kenapa memilih Golang untuk Backend?"
> **Jawaban**: Golang sangat cepat, ringan, dan memiliki fitur **Goroutine** (multithreading bawaan). Dengan Goroutine, backend GAMON bisa memantau puluhan perangkat jaringan secara bersamaan (konkuren) tanpa membebani memori server.

### ❓ Q2: "Kenapa menggunakan SQLite sebagai database?"
> **Jawaban**: SQLite bersifat *embedded* (tersimpan dalam 1 berkas `gamon.db` tanpa perlu install server database terpisah seperti MySQL). Sangat pas untuk aplikasi pemantauan jaringan internal yang membutuhkan kemudahan instalasi dan pemeliharaan.

### ❓ Q3: "Bagaimana cara kerja metode ICMP Ping di aplikasi ini?"
> **Jawaban**: Aplikasi memanfaatkan perintah native `ping` OS melalui paket `os/exec` Go. Aplikasi mengirim 1 paket data ke IP target, lalu menghitung berapa milidetik ($ms$) waktu tempuh balik paket tersebut. Jika paket tidak balik, berarti perangkat terputus/mati.

### ❓ Q4: "Bagaimana cara membuat halaman web ter-update otomatis tanpa me-refresh?"
> **Jawaban**: Menggunakan teknologi **WebSocket**. WebSocket membuka jalur komunikasi dua arah yang terus terbuka antara server Go dan browser React. Saat ada data ping baru, server langsung "mendorong" (push) data tersebut ke browser.

### ❓ Q5: "Bagaimana Telegram Bot bisa tersambung dengan sistem ini?"
> **Jawaban**: Pengguna melakukan *pairing* dengan mengetik kode token di Telegram Bot (`/pair GMN-XXXX`). Chat ID penggunanya akan disimpan di database. Ketika engine menemukan perangkat mati, aplikasi mengirim pesan HTTP POST ke **Telegram Bot API** resmi yang otomatis diteruskan ke HP admin.

---

## 5. Tips Sukses Presentasi Sidang PKL

1. **Gunakan Istilah Sederhana**: Hindari istilah yang terlalu rumit. Gunakan padanan kata seperti "Status Merespon" alih-alih istilah teknis yang berbelit-belit.
2. **Demonstrasikan Fitur**: Tunjukkan bagaimana saat IP perangkat diping dan bagaimana grafik latensi di dashboard bergerak secara *real-time*.
3. **Tunjukkan Notifikasi Telegram**: Perlihatkan pesan Telegram yang masuk saat perangkat di-simulasikan mati.
