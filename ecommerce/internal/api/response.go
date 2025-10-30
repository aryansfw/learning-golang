package api

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool `json:"success"`
	Data    any  `json:"data,omitempty"`
	Error   any  `json:"error,omitempty"`
}

func (a *api) successResponse(w http.ResponseWriter, payload any) {
	res := response{
		Success: true,
		Data:    payload,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (a *api) errorResponse(w http.ResponseWriter, status int, err any) {
	res := response{
		Success: false,
		Error:   err,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
