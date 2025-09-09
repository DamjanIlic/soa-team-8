package handler

import (
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
	json.NewDecoder(r.Body).Decode(&blog)
	h.BlogService.Create(&blog)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.BlogService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	blog, err := h.BlogService.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) Like(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.BlogService.Like(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blog liked"))
}

func (h *BlogHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.BlogService.Unlike(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Blog unliked"))
}
