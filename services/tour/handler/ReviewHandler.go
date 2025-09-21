package handler

import (
	"encoding/json"
	"net/http"
	"tour/middleware"
	"tour/model"
	"tour/service"

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

	var req model.ReviewRequest

	if req.Rating < 1 || req.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	tourID := mux.Vars(r)["tourId"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.ReviewService.CreateReview(tourID, userID.(string), &req)
	if err != nil {
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
