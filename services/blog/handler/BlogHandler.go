package handler

import (
	"blog/middleware" // ili gde god je tvoj middleware
	"blog/model"
	"blog/service"
	"encoding/json"
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

	userIDCtx := r.Context().Value(middleware.ContextUserID)
	if userIDCtx == nil {
		http.Error(w, "missing userID in context", http.StatusUnauthorized)
		return
	}
	userID := userIDCtx.(string)

	count, err := h.BlogService.Like(blogID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": count})
}

func (h *BlogHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	userIDCtx := r.Context().Value(middleware.ContextUserID)
	if userIDCtx == nil {
		http.Error(w, "missing userID in context", http.StatusUnauthorized)
		return
	}
	userID := userIDCtx.(string)

	count, err := h.BlogService.Unlike(blogID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"likes": count})
}
