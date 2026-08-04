package notification

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

type TelegramNotifier struct {
	botToken string
	db       *sql.DB
	enabled  bool
	client   *http.Client
}

func NewTelegramNotifier(db *sql.DB) *TelegramNotifier {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	return &TelegramNotifier{
		botToken: botToken,
		db:       db,
		enabled:  botToken != "",
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *TelegramNotifier) IsEnabled() bool {
	return t.enabled
}

func (t *TelegramNotifier) SendAlert(deviceName, deviceIP string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi aktif, skip notifikasi")
		return
	}

	msg := fmt.Sprintf(
		"🚨 *GAMON ALERT*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* OFFLINE\n"+
			"*Time:* %s\n\n"+
			"Device tidak merespons ping.",
		deviceName, deviceIP, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, msg)
}

func (t *TelegramNotifier) SendRecovery(deviceName, deviceIP string) {
	if !t.enabled {
		return
	}

	chatID := t.getActiveChatID()
	if chatID == "" {
		log.Println("[Telegram] Tidak ada koneksi aktif, skip notifikasi")
		return
	}

	msg := fmt.Sprintf(
		"✅ *GAMON RECOVERY*\n\n"+
			"*Device:* %s\n"+
			"*IP:* %s\n"+
			"*Status:* ONLINE\n"+
			"*Time:* %s\n\n"+
			"Device kembali online.",
		deviceName, deviceIP, time.Now().Format("02 Jan 2006 15:04:05"),
	)

	t.send(chatID, msg)
}

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

func (t *TelegramNotifier) send(chatID string, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", text)
	data.Set("parse_mode", "Markdown")

	resp, err := t.client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("[Telegram] Gagal kirim: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Telegram] Error %d: %s", resp.StatusCode, string(body))
		return
	}

	log.Printf("[Telegram] Pesan terkirim ke chat %s", chatID)
}
