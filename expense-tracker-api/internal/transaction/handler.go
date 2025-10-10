package transaction

import (
	"encoding/json"
	"log"
	"net/http"
	"spendime/internal/contextkey"
	"spendime/internal/response"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filters := TransactionFilters{
		DateFrom: parseTimePtr(q.Get("date_from")),
		DateTo:   parseTimePtr(q.Get("date_to")),
		Category: parseBoolPtr(q.Get("category")),
	}

	transactions, err := h.svc.List(r.Context(), filters)
	if err != nil {
		response.Error(w, "Error fetching transactions", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully retrieved transactions", transactions)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var transaction = Transaction{UserID: r.Context().Value(contextkey.UserID).(uuid.UUID)}

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	log.Println(transaction)

	if transaction.Name == "" || transaction.Type == "" || transaction.Amount == 0 || transaction.Date.IsZero() || transaction.CategoryID == nil {
		log.Println(transaction.Name == "", transaction.Type == "", transaction.Amount == 0, transaction.Date.IsZero(), transaction.CategoryID == nil)
		response.Error(w, "Fields must not be empty", http.StatusBadRequest, "empty fields")
		return
	}

	if !IsValidType(transaction.Type) {
		response.Error(w, "Invalid transaction type", http.StatusBadRequest, "invalid type")
		return
	}

	err := h.svc.Create(r.Context(), &transaction)
	if err != nil {
		response.Error(w, "Error creating new transaction", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully created transaction", transaction)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/transactions/")
	id, err := uuid.Parse(idString)
	if err != nil {
		response.Error(w, "Error parsing id", http.StatusBadRequest, err.Error())
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.Error(w, "Error deleting transaction", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully deleted transaction", nil)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/transactions/")
	id, err := uuid.Parse(idString)
	if err != nil {
		response.Error(w, "Error parsing id", http.StatusBadRequest, err.Error())
		return
	}

	var transaction = Transaction{ID: id}

	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if !IsValidType(transaction.Type) {
		response.Error(w, "Invalid transaction type", http.StatusBadRequest, "invalid type")
		return
	}

	if err := h.svc.Update(r.Context(), &transaction); err != nil {
		response.Error(w, "Error updating transaction", http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, "Successfully update transaction", transaction)
}

func parseTimePtr(s string) *time.Time {
	if s == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func parseBoolPtr(s string) *bool {
	if s == "" {
		return nil
	}

	var b bool

	switch s {
	case "true":
		b = true
	case "false":
		b = false
	default:
		return nil
	}

	return &b
}
