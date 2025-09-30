package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type BlogHandler struct {
	Service *BlogService
}

func NewBlogHandler(s *BlogService) *BlogHandler {
	return &BlogHandler{Service: s}
}

func (h *BlogHandler) HandleGetBlogs(w http.ResponseWriter, r *http.Request) {
	blogs := h.Service.GetBlogs()

	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		log.Println("Failed to encode", blogs)
		return
	}
	log.Println("Successfully returned", blogs)
}

func (h *BlogHandler) HandleGetBlogById(w http.ResponseWriter, r *http.Request) {
	var id = r.PathValue("id")

	blog := h.Service.GetBlogById(id)

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Println("Failed to encode", blog)
		return
	}
	log.Println("Successfully returned", blog)
}

func (h *BlogHandler) HandleCreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog Blog

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		log.Println("Failed to decode", r.Body)
		return
	}

	blog = h.Service.CreateBlog(blog)

	if err := json.NewEncoder(w).Encode(blog); err != nil {
		log.Println("Failed to encode", blog)
		return
	}
	log.Println("Successfully returned", blog)
}

func (h *BlogHandler) HandleUpdateBlog(w http.ResponseWriter, r *http.Request) {
	var id = r.PathValue("id")
	var blog Blog

	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		log.Println("Failed to decode", r.Body)
		return
	}

	switch r.Method {
	case http.MethodPatch:
		blog, err := h.Service.UpdateBlog(id, blog)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully updated blog with id", id)
	case http.MethodPut:
		blog, err := h.Service.ReplaceBlog(id, blog)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		log.Println("Successfully updated blog with id", id)
	}
}

func (h *BlogHandler) HandleDeleteBlog(w http.ResponseWriter, r *http.Request) {
	var id = r.PathValue("id")

	h.Service.DeleteBlog(id)

	log.Println("Successfully deleted blog with id", id)
}
