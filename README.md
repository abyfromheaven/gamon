<div align="center">

<h1>GARDA MONITORING (GAMON)</h1>
<p>Aplikasi Pemantau Jaringan & Perangkat Berbasis Web</p>

![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)
![Version](https://img.shields.io/badge/Version-0.4-green?style=flat-square)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![React](https://img.shields.io/badge/React-61DAFB?style=for-the-badge&logo=react&logoColor=black)
![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)

</div>

---

**GAMON (Garda Monitoring)** adalah aplikasi pemantau jaringan dan perangkat berbasis web yang dibangun dengan Go (backend) dan React + TypeScript (frontend). Aplikasi ini memungkinkan administrator jaringan untuk memantau status, latensi, dan ketersediaan perangkat jaringan secara real-time melalui dashboard web, lengkap dengan notifikasi Telegram otomatis saat terjadi gangguan.

## Fitur

- **Pemantauan Real-Time** — Status perangkat diperbarui secara langsung melalui WebSocket tanpa perlu refresh halaman
- **ICMP Ping Monitoring** — Mendeteksi ketersediaan dan latensi perangkat menggunakan ICMP ping
- **Dashboard Ringkas** — Ringkasan jumlah perangkat online, offline, dan warning dalam satu tampilan
- **Manajemen Perangkat** — Tambah, edit, dan hapus perangkat jaringan yang ingin dipantau
- **Sistem Alert** — Peringatan otomatis saat perangkat bermasalah (offline)
- **Notifikasi Telegram** — Kirim peringatan langsung ke Telegram melalui bot, termasuk pairing otomatis
- **Riwayat Pengecekan** — Simpan dan tampilkan riwayat latensi serta status perangkat dari waktu ke waktu
- **Konfigurasi Fleksibel** — Atur interval pengecekan, timeout, dan threshold per perangkat

## Instalasi

### Prasyarat

- [Go](https://go.dev/dl/) 1.21 atau lebih baru
- [Node.js](https://nodejs.org/) 18+ dan npm/yarn
- SQLite (otomatis ter-install sebagai dependency Go)

### Backend

```bash
# Clone repository
git clone https://github.com/abyfromheaven/gamon.git
cd gamon

# Install dependency Go
go mod download

# Jalankan server backend
go run main.go
```

Server backend akan berjalan di `http://localhost:8080`.

### Frontend

```bash
# Pindah ke direktori frontend
cd frontend

# Install dependency
npm install

# Jalankan development server
npm run dev
```

Frontend akan berjalan di `http://localhost:5173` (Vite default).

## Penggunaan

1. Buka browser dan akses `http://localhost:5173`
2. Tambahkan perangkat jaringan baru melalui menu **Devices**
3. Pantau status perangkat secara real-time di **Dashboard**
4. (Opsional) Aktifkan notifikasi Telegram dengan mengatur environment variable `TELEGRAM_BOT_TOKEN`

## Konfigurasi

| Nama | Tipe | Default | Deskripsi |
|------|------|---------|-----------|
| `TELEGRAM_BOT_TOKEN` | `string` | _(kosong)_ | Token bot Telegram untuk mengirim notifikasi. Jika tidak diset, notifikasi Telegram nonaktif |
| Server Port | `int` | `8080` | Port yang digunakan oleh backend server Go |
| Frontend Port | `int` | `5173` | Port development server Vite (frontend) |
| `check_interval` | `int` | `3` | Interval pengecekan default per perangkat (dalam detik) |

## API Reference

### Devices

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/devices` | `GET` | Ambil semua perangkat |
| `/api/devices` | `POST` | Tambah perangkat baru |
| `/api/devices/{id}` | `GET` | Ambil detail perangkat |
| `/api/devices/{id}` | `PUT` | Update perangkat |
| `/api/devices/{id}` | `DELETE` | Hapus perangkat |

### Monitoring

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/monitoring` | `GET` | Ambil status monitoring semua perangkat |
| `/api/monitoring/{id}` | `GET` | Ambil status monitoring perangkat tertentu |

### Alerts

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/alerts` | `GET` | Ambil semua alert |
| `/api/alerts/{id}` | `GET` | Ambil detail alert |
| `/api/alerts/{id}` | `PUT` | Update status alert (resolve/acknowledge) |

### Dashboard

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/dashboard` | `GET` | Ambil ringkasan dashboard (total device, online, offline, warning) |

### Telegram

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/telegram/pair` | `POST` | Pairing bot Telegram dengan chat |
| `/api/telegram/status` | `GET` | Cek status koneksi Telegram |
| `/api/telegram/disconnect` | `POST` | Putuskan koneksi Telegram |

### WebSocket

| Endpoint | Deskripsi |
|----------|-----------|
| `/ws` | Koneksi WebSocket untuk menerima update status perangkat secara real-time |

### Settings & Health

| Endpoint | Method | Deskripsi |
|----------|--------|-----------|
| `/api/settings` | `GET/PUT` | Ambil atau update pengaturan aplikasi |
| `/api/health` | `GET` | Health check endpoint |

## Roadmap

- [ ] Grafik latensi perangkat dalam periode waktu tertentu
- [ ] Export data riwayat pengecekan ke CSV/PDF
- [ ] Autentikasi user (login/register)
- [ ] Multi-user dengan role-based access control
- [ ] Email notification selain Telegram
- [ ] Import perangkat dari file CSV
- [ ] Dark mode untuk dashboard
- [ ] Responsive design untuk mobile

## Acknowledgments

- [Go](https://go.dev/) — Bahasa pemrograman backend
- [Gorilla WebSocket](https://github.com/gorilla/websocket) — Library WebSocket untuk Go
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — Driver SQLite pure-Go
- [React](https://react.dev/) — Library frontend
- [Vite](https://vitejs.dev/) — Build tool frontend
- [Tailwind CSS](https://tailwindcss.com/) — CSS framework
- [Recharts](https://recharts.org/) — Library charting untuk React

## License

Project ini dilisensikan di bawah [MIT License](LICENSE).
