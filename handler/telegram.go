package handler

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type TelegramHandler struct {
	db *sql.DB
}

func NewTelegramHandler(db *sql.DB) *TelegramHandler {
	return &TelegramHandler{db: db}
}

func (h *TelegramHandler) HandlePair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Cek apakah sudah ada pairing aktif
	var activeCount int
	err := h.db.QueryRow("SELECT COUNT(*) FROM telegram_pairing WHERE status = 'connected'").Scan(&activeCount)
	if err == nil && activeCount > 0 {
		respondError(w, http.StatusConflict, "Telegram sudah terhubung. Disconnect terlebih dahulu.")
		return
	}

	// Generate token
	token, err := generateToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal generate token")
		return
	}

	// Simpan ke database (berlaku 30 hari)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	_, err = h.db.Exec(
		"INSERT INTO telegram_pairing (token, status, expires_at) VALUES (?, 'pending', ?)",
		token, expiresAt.Format(time.RFC3339),
	)
	if err != nil {
		log.Printf("Error creating pairing token: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal menyimpan token")
		return
	}

	log.Printf("[Telegram] Pairing token generated: %s", token)
	respondData(w, map[string]string{
		"token":      token,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *TelegramHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	type StatusResponse struct {
		Status   string  `json:"status"`
		ChatID   string  `json:"chat_id"`
		PairedAt *string `json:"paired_at"`
	}

	var status, chatID string
	var pairedAt sql.NullTime

	err := h.db.QueryRow(
		"SELECT status, chat_id, paired_at FROM telegram_pairing WHERE status = 'connected' ORDER BY paired_at DESC LIMIT 1",
	).Scan(&status, &chatID, &pairedAt)

	if err == sql.ErrNoRows {
		respondData(w, StatusResponse{Status: "disconnected"})
		return
	}
	if err != nil {
		log.Printf("Error getting telegram status: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil status")
		return
	}

	var pairedAtStr *string
	if pairedAt.Valid {
		s := pairedAt.Time.Format("02 Jan 2006 15:04")
		pairedAtStr = &s
	}

	respondData(w, StatusResponse{
		Status:   status,
		ChatID:   chatID,
		PairedAt: pairedAtStr,
	})
}

func (h *TelegramHandler) HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	result, err := h.db.Exec(
		"UPDATE telegram_pairing SET status = 'disconnected', chat_id = '' WHERE status = 'connected'",
	)
	if err != nil {
		log.Printf("Error disconnecting telegram: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal disconnect")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Tidak ada koneksi aktif")
		return
	}

	log.Printf("[Telegram] Account disconnected")
	respondSuccess(w, "Telegram disconnected")
}

func (h *TelegramHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Endpoint ini akan dipanggil oleh poller untuk memproses pairing
	// Bukan webhook dari Telegram langsung
	respondSuccess(w, "OK")
}

func (h *TelegramHandler) CompletePairing(token string, chatID string) error {
	// Validasi token
	var status string
	var expiresAt time.Time
	err := h.db.QueryRow(
		"SELECT status, expires_at FROM telegram_pairing WHERE token = ?", token,
	).Scan(&status, &expiresAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("token tidak ditemukan")
	}
	if err != nil {
		return err
	}

	if status != "pending" {
		return fmt.Errorf("token sudah digunakan atau tidak valid")
	}

	if time.Now().After(expiresAt) {
		h.db.Exec("UPDATE telegram_pairing SET status = 'expired' WHERE token = ?", token)
		return fmt.Errorf("token sudah kedaluwarsa")
	}

	// Simpan chat ID dan update status
	now := time.Now()
	_, err = h.db.Exec(
		"UPDATE telegram_pairing SET status = 'connected', chat_id = ?, paired_at = ? WHERE token = ?",
		chatID, now.Format(time.RFC3339), token,
	)
	if err != nil {
		return err
	}

	log.Printf("[Telegram] Pairing berhasil: token=%s, chat_id=%s", token, chatID)
	return nil
}

func (h *TelegramHandler) GetActiveChatID() string {
	var chatID string
	err := h.db.QueryRow(
		"SELECT chat_id FROM telegram_pairing WHERE status = 'connected' AND chat_id != '' ORDER BY paired_at DESC LIMIT 1",
	).Scan(&chatID)

	if err != nil {
		return ""
	}
	return chatID
}

func (h *TelegramHandler) IsActive() bool {
	var count int
	h.db.QueryRow("SELECT COUNT(*) FROM telegram_pairing WHERE status = 'connected' AND chat_id != ''").Scan(&count)
	return count > 0
}

func generateToken() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Format: GMN-XXXX-XXXX
	hexStr := hex.EncodeToString(b)
	upper := strings.ToUpper(hexStr)
	return fmt.Sprintf("GMN-%s-%s", upper[:4], upper[4:8]), nil
}
