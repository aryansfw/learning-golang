package handler

import (
	"errors"
	"fmt"
	"io"
	"md-note-api/model"
	"md-note-api/repository"
	"md-note-api/response"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Note struct {
	repo *repository.Note
}

func NewNote(db *pgxpool.Pool) *Note {
	return &Note{repo: repository.NewNote(db)}
}

func (h *Note) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(4 << 20); err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	uf, ufh, err := r.FormFile("note")
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest, err)
		return
	}
	defer uf.Close()

	ID := uuid.New()

	path := filepath.Join(os.Getenv("STORAGE_PATH"), ID.String())

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, uf); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	note := model.Note{
		ID:       ID,
		FileName: ufh.Filename,
	}

	if err := h.repo.Create(r.Context(), note); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	response.Success(w, "Successfully uploaded file", map[string]string{"id": ID.String()})
}

func (h *Note) Download(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := uuid.Parse(idString)

	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}

	note, err := h.repo.FindById(r.Context(), id)
	if err != nil {
		response.Error(w, err.Error(), http.StatusBadRequest, err)
		return
	}

	path := filepath.Join(os.Getenv("STORAGE_PATH"), idString)

	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		response.Error(w, "No file with given id found", http.StatusBadRequest, err)
		return
	}

	f, err := os.Open(path)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}
	defer f.Close()

	contentDisposition := fmt.Sprintf("attachment; filename=%s", note.FileName)
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError, err)
		return
	}
}
