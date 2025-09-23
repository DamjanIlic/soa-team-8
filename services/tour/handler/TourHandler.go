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

// zajedno sa getmytours nabavlja ture od vodica
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

func (h *TourHandler) GetAuthorTours(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := r.Context().Value(middleware.ContextRole).(string)
	if !ok || role != "guide" {
		http.Error(w, "Only guides can access this endpoint", http.StatusForbidden)
		return
	}

	tours, err := h.TourService.GetToursByAuthor(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tours) // i ako je [] frontend dobija 200 + []
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

func (h *TourHandler) PublishTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["id"]
	authorID := r.Context().Value(middleware.ContextUserID).(string)

	tour, err := h.TourService.PublishTour(tourID, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) ArchiveTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["id"]
	authorID := r.Context().Value(middleware.ContextUserID).(string)

	tour, err := h.TourService.ArchiveTour(tourID, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(tour)
}

func (h *TourHandler) ReactivateTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["id"]
	authorID := r.Context().Value(middleware.ContextUserID).(string)

	tour, err := h.TourService.ReactivateTour(tourID, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(tour)
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

func (h *TourHandler) UpdateDistance(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["id"]

	var payload struct {
		DistanceKm float64 `json:"distance_km"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.UpdateDistance(tourID, payload.DistanceKm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

// handler/tour_handler.go
func (h *TourHandler) UpdatePrice(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["id"]
	authorID := r.Context().Value(middleware.ContextUserID).(string)

	// Parse request body
	var payload struct {
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.Price <= 0 {
		http.Error(w, "Price must be greater than 0", http.StatusBadRequest)
		return
	}

	// Check author
	tour, err := h.TourService.GetTour(tourID)
	if err != nil || tour.AuthorID != authorID {
		http.Error(w, "Unauthorized or tour not found", http.StatusForbidden)
		return
	}

	// Update price
	updatedTour, err := h.TourService.UpdatePrice(tourID, payload.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTour)
}
