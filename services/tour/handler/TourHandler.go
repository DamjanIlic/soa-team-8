package handler

import (
	"encoding/json"
	"net/http"
	"tour/model"
	"tour/service"

	"github.com/gorilla/mux"
)

type TourHandler struct {
	TourService *service.TourService
}

func (h *TourHandler) CreateTour(w http.ResponseWriter, r *http.Request) {
	// hardkodovan authorID, kasnije dodati auth
	authorID := r.Header.Get("Author-ID")
	if authorID == "" {
		http.Error(w, "Author-ID header required", http.StatusBadRequest)
		return
	}

	var req model.TourRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tour, err := h.TourService.CreateTour(authorID, &req)
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