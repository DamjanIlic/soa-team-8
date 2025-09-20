package handler

import (
	"encoding/json"
	"net/http"
	"tour/middleware"
	"tour/model"
	"tour/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TourHandler struct {
	TourService *service.TourService
}

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	role := r.Context().Value(middleware.ContextRole)

	if userID == nil || role == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if role.(string) != "guide" {
		http.Error(w, "Forbidden: only guides can create tours", http.StatusForbidden)
		return
	}

	var req model.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.CreateTour(userID.(string), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) GetTour(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	tour, err := h.TourService.GetTour(id)
	if err != nil {
		http.Error(w, "Tour not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) GetToursByAuthor(w http.ResponseWriter, r *http.Request) {
	authorID := mux.Vars(r)["authorId"]

	tours, err := h.TourService.GetToursByAuthor(authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {
	tours, err := h.TourService.GetAllTours()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

func (h *TourHandler) GetTourStatus(w http.ResponseWriter, r *http.Request) {
	tourIDStr := mux.Vars(r)["id"]

	_, err := uuid.Parse(tourIDStr)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	tourResponse, err := h.TourService.GetTour(tourIDStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "not_found"})
		return
	}

	var status string
	switch tourResponse.Status {
	case "published":
		status = "available" // moze da se kupi
	case "archived":
		status = "archived"
	case "draft":
		status = "draft"
	default:
		status = "not_available"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}
