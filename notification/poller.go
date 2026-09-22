package notification

// ============================================================================
// MODUL TELEGRAM POLLER (notification/poller.go)
// ============================================================================
// Modul ini bertugas mengecek pesan masuk (perintah slash command) dari pengguna
// di aplikasi Telegram secara periodik (polling setiap 2 detik).
// Perintah yang didukung:
// 1. /pair GMN-XXXX-XXXX : Menghubungkan Telegram ke sistem GAMON.
// 2. /status              : Mengecek status koneksi Telegram.
// 3. /help                : Menampilkan daftar bantuan perintah bot.
// ============================================================================

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

// TelegramPoller mengelola perulangan polling perintah pesan dari Telegram API.
type TelegramPoller struct {
	botToken   string
	client     *http.Client
	pairFunc   func(token, chatID string) error
	statusFunc func(chatID string) bool
	offset     int
	chatIDMap  map[int64]string
}

// Format data internal pembaruan pesan dari Telegram API.
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

// Format respons pembungkus dari API Telegram.
type telegramResponse struct {
	Ok     bool             `json:"ok"`
	Result []telegramUpdate `json:"result"`
}

// NewTelegramPoller membuat instansi baru TelegramPoller.
func NewTelegramPoller(botToken string, pairFunc func(token, chatID string) error, statusFunc func(chatID string) bool) *TelegramPoller {
	return &TelegramPoller{
		botToken:   botToken,
		client:     &http.Client{Timeout: 10 * time.Second},
		pairFunc:   pairFunc,
		statusFunc: statusFunc,
		chatIDMap:  make(map[int64]string),
	}
}

// Start menjalankan goroutine polling interval setiap 2 detik di latar belakang.
func (p *TelegramPoller) Start() {
	if p.botToken == "" {
		log.Println("[Telegram Poller] Nonaktif (Bot token tidak diatur)")
		return
	}

	log.Println("[Telegram Poller] Berjalan (Polling pesan masuk setiap 2 detik)")

	penghitungWaktu := time.NewTicker(2 * time.Second)
	defer penghitungWaktu.Stop()

	for range penghitungWaktu.C {
		daftarPesanBaru, err := p.getUpdates()
		if err != nil {
			log.Printf("[Telegram Poller] Error mengambil pesan: %v", err)
			continue
		}

		for _, pesan := range daftarPesanBaru {
			p.offset = pesan.UpdateID + 1
			p.handleUpdate(pesan)
		}
	}
}

// getUpdates memanggil API Telegram endpoint 'getUpdates' untuk mengambil pesan baru.
func (p *TelegramPoller) getUpdates() ([]telegramUpdate, error) {
	alamatAPI := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", p.botToken)

	dataForm := url.Values{}
	dataForm.Set("offset", fmt.Sprintf("%d", p.offset))
	dataForm.Set("timeout", "1")
	dataForm.Set("allowed_updates", `["message"]`)

	respon, err := p.client.Post(alamatAPI, "application/x-www-form-urlencoded", strings.NewReader(dataForm.Encode()))
	if err != nil {
		return nil, err
	}
	defer respon.Body.Close()

	badanRespon, err := io.ReadAll(respon.Body)
	if err != nil {
		return nil, err
	}

	var hasil telegramResponse
	if err := json.Unmarshal(badanRespon, &hasil); err != nil {
		return nil, err
	}

	if !hasil.Ok {
		return nil, fmt.Errorf("API Telegram mengembalikan status ok=false")
	}

	return hasil.Result, nil
}

// handleUpdate mengevaluasi isi pesan yang diketik oleh pengguna di Telegram.
func (p *TelegramPoller) handleUpdate(pesan telegramUpdate) {
	teksPesan := pesan.Message.Text
	chatID := pesan.Message.Chat.ID
	namaPengirim := pesan.Message.From.FirstName

	if teksPesan == "" {
		return
	}

	// Evaluasi jenis perintah (slash command)
	if strings.HasPrefix(teksPesan, "/pair") {
		p.handlePair(teksPesan, chatID, namaPengirim)
	} else if teksPesan == "/status" {
		p.handleStatus(chatID)
	} else if teksPesan == "/unpair" {
		p.handleUnpair(chatID, namaPengirim)
	} else if teksPesan == "/start" || teksPesan == "/help" {
		p.handleHelp(chatID)
	}
}

