package handler

import (
	"blog/model"
	"blog/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type BlogHandler struct {
	BlogService *service.BlogService
}

func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var blog model.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ako ID nije setovan, kreiraj novi
	if blog.ID == "" {
		blog.ID = model.NewBlog(blog.Title, blog.Content, blog.ImageURL).ID
	}

	if err := h.BlogService.Create(&blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.BlogService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	blog, err := h.BlogService.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) Like(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "missing userID", http.StatusBadRequest)
		return
	}

	count, err := h.BlogService.Like(blogID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Blog liked, total likes: %d", count)))
}

func (h *BlogHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "missing userID", http.StatusBadRequest)
		return
	}

	count, err := h.BlogService.Unlike(blogID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Blog unliked, total likes: %d", count)))
}
