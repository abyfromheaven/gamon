package database

// ============================================================
// MODULE DATABASE
// ============================================================
// Module ini bertanggung jawab untuk:
// 1. Menghubungkan aplikasi dengan database SQLite
// 2. Membuat tabel-tabel yang dibutuhkan secara otomatis
// 3. Menyimpan dan mengambil pengaturan aplikasi
//
// Database digunakan untuk menyimpan data perangkat, riwayat ping,
// alert/peringatan, dan pengaturan sistem GAMON.
// ============================================================

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Driver SQLite untuk Go tanpa perlu install额外 software
)

// dbFile adalah lokasi penyimpanan berkas database di dalam folder data/
const dbFile = "data/gamon.db"

// NewDB adalah fungsi utama untuk memulai koneksi database.
// Fungsi ini akan:
// - Membuat folder data/ jika belum ada
// - Membuka koneksi ke database SQLite
// - Menjalankan pembuatan tabel otomatis (migrasi)
// - Mengembalikan koneksi database yang siap digunakan
func NewDB() (*sql.DB, error) {
	// Buat folder data/ jika belum ada
	if err := os.MkdirAll(filepath.Dir(dbFile), 0755); err != nil {
		return nil, fmt.Errorf("gagal membuat direktori data: %w", err)
	}

	// Buka koneksi database dengan mode aman (WAL)
	// WAL = Write-Ahead Logging, agar data aman saat ada banyak akses bersamaan
	db, err := sql.Open("sqlite", dbFile+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("gagal membuka database: %w", err)
	}

	// Batasi jumlah koneksi ke 1 agar tidak terjadi konflik data
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Tes apakah koneksi database berhasil
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal melakukan tes koneksi ke database: %w", err)
	}

	// Jalankan pembuatan tabel otomatis
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("gagal menjalankan migrasi database: %w", err)
	}

	log.Println("Database SQLite terhubung dan berhasil dimigrasi")
	return db, nil
}

// migrate adalah fungsi untuk membuat semua tabel yang dibutuhkan.
// Fungsi ini hanya membuat tabel jika tabel belum ada (IF NOT EXISTS).
// Jika tabel sudah ada, fungsi ini tidak akan mengubah apapun.
func migrate(db *sql.DB) error {
	queries := []string{
		// ============================================
		// TABEL DEVICES (Data Perangkat yang Dipantau)
		// ============================================
		// Tabel ini menyimpan informasi tentang semua perangkat jaringan
		// yang ingin kita pantau statusnya (online/offline)
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
			last_online DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// ============================================
		// TABEL PING_HISTORY (Riwayat Hasil Pengecekan)
		// ============================================
		// Setiap kali perangkat dicek, hasilnya disimpan di sini.
		// Data ini digunakan untuk melihat riwayat kapan perangkat online/offline
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

		// ============================================
		// TABEL ALERTS (Peringatan/Sistem Alert)
		// ============================================
		// Tabel ini menyimpan catatan ketika perangkat bermasalah (offline).
		// Alert akan muncul di dashboard untuk memberitahu admin
		`CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'ongoing',
			severity TEXT NOT NULL DEFAULT 'low',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			resolved_at DATETIME,
			description TEXT DEFAULT '',
			alert_type TEXT DEFAULT 'critical',
			acknowledged BOOLEAN DEFAULT FALSE,
			acknowledged_at DATETIME,
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		)`,

		// ============================================
		// INDEX (Indeks untuk Mempercepat Pencarian)
		// ============================================
		// Index membuat query database lebih cepat, seperti indeks buku
		`CREATE INDEX IF NOT EXISTS idx_ping_history_device_id ON ping_history(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ping_history_timestamp ON ping_history(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_device_id ON alerts(device_id)`,
		`CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status)`,

		// ============================================
		// TABEL TELEGRAM_PAIRING (Koneksi Telegram Bot)
		// ============================================
		// Tabel ini menyimpan data koneksi antara GAMON dengan Telegram Bot
		// untuk mengirim notifikasi ke Telegram
		`CREATE TABLE IF NOT EXISTS telegram_pairing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			chat_id TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			paired_at DATETIME
		)`,

		// ============================================
		// TABEL SETTINGS (Pengaturan Aplikasi)
		// ============================================
		// Menyimpan pengaturan seperti ambang batas kegagalan ping
		`CREATE TABLE IF NOT EXISTS settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// Jalankan semua perintah pembuatan tabel
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("migrasi skema gagal: %w", err)
		}
	}

	// Migrasi tambahan untuk versi lama (menambah kolom baru jika belum ada)
	// PERINGATAN: Jangan hapus baris ini karena dibutuhkan untuk update dari versi lama
	_, _ = db.Exec("ALTER TABLE devices ADD COLUMN status TEXT DEFAULT 'active'")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN alert_type TEXT DEFAULT 'critical'")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN acknowledged BOOLEAN DEFAULT FALSE")
	_, _ = db.Exec("ALTER TABLE alerts ADD COLUMN acknowledged_at DATETIME")
	_, _ = db.Exec("ALTER TABLE devices ADD COLUMN last_online DATETIME")

	return nil
}

// ============================================================
// FUNGSI-FUNGSI PENGELOLA PENGATURAN
// ============================================================

// GetSetting mengambil nilai pengaturan dari database berdasarkan namanya.
// Contoh: GetSetting(db, "failure_threshold", "3") akan mengembalikan "3"
// jika pengaturan belum ada di database.
func GetSetting(db *sql.DB, key, defaultValue string) string {
	var value string
	err := db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetSetting menyimpan atau memperbarui pengaturan ke database.
// Jika pengaturan sudah ada, nilainya akan diupdate.
// Jika belum ada, pengaturan baru akan dibuat.
func SetSetting(db *sql.DB, key, value string) error {
	_, err := db.Exec(
		"INSERT INTO settings (key, value, updated_at) VALUES (?, ?, CURRENT_TIMESTAMP) ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP",
		key, value, value,
	)
	return err
}
