package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"gamon/database"
)

type SettingsHandler struct {
	db *sql.DB
}

type SettingsResponse struct {
	FailureThreshold     int  `json:"failure_threshold"`
	NotificationsEnabled bool `json:"notifications_enabled"`
}

func NewSettingsHandler(db *sql.DB) *SettingsHandler {
	return &SettingsHandler{db: db}
}

func (h *SettingsHandler) HandleSettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSettings(w, r)
	case http.MethodPut:
		h.updateSettings(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *SettingsHandler) getSettings(w http.ResponseWriter, _ *http.Request) {
	failureThreshold := database.GetSetting(h.db, "failure_threshold", "3")
	notificationsEnabled := database.GetSetting(h.db, "notifications_enabled", "true")

	ft, _ := strconv.Atoi(failureThreshold)
	ne := notificationsEnabled == "true"

	respondData(w, SettingsResponse{
		FailureThreshold:     ft,
		NotificationsEnabled: ne,
	})
}

func (h *SettingsHandler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FailureThreshold     *int  `json:"failure_threshold"`
		NotificationsEnabled *bool `json:"notifications_enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.FailureThreshold != nil {
		if *req.FailureThreshold < 1 || *req.FailureThreshold > 10 {
			respondError(w, http.StatusBadRequest, "Failure threshold harus antara 1-10")
			return
		}
		database.SetSetting(h.db, "failure_threshold", strconv.Itoa(*req.FailureThreshold))
	}

	if req.NotificationsEnabled != nil {
		val := "true"
		if !*req.NotificationsEnabled {
			val = "false"
		}
		database.SetSetting(h.db, "notifications_enabled", val)
	}

	respondSuccess(w, "Settings updated")
}
