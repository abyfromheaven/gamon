package handler

// ============================================================================
// MODUL HANDLER PENGATURAN (handler/settings.go)
// ============================================================================
// Modul ini mengelola pengaturan aplikasi (Settings) yang tersimpan di SQLite,
// seperti:
// 1. failure_threshold: Ambang batas kegagalan ping berturut-turut sebelum dianggap Offline (default: 3 kali).
// 2. notifications_enabled: Pengaturan pengaktifan notifikasi.
// ============================================================================

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"gamon/database"
)

// SettingsHandler menyimpan koneksi database SQLite.
type SettingsHandler struct {
	db *sql.DB
}

// SettingsResponse format data pengaturan aplikasi untuk dikirim ke frontend.
type SettingsResponse struct {
	FailureThreshold     int  `json:"failure_threshold"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

// NewSettingsHandler membuat instansi baru SettingsHandler.
func NewSettingsHandler(db *sql.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

// HandleSettings memproses rute /api/settings (GET untuk membaca, PUT untuk memperbarui).
func (h *SettingsHandler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSettings(w, r)
	case http.MethodPut:
		h.updateSettings(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// getSettings membaca pengaturan saat ini dari database SQLite.
func (h *SettingsHandler) getSettings(w http.ResponseWriter, _ *http.Request) {
	ambangGagalTeks := database.GetSetting(h.db, "failure_threshold", "3")
	notifAktifTeks := database.GetSetting(h.db, "notifications_enabled", "true")

	ambangGagal, _ := strconv.Atoi(ambangGagalTeks)
	notifAktif := notifAktifTeks == "true"

	respondData(w, SettingsResponse{
		FailureThreshold:     ambangGagal,
		NotificationsEnabled: notifAktif,
	})
}

// updateSettings memperbarui nilai pengaturan aplikasi di database SQLite.
func (h *SettingsHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FailureThreshold     *int  `json:"failure_threshold"`
		NotificationsEnabled *bool `json:"notifications_enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if req.FailureThreshold != nil {
		if *req.FailureThreshold < 1 || *req.FailureThreshold > 10 {
			respondError(w, http.StatusBadRequest, "Ambang batas kegagalan harus antara 1-10")
			return
		}
		database.SetSetting(h.db, "failure_threshold", strconv.Itoa(*req.FailureThreshold))
	}

	if req.NotificationsEnabled != nil {
		nilaiNotif := "true"
		if !*req.NotificationsEnabled {
			nilaiNotif = "false"
		}
		database.SetSetting(h.db, "notifications_enabled", nilaiNotif)
	}

	respondSuccess(w, "Pengaturan berhasil diperbarui")
}
