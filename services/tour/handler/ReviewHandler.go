package handler

import (
	"encoding/json"
	"net/http"
	"tour/middleware"
	"tour/model"
	"tour/service"

	"log"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	ReviewService *service.ReviewService
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	role := r.Context().Value(middleware.ContextRole)

	if userID == nil || role == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if role.(string) != "tourist" {
		http.Error(w, "Forbidden: only tourists can leave reviews", http.StatusForbidden)
		return
	}

	tourID := mux.Vars(r)["tourId"]

	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Decoded raw map: %+v", raw)

	var req model.ReviewRequest
	if rating, ok := raw["rating"].(float64); ok {
		req.Rating = int(rating)
	}
	if comment, ok := raw["comment"].(string); ok {
		req.Comment = comment
	}
	if visitedAt, ok := raw["visited_at"].(string); ok {
		req.VisitedAt = visitedAt
	}
	if images, ok := raw["images"].([]interface{}); ok {
		for _, img := range images {
			if s, ok := img.(string); ok {
				req.Images = append(req.Images, s)
			}
		}
	}

	log.Printf("Mapped ReviewRequest: %+v", req)

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	review, err := h.ReviewService.CreateReview(tourID, userID.(string), &req)
	if err != nil {
		log.Printf("CreateReview error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetReviewsByTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["tourId"]

	reviews, err := h.ReviewService.GetReviewsByTour(tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// reviews za admina
func (h *ReviewHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.ContextRole)
	if role == nil || role.(string) != "admin" {
		http.Error(w, "Forbidden: admin only", http.StatusForbidden)
		return
	}

	reviews, err := h.ReviewService.GetAllReviews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// reviews od vodica
func (h *ReviewHandler) GetReviewsForGuide(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	role := r.Context().Value(middleware.ContextRole)

	if userID == nil || role == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if role.(string) != "guide" {
		http.Error(w, "Forbidden: guides only", http.StatusForbidden)
		return
	}

	reviews, err := h.ReviewService.GetReviewsForGuide(userID.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// reviews za turistu
func (h *ReviewHandler) GetMyReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	role := r.Context().Value(middleware.ContextRole)

	if userID == nil || role == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if role.(string) != "tourist" {
		http.Error(w, "Forbidden: tourists only", http.StatusForbidden)
		return
	}

	touristID, err := uuid.Parse(userID.(string))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	reviews, err := h.ReviewService.GetReviewsByTourist(touristID.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
