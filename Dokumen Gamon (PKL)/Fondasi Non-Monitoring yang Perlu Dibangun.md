
1. Database Layer (SQLite)
Ini paling dasar. Tanpa database, semua data hanya di memori dan hilang saat restart.
Yang perlu:
├── Koneksi SQLite + migrations
├── Model: devices, ping_history, settings
└── Write queue (buffer sebelum batch insert)
Kenapa dulu: Semua fitur lain (device management, history, config) bergantung pada ini.
2. Device Management (CRUD)
Sekarang cuma bisa input IP lewat form, tidak ada permanensi. Butuh CRUD:
API yang perlu:
├── GET    /api/devices        → list semua device
├── POST   /api/devices        → tambah device baru
├── PUT    /api/devices/:id    → update device
├── DELETE /api/devices/:id    → hapus device
└── POST   /api/devices/:id/monitor → start/stop monitoring
Yang disimpan:
- IP address, hostname (opsional), deskripsi
- Status monitoring (aktif/nonaktif)
- Group/label (misal: "Router", "Server", "PC")
3. Configuration System
Sekarang semua hardcoded (port 8080, interval 3s). Butuh config:
config.yaml:
├── server:
│   ├── port: 8080
│   └── host: "0.0.0.0"
├── scan:
│   ├── default_interval: 3s
│   ├── ping_timeout: 3s
│   └── max_concurrent: 50
├── database:
│   └── path: "./data/gamon.db"
└── retention:
    ├── history_days: 7
    └── auto_cleanup: true
Kenapa penting: User bisa sesuaikan tanpa rebuild.
4. Simple Auth (Opsional tapi Disarankan)
Sekarang API dan WebSocket terbuka total. Minimal:
┌─────────────────────────────┐
│  Opsi A: API Key            │
│  - Header: X-API-Key: xxx   │
│  - Simpan di config         │
│  - Simple, cukup untuk MVP  │
├─────────────────────────────┤
│  Opsi B: Basic Auth         │
│  - Username + Password      │
│  - Lebih proper             │
│  - Tapi lebih kompleks      │
└─────────────────────────────┘
Untuk MVP: API Key sudah cukup.
5. Structured Logging
Sekarang cuma log.Printf ke stdout. Butuh:
Level logging:
├── INFO  → "Device 192.168.1.1 added"
├── WARN  → "Device 192.168.1.5 consecutive fails: 5"
├── ERROR → "Failed to ping 10.0.0.1: permission denied"
└── DEBUG → "Ping seq=12 latency=12.3ms"
Simpan ke file juga, bukan hanya stdout. Bisa pakai slog (bawaan Go 1.21+).
6. Dashboard Stats / History View
Ini terkait database, tapi fokus ke tampilan:
Yang perlu ditampilkan:
├── Daftar device + status terkini
├── Uptime per device (persentase)
├── Latency chart (grafik sederhana)
└── Last 24h activity log
Urutan Prioritas Rekomendasi
Phase 1 (Fondasi)
  7. Database Layer (SQLite)     ← paling dasar
  8. Configuration System        ← supaya tidak hardcoded

Phase 2 (Manajemen)
  3. Device Management (CRUD)    ← kelola device
  4. Simple Auth (API Key)       ← keamanan minimal

Phase 3 (Observability)
  5. Structured Logging          ← traceability
  6. Dashboard Stats / History   ← visibilitas data
Yang Tidak Perlu Dulu
- User management (multi-user, role) → terlalu over-engineering untuk MVP
- Email/SMS alerting → fitur lanjutan, fondasi belum siap
- Web UI untuk config → config file cukup untuk sekarang