package category

import (
	"net/http"
	"spendime/internal/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.List(r.Context())

	if err != nil {
		response.Error(w, "Error fetching categories", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully retrieved categories", categories)
}
