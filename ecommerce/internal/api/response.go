package api

import (
	"encoding/json"
	"net/http"
)

func (a *api) successResponse(w http.ResponseWriter, payload any) {
	res := map[string]any{"success": true, "data": payload}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (a *api) errorResponse(w http.ResponseWriter, status int, err any) {
	res := map[string]any{"success": false, "error": err}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(res)
}
