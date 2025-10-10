package response

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    any `json:"data,omitempty"`
	Error   any `json:"error,omitempty"`
}

func sendJson(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to send data", http.StatusInternalServerError)
	}
}

func Success(w http.ResponseWriter, message string, data any) {
	sendJson(w, http.StatusOK, apiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, message string, status int, err any) {
	sendJson(w, status, apiResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
