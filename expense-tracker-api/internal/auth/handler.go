package auth

import (
	"encoding/json"
	"net/http"
	"spendime/internal/config"
	"spendime/internal/response"
)

type Handler struct {
	svc *Service
	cfg *config.Config
}

func NewHandler(svc *Service, cfg *config.Config) *Handler {
	return &Handler{svc, cfg}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Invalid method", http.StatusMethodNotAllowed, "Invalid method")
		return
	}

	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, "Fields must not be empty", http.StatusBadRequest, "unauthorized")
		return
	}

	token, err := h.svc.Login(req.Email, req.Password, h.cfg.JWTSecret)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest, err.Error())
		return
	}

	data := map[string]string{"token": token}
	response.Success(w, "Login successful", data)
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, "Invalid method", http.StatusMethodNotAllowed, "invalid method")
		return
	}

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		response.Error(w, "Fields must not be empty", http.StatusBadRequest, "unauthorized")
		return
	}

	token, userID, err := h.svc.Register(req.Name, req.Email, req.Password, h.cfg.JWTSecret)
	if err != nil {
		response.Error(w, "Error registering account", http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]string{"token": token, "user_id": userID.String()}
	response.Success(w, "Registration successful", data)
}