// handlePair memproses token pairing (contoh: /pair GMN-1234-5678).
func (p *TelegramPoller) handlePair(teksPesan string, chatID int64, namaPengirim string) {
	bagianKata := strings.Fields(teksPesan)
	if len(bagianKata) < 2 {
		p.sendMessage(chatID, "Gunakan format: /pair GMN-XXXX-XXXX")
		return
	}

	tokenPairing := bagianKata[1]
	err := p.pairFunc(tokenPairing, fmt.Sprintf("%d", chatID))
	if err != nil {
		p.sendMessage(chatID, fmt.Sprintf("❌ Pairing gagal: %s", err.Error()))
		return
	}

	nama := namaPengirim
	if nama == "" {
		nama = "Admin"
	}
	p.sendMessage(chatID, fmt.Sprintf("✅ Pairing berhasil!\n\nHalo %s, Telegram kamu sudah terhubung dengan GAMON.\nSekarang kamu akan menerima notifikasi jika ada perangkat yang offline.", nama))
	log.Printf("[Telegram Poller] Pairing berhasil: chat_id=%d, dari=%s", chatID, namaPengirim)
}

// handleStatus mengecek status terhubung/tidaknya ID chat saat ini.
func (p *TelegramPoller) handleStatus(chatID int64) {
	chatIDTeks := fmt.Sprintf("%d", chatID)
	if p.statusFunc != nil && p.statusFunc(chatIDTeks) {
		p.sendMessage(chatID, "📊 Status: Terhubung ✅\n\nTelegram kamu aktif dan akan menerima notifikasi dari GAMON.")
	} else {
		p.sendMessage(chatID, "📊 Status: Belum Terhubung ❌\n\nTelegram belum terhubung. Gunakan /pair GMN-XXXX-XXXX untuk menghubungkan.")
	}
}

// handleUnpair memberikan arahan pemutusan koneksi Telegram via website.
func (p *TelegramPoller) handleUnpair(chatID int64, namaPengirim string) {
	p.sendMessage(chatID, "Untuk memutus hubungan Telegram, buka halaman Settings di website GAMON.")
}

// handleHelp menampilkan menu bantuan perintah bot Telegram.
func (p *TelegramPoller) handleHelp(chatID int64) {
	pesanBantuan := `🤖 GAMON Bot

Perintah yang tersedia:

/pair GMN-XXXX-XXXX
Hubungkan Telegram dengan GAMON

/status
Cek status koneksi

/help
Tampilkan bantuan ini

Untuk memulai, buka website GAMON → Settings → Telegram Integration → Connect Telegram, lalu kirim perintah /pair dengan token yang diberikan.`

	p.sendMessage(chatID, pesanBantuan)
}

// sendMessage mengirimkan pesan balasan dari bot ke pengguna Telegram.
func (p *TelegramPoller) sendMessage(chatID int64, isiPesan string) {
	alamatAPI := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", p.botToken)

	dataForm := url.Values{}
	dataForm.Set("chat_id", fmt.Sprintf("%d", chatID))
	dataForm.Set("text", isiPesan)

	respon, err := p.client.Post(alamatAPI, "application/x-www-form-urlencoded", strings.NewReader(dataForm.Encode()))
	if err != nil {
		log.Printf("[Telegram Poller] Gagal kirim pesan ke %d: %v", chatID, err)
		return
	}
	defer respon.Body.Close()

	if respon.StatusCode != http.StatusOK {
		badanRespon, _ := io.ReadAll(respon.Body)
		log.Printf("[Telegram Poller] Error HTTP %d: %s", respon.StatusCode, string(badanRespon))
	}
}
