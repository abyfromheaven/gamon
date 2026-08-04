package notification

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TelegramPoller struct {
	botToken   string
	client     *http.Client
	pairFunc   func(token, chatID string) error
	statusFunc func(chatID string) bool
	offset     int
	chatIDMap  map[int64]string // chatID -> pairing status
}

type telegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		From struct {
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
	} `json:"message"`
}

type telegramResponse struct {
	Ok     bool              `json:"ok"`
	Result []telegramUpdate  `json:"result"`
}

func NewTelegramPoller(botToken string, pairFunc func(token, chatID string) error, statusFunc func(chatID string) bool) *TelegramPoller {
	return &TelegramPoller{
		botToken:   botToken,
		client:     &http.Client{Timeout: 10 * time.Second},
		pairFunc:   pairFunc,
		statusFunc: statusFunc,
		chatIDMap:  make(map[int64]string),
	}
}

func (p *TelegramPoller) Start() {
	if p.botToken == "" {
		log.Println("[Telegram Poller] Disabled (no bot token)")
		return
	}

	log.Println("[Telegram Poller] Started (polling setiap 2 detik)")

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updates, err := p.getUpdates()
		if err != nil {
			log.Printf("[Telegram Poller] Error getting updates: %v", err)
			continue
		}

		for _, update := range updates {
			p.offset = update.UpdateID + 1
			p.handleUpdate(update)
		}
	}
}

func (p *TelegramPoller) getUpdates() ([]telegramUpdate, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", p.botToken)

	data := url.Values{}
	data.Set("offset", fmt.Sprintf("%d", p.offset))
	data.Set("timeout", "1")
	data.Set("allowed_updates", `["message"]`)

	resp, err := p.client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result telegramResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Ok {
		return nil, fmt.Errorf("API returned ok=false")
	}

	return result.Result, nil
}

func (p *TelegramPoller) handleUpdate(update telegramUpdate) {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	fromName := update.Message.From.FirstName

	if text == "" {
		return
	}

	// Handle commands
	if strings.HasPrefix(text, "/pair") {
		p.handlePair(text, chatID, fromName)
	} else if text == "/status" {
		p.handleStatus(chatID)
	} else if text == "/unpair" {
		p.handleUnpair(chatID, fromName)
	} else if text == "/start" || text == "/help" {
		p.handleHelp(chatID)
	}
}

func (p *TelegramPoller) handlePair(text string, chatID int64, fromName string) {
	parts := strings.Fields(text)
	if len(parts) < 2 {
		p.sendMessage(chatID, "Gunakan: /pair GMN-XXXX-XXXX")
		return
	}

	token := parts[1]
	err := p.pairFunc(token, fmt.Sprintf("%d", chatID))
	if err != nil {
		p.sendMessage(chatID, fmt.Sprintf("❌ Pairing gagal: %s", err.Error()))
		return
	}

	name := fromName
	if name == "" {
		name = "Admin"
	}
	p.sendMessage(chatID, fmt.Sprintf("✅ Pairing berhasil!\n\nHalo %s, Telegram kamu sudah terhubung dengan GAMON.\nSekarang kamu akan menerima notifikasi jika ada device yang offline.", name))
	log.Printf("[Telegram Poller] Pairing berhasil: chat_id=%d, from=%s", chatID, fromName)
}

func (p *TelegramPoller) handleStatus(chatID int64) {
	chatIDStr := fmt.Sprintf("%d", chatID)
	if p.statusFunc != nil && p.statusFunc(chatIDStr) {
		p.sendMessage(chatID, "📊 Status: Connected ✅\n\nTelegram kamu aktif dan akan menerima notifikasi dari GAMON.")
	} else {
		p.sendMessage(chatID, "📊 Status: Not Connected ❌\n\nTelegram belum terhubung. Gunakan /pair GMN-XXXX-XXXX untuk menghubungkan.")
	}
}

func (p *TelegramPoller) handleUnpair(chatID int64, fromName string) {
	// Note: Untuk unpair, kita perlu akses database. 
	// Sementara cukup beri info bahwa unpair harus dilakukan dari web.
	p.sendMessage(chatID, "Untuk memutus pairing, buka halaman Settings di website GAMON.")
}

func (p *TelegramPoller) handleHelp(chatID int64) {
	msg := `🤖 GAMON Bot

Perintah yang tersedia:

/pair GMN-XXXX-XXXX
Hubungkan Telegram dengan GAMON

/status
Cek status koneksi

/help
Tampilkan bantuan ini

Untuk memulai, buka website GAMON → Settings → Telegram Integration → Connect Telegram, lalu kirim perintah /pair dengan token yang diberikan.`

	p.sendMessage(chatID, msg)
}

func (p *TelegramPoller) sendMessage(chatID int64, text string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", p.botToken)

	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", chatID))
	data.Set("text", text)

	resp, err := p.client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("[Telegram Poller] Gagal kirim pesan ke %d: %v", chatID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[Telegram Poller] Error %d: %s", resp.StatusCode, string(body))
	}
}
