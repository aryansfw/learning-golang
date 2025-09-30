package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"todo/internal/auth"
	"todo/internal/models"
	"todo/internal/services"

	"github.com/google/uuid"
)

type TaskHandler struct {
	svc *services.TaskService
}

func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIdKey).(uuid.UUID)
	log.Println(userId)
	tasks, err := h.svc.ListTasks(userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !models.IsValidStatus(req.Status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(auth.UserIdKey).(uuid.UUID)

	task, err := h.svc.CreateTask(req.Title, req.Description, models.Status(req.Status), userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" || req.Status == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	if !models.IsValidStatus(req.Status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	task, err := h.svc.UpdateTask(id, req.Title, req.Description, models.Status(req.Status))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := uuid.Parse(idString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// w.Header().Set("Content-Type", "application/json")
	// if err := json.NewEncoder(w).Encode(task); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
}
