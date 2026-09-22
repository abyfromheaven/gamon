package handler

// ============================================================================
// MODUL PEMBANTU RESPONSE API (handler/helpers.go)
// ============================================================================
// Modul ini menyediakan fungsi-fungsi standar untuk pembentukan format respons
// JSON API ke frontend, sehingga seluruh endpoint menghasilkan format pesan
// yang seragam (Clean Code & DRY - Don't Repeat Yourself).
// ============================================================================

import (
	"encoding/json"
	"net/http"
)

// APIResponse format standar balasan pesan status sukses/gagal tanpa data.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DataResponse format standar balasan pesan yang menyertakan data (payload).
type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// respondJSON mengirimkan respons berformat JSON dengan kode status HTTP yang ditentukan.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// respondError mengirimkan respons error JSON (Success: false).
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIResponse{Success: false, Message: message})
}

// respondSuccess mengirimkan respons sukses JSON (Success: true).
func respondSuccess(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: message})
}

// respondData mengirimkan respons data JSON (Success: true, Data: ...).
func respondData(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, DataResponse{Success: true, Data: data})
}
