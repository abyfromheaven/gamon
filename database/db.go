package database

// ============================================================================
// MODUL DATABASE (database/db.go)
// ============================================================================
// Modul ini bertugas mengelola koneksi ke database SQLite (data/gamon.db),
// melakukan pembuat skema tabel otomatis (migrasi), serta fungsi pembantu
// untuk membaca/menyimpan pengaturan aplikasi (settings).
//
// Catatan Sidang PKL:
// Database yang digunakan adalah SQLite murni (tanpa GCC/CGO) dengan mode WAL
// (Write-Ahead Logging) agar proses baca-tulis data latensi jaringan cepat dan aman.
// ============================================================================

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Driver SQLite murni dalam bahasa Go (tanpa dependency CGO)
)

// berkasDatabase lokasi file penyimpanan database SQLite GAMON.
const berkasDatabase = "data/gamon.db"

// NewDB membuat folder data jika belum ada, membuka koneksi ke file SQLite,
// mengatur konfigurasi batas koneksi, serta otomatis memuat struktur tabel.
func NewDB() (*sql.DB, error) {
	// 1. Pastikan folder penyimpanan data ('data/') sudah tersedia
	folderPenyimpanan := filepath.Dir(berkasDatabase)
	if err := os.MkdirAll(folderPenyimpanan, 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori data: %w", err)
	}

	// 2. Buka koneksi database dengan mode WAL (Write-Ahead Logging) & Timeout 5000ms
	// Mode WAL memungkinkan pembacaan data (select) tidak terhalang oleh penulisan data (insert/update).
	koneksiDb, err := sql.Open("sqlite", berkasDatabase+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("gagal membuka database: %w", err)
	}

	// 3. Batasi koneksi maksimal 1 untuk SQLite demi menghindari konflik 'database locked'
	koneksiDb.SetMaxOpenConns(1)
	koneksiDb.SetMaxIdleConns(1)

	// 4. Pengujian koneksi awal ke database
	if err := koneksiDb.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan tes koneksi ke database: %w", err)
	}

	// 5. Jalankan proses pembuatan tabel otomatis (migrasi skema database)
	if err := migrate(koneksiDb); err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi database: %w", err)
	}

	log.Println("Database SQLite terhubung dan berhasil dimigrasi")
	return koneksiDb, nil
}

// migrate membuat tabel-tabel utama jika belum ada di database SQLite.
func migrate(koneksiDb *sql.DB) error {
	daftarPerintahSQL := []string{
		// Tabel Perangkat (devices) - Menyimpan data perangkat jaringan yang dipantau
		`CREATE TABLE IF NOT EXISTS devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			ip TEXT NOT NULL,
			url TEXT DEFAULT '',
			port INTEGER,
			method TEXT NOT NULL DEFAULT 'ICMP Ping',
			location TEXT DEFAULT '',
			check_interval INTEGER DEFAULT 3,
			status TEXT DEFAULT 'active',
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// Tabel Riwayat Ping (ping_history) - Menyimpan log latensi & status ping
		`CREATE TABLE IF NOT EXISTS ping_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			latency_ms REAL DEFAULT 0,
			ttl INTEGER DEFAULT 0,
			seq INTEGER DEFAULT 0,
			details TEXT DEFAULT '{}',
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		)`,
		// Tabel Peringatan (alerts) - Menyimpan riwayat kejadian perangkat down/offline
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'ongoing',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			resolved_at DATETIME,
			description TEXT DEFAULT '',
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		)`,
		// Pembuatan Indeks Pencarian (Index SQL) untuk mempercepat pencarian data
		`CREATE INDEX IF NOT EXISTS idx_ping_history_device_id ON ping_history(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ping_history_timestamp ON ping_history(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_device_id ON alerts(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,

		// Tabel Telegram Pairing - Menyimpan token hubung bot Telegram
		`CREATE TABLE IF NOT EXISTS telegram_pairing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			chat_id TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			paired_at DATETIME
		)`,

		// Tabel Pengaturan (settings) - Menyimpan opsi konfigurasi sistem (Kunci-Nilai / Key-Value)
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// Eksekusi seluruh skema tabel secara berurutan
	for _, perintah := range daftarPerintahSQL {
		if _, err := koneksiDb.Exec(perintah); err != nil {
			return fmt.Errorf("migrasi skema gagal: %w", err)
		}
	}

	// Migrasi tambahan untuk penyesuaian versi sebelumnya
	_, _ = koneksiDb.Exec("ALTER TABLE devices ADD COLUMN status TEXT DEFAULT 'active'")
	_, _ = koneksiDb.Exec("ALTER TABLE alerts ADD COLUMN alert_type TEXT DEFAULT 'critical'")
	_, _ = koneksiDb.Exec("ALTER TABLE alerts ADD COLUMN acknowledged BOOLEAN DEFAULT FALSE")
	_, _ = koneksiDb.Exec("ALTER TABLE alerts ADD COLUMN acknowledged_at DATETIME")

	return nil
}

// GetSetting mengambil nilai pengaturan berdasarkan kata kunci (kunciPengaturan).
func GetSetting(koneksiDb *sql.DB, kunciPengaturan, nilaiBawaan string) string {
	var nilaiPengaturan string
	err := koneksiDb.QueryRow("SELECT value FROM settings WHERE key = ?", kunciPengaturan).Scan(&nilaiPengaturan)
	if err != nil {
		return nilaiBawaan
	}
	return nilaiPengaturan
}

// SetSetting menyimpan atau memperbarui nilai pengaturan ke dalam tabel settings.
func SetSetting(koneksiDb *sql.DB, kunciPengaturan, nilaiPengaturan string) error {
	_, err := koneksiDb.Exec(
		"INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP",
		kunciPengaturan, nilaiPengaturan, nilaiPengaturan,
	)
	return err
}

