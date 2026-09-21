package database

// Module Database bertanggung jawab mengelola koneksi ke SQLite database (data/gamon.db),
// melakukan migrasi otomatis tabel-tabel, serta fungsi pembantu pengaturan aplikasi (settings).

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Driver SQLite murni dalam bahasa Go (tanpa CGO/GCC)
)

// dbFile adalah lokasi penyimpanan berkas database SQLite GAMON.
const dbFile = "data/gamon.db"

// NewDB membuat folder data jika belum ada, membuka koneksi ke berkas SQLite,
// mengatur konfigurasi koneksi, serta otomatis menjalankan struktur tabel (migrasi).
func NewDB() (*sql.DB, error) {
	// 1. Pastikan direktori 'data/' sudah dibuat
	if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori data: %w", err)
	}

	// 2. Buka koneksi database dengan mode WAL (Write-Ahead Logging) agar cepat & aman saat concurrent access
	db, err := sql.Open("sqlite", dbFile+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("gagal membuka database: %w", err)
	}

	// 3. Batasi koneksi maksimal 1 untuk SQLite menghindari masalah database locked
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// 4. Tes koneksi database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan tes koneksi ke database: %w", err)
	}

	// 5. Jalankan pembuat skema tabel otomatis (migrasi)
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi database: %w", err)
	}

	log.Println("Database SQLite terhubung dan berhasil dimigrasi")
	return db, nil
}

// migrate membuat tabel-tabel utama database jika tabel belum ada saat pertama kali aplikasi dijalankan.
func migrate(db *sql.DB) error {
	queries := []string{
		// Tabel Perangkat (devices)
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
		// Tabel Riwayat Ping (ping_history)
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
		// Tabel Peringatan / Masalah (alerts)
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
		// Pembuatan Indeks Pencarian (Index) agar query cepat
		`CREATE INDEX IF NOT EXISTS idx_ping_history_device_id ON ping_history(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ping_history_timestamp ON ping_history(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_device_id ON alerts(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,

		// Tabel Pairing Telegram Bot (telegram_pairing)
		`CREATE TABLE IF NOT EXISTS telegram_pairing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			chat_id TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			paired_at DATETIME
		)`,

		// Tabel Pengaturan Aplikasi (settings)
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migrasi skema gagal: %w", err)
		}
	}

	// Migrasi bertahap jika ada penambahan kolom pada versi sebelumnya
	_, _ = db.Exec("ALTER TABLE devices ADD COLUMN status TEXT DEFAULT 'active'")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN alert_type TEXT DEFAULT 'critical'")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN acknowledged BOOLEAN DEFAULT FALSE")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN acknowledged_at DATETIME")

	return nil
}

// GetSetting mengambil nilai pengaturan berdasarkan kunci (key) dari tabel settings.
func GetSetting(db *sql.DB, key, defaultValue string) string {
	var value string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetSetting menyimpan atau memperbarui nilai pengaturan ke tabel settings.
func SetSetting(db *sql.DB, key, value string) error {
	_, err := db.Exec(
		"INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP",
		key, value, value,
	)
	return err
}
