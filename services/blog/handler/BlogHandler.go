package handler

import (
	"blog/metrics"
	"blog/middleware" // ili gde god je tvoj middleware
	"blog/model"
	"blog/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type BlogHandler struct {
	BlogService *service.BlogService
}

var tracer = otel.Tracer("blog-service")

func (h *BlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Dekodiraj body u privremeni struct
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Uzmi userID iz context-a
	userIDCtx := r.Context().Value(middleware.ContextUserID)
	if userIDCtx == nil {
		http.Error(w, "missing userID in context", http.StatusUnauthorized)
		return
	}
	userID := userIDCtx.(string)

	// Kreiraj novi blog sa userID-om
	blog := model.NewBlog(input.Title, input.Content, input.ImageURL, userID)

	// Sačuvaj u servisu
	if err := h.BlogService.Create(blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blog)

	metrics.BlogPostsProcessed.Inc()
}

func (h *BlogHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Start tracing span
	_, span := tracer.Start(r.Context(), "GetAllBlogs")
	defer span.End()

	blogs, err := h.BlogService.GetAll() // možeš proslediti ctx ako želiš da pratiš deeper calls
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

func (h *BlogHandler) GetForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		http.Error(w, "userId path param required", http.StatusBadRequest)
		return
	}

	blogs, err := h.BlogService.GetForUser(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}

func (h *BlogHandler) GetMyBlogs(w http.ResponseWriter, r *http.Request) {
	userIDCtx := r.Context().Value(middleware.ContextUserID)
	if userIDCtx == nil {
		http.Error(w, "missing userID in context", http.StatusUnauthorized)
		return
	}
	userID := userIDCtx.(string)

	blogs, err := h.BlogService.GetMyBlogs(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blogs)
}
