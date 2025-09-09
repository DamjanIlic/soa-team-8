package handler

import (
	"blog/model"
	"blog/service"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CommentHandler struct {
	CommentService *service.CommentService
}

// POST /blogs/{id}/comments
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]

	var input struct {
		UserID string `json:"user_id"`
		Text   string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	comment := &model.Comment{
		BlogID: parseUUID(blogID),
		UserID: parseUUID(input.UserID),
		Text:   input.Text,
	}

	if err := h.CommentService.CreateComment(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GET /blogs/{id}/comments
func (h *CommentHandler) GetByBlogID(w http.ResponseWriter, r *http.Request) {
	blogID := mux.Vars(r)["id"]
	comments, err := h.CommentService.GetComments(blogID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

// Helper funkcija za parsiranje UUID iz stringa
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
