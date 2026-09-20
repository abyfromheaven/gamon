package notification

// Module TelegramNotifier bertanggung jawab mengirimkan pesan notifikasi peringatan (alert)
// dan notifikasi pemulihan (recovery) langsung ke aplikasi Telegram admin melalui Telegram Bot API.

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

// TelegramNotifier menyimpan konfigurasi token bot, koneksi database, dan HTTP client.
type TelegramNotifier struct {
	botToken string
	db       *sql.DB
	enabled  bool
	client   *http.Client
}

// NewTelegramNotifier menginisialisasi notifier Telegram.
func NewTelegramNotifier(db *sql.DB) *TelegramNotifier {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	return &TelegramNotifier{
		botToken: botToken,
		db:       db,
		enabled:  botToken != "",
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// IsEnabled mengecek apakah fitur bot Telegram diaktifkan.
func (t *TelegramNotifier) IsEnabled() bool {
	return t.enabled
}

// SendAlert mengirimkan pesan notifikasi PERINGATAN (OFFLINE) saat ada perangkat jaringan yang mati.
func (t *TelegramNotifier) SendAlert(deviceName, deviceIP string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi Telegram aktif, skip notifikasi")
		return
	}

	msg := fmt.Sprintf(
		"🚨 *GAMON ALERT*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* OFFLINE\n"+
			"*Waktu:* %s\n\n"+
			"Perangkat tidak merespons ICMP Ping.",
		deviceName, deviceIP, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, msg)
}

// SendRecovery mengirimkan pesan notifikasi PEMULIHAN (ONLINE) saat perangkat yang sebelumnya mati kembali normal.
func (t *TelegramNotifier) SendRecovery(deviceName, deviceIP string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi Telegram aktif, skip notifikasi")
		return
	}

	msg := fmt.Sprintf(
		"✅ *GAMON RECOVERY*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* ONLINE\n"+
			"*Waktu:* %s\n\n"+
			"Perangkat kembali online dan merespons normal.",
		deviceName, deviceIP, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, msg)
}

// getActiveChatID mengambil ID Chat Telegram milik admin yang telah berhasil dipairing di database.
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

// send mengirimkan permintaan HTTP POST ke Telegram Bot API untuk meneruskan pesan ke HP admin.
func (t *TelegramNotifier) send(chatID string, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", text)
	data.Set("parse_mode", "Markdown")

	resp, err := t.client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("[Telegram] Gagal mengirim pesan: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Telegram] Error HTTP %d: %s", resp.StatusCode, string(body))
		return
	}

	log.Printf("[Telegram] Berhasil mengirim pesan notifikasi ke Chat ID %s", chatID)
}
