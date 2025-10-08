package transaction

import (
	"encoding/json"
	"net/http"
	"spendime/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

type CreateRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Type == "" || req.Amount == 0 {
		response.Error(w, "Fields must not be empty", http.StatusBadRequest, "emtpy fields")
		return
	}

	transaction, err := h.svc.Create(req.Name, req.Amount, req.Type)
	if err != nil {
		response.Error(w, "Error creating new transaction", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully created transaction", transaction)
}
