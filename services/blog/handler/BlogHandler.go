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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if err := h.BlogService.Create(&blog); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.BlogService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	blog, err := h.BlogService.Get(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(blog)
}

func (h *BlogHandler) Like(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	// primer dobijanja userID; zameni po potrebi (npr. iz JWT tokena)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing userID"))
		return
	}

	count, err := h.BlogService.Like(blogID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Blog liked, total likes: %d", count)))
}

func (h *BlogHandler) Unlike(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing userID"))
		return
	}

	count, err := h.BlogService.Unlike(blogID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Blog unliked, total likes: %d", count)))
}
