package handler

import (
	"encoding/json"
	"net/http"
	"tour/middleware"
	"tour/model"
	"tour/service"

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
