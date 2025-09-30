package handlers

import (
	"encoding/json"
	"net/http"
	"todo/internal/auth"
	"todo/internal/config"
	"todo/internal/services"
)

type AuthHandler struct {
	svc *services.AuthService
	cfg *config.Config
}

func NewAuthHandler(svc *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Fields must not be empty", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Login(req.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if req.Email != user.Email || req.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(*user, *h.cfg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Fields must not be empty", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(req.Name, req.Email, req.Password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateToken(*user, *h.cfg)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
