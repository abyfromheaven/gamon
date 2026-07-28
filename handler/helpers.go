package handler

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, APIResponse{Success: false, Message: message})
}

func respondSuccess(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusOK, APIResponse{Success: true, Message: message})
}

func respondData(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, DataResponse{Success: true, Data: data})
}
