package notification

// ============================================================================
// MODUL NOTIFIKASI TELEGRAM BOT (notification/telegram.go)
// ============================================================================
// Modul ini bertugas mengirimkan pesan notifikasi Peringatan (Alert) dan Pemulihan (Recovery)
// langsung ke aplikasi Telegram milik Admin / Tim IT melalui Telegram Bot API.
//
// Alur Kerja:
// 1. Saat perangkat mati (offline), fungsi SendAlert dipanggil.
// 2. Notifier membaca Chat ID aktif dari tabel 'telegram_pairing' di SQLite.
// 3. Notifier melakukan permintaan HTTP POST ke URL API Telegram (https://api.telegram.org/bot<TOKEN>/sendMessage).
// 4. Pesan dalam format Markdown terkirim langsung ke HP Admin.
// ============================================================================

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// TelegramNotifier menyimpan token bot Telegram, koneksi database, dan HTTP client.
type TelegramNotifier struct {
	botToken string
	db       *sql.DB
	enabled  bool
	client   *http.Client
}

// NewTelegramNotifier menginisialisasi Notifier Telegram.
func NewTelegramNotifier(db *sql.DB) *TelegramNotifier {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	return &TelegramNotifier{
		botToken: botToken,
		db:       db,
		enabled:  botToken != "",
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// IsEnabled mengecek apakah variabel lingkungan TELEGRAM_BOT_TOKEN diatur (aktif/nonaktif).
func (t *TelegramNotifier) IsEnabled() bool {
	return t.enabled
}

// SendAlert mengirimkan pesan notifikasi PERINGATAN (OFFLINE) ke Telegram Admin.
func (t *TelegramNotifier) SendAlert(namaPerangkat, ipPerangkat string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi Telegram aktif, lewati pengiriman notifikasi")
		return
	}

	pesan := fmt.Sprintf(
		"🚨 *GAMON ALERT*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* OFFLINE\n"+
			"*Waktu:* %s\n\n"+
			"Perangkat tidak merespons ICMP Ping.",
		namaPerangkat, ipPerangkat, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, pesan)
}

// SendRecovery mengirimkan pesan notifikasi PEMULIHAN (ONLINE) saat perangkat kembali normal.
func (t *TelegramNotifier) SendRecovery(namaPerangkat, ipPerangkat string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi Telegram aktif, lewati pengiriman notifikasi")
		return
	}

	pesan := fmt.Sprintf(
		"✅ *GAMON RECOVERY*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* ONLINE\n"+
			"*Waktu:* %s\n\n"+
			"Perangkat kembali online dan merespons normal.",
		namaPerangkat, ipPerangkat, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, pesan)
}

// getActiveChatID mengambil ID Chat Telegram Admin dari database SQLite yang sudah di-pair.
func (t *TelegramNotifier) getActiveChatID() string {
	var chatID string
	err := t.db.QueryRow(
		"SELECT chat_id FROM telegram_pairing WHERE status = 'connected' AND chat_id != '' ORDER BY paired_at DESC LIMIT 1",
	).Scan(&chatID)

	if err != nil {
		return ""
	}
	return chatID
}

// send mengirimkan permintaan HTTP POST ke Telegram Bot API endpoint 'sendMessage'.
func (t *TelegramNotifier) send(chatID string, isiPesan string) {
	alamatAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	dataForm := url.Values{}
	dataForm.Set("chat_id", chatID)
	dataForm.Set("text", isiPesan)
	dataForm.Set("parse_mode", "Markdown")

	respon, err := t.client.Post(alamatAPI, "application/x-www-form-urlencoded", strings.NewReader(dataForm.Encode()))
	if err != nil {
		log.Printf("[Telegram] Gagal mengirim pesan: %v", err)
		return
	}
	defer respon.Body.Close()

	if respon.StatusCode != http.StatusOK {
		badanRespon, _ := io.ReadAll(respon.Body)
		log.Printf("[Telegram] Error HTTP %d: %s", respon.StatusCode, string(badanRespon))
		return
	}

	log.Printf("[Telegram] Berhasil mengirim pesan notifikasi ke Chat ID %s", chatID)
}
