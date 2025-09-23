package handler

import (
	"encoding/json"
	"net/http"
	"tour/service"

	"github.com/gorilla/mux"
)

type DurationHandler struct {
	DurationService *service.DurationService
}

func (h *DurationHandler) AddDuration(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["tourId"]

	var req struct {
		Transport string `json:"transport"`
		Minutes   int    `json:"minutes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	duration, err := h.DurationService.AddDuration(tourID, req.Transport, req.Minutes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(duration)
}

func (h *DurationHandler) GetDurationsByTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["tourId"]

	durations, err := h.DurationService.GetDurationsByTour(tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(durations)
}
