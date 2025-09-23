package handler

import (
	"encoding/json"
	"net/http"

	"tour/middleware"
	"tour/service"

	"github.com/google/uuid"
)

type TourExecutionHandler struct {
	service *service.TourExecutionService
}

func NewTourExecutionHandler(s *service.TourExecutionService) *TourExecutionHandler {
	return &TourExecutionHandler{service: s}
}

// helper funkcija za proveru role
func isTourist(r *http.Request) bool {
	role, ok := r.Context().Value(middleware.ContextRole).(string)
	return ok && role == "tourist"
}

func (h *TourExecutionHandler) StartTour(w http.ResponseWriter, r *http.Request) {
	if !isTourist(r) {
		http.Error(w, "forbidden: only tourists can start tours", http.StatusForbidden)
		return
	}

	userIDStr := r.Context().Value(middleware.ContextUserID).(string)
	var req struct {
		TourID string `json:"tour_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userID, _ := uuid.Parse(userIDStr)
	tourID, _ := uuid.Parse(req.TourID)

	exec, err := h.service.StartTour(userID, tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(exec)
}

func (h *TourExecutionHandler) CompleteTour(w http.ResponseWriter, r *http.Request) {
	if !isTourist(r) {
		http.Error(w, "forbidden: only tourists can complete tours", http.StatusForbidden)
		return
	}

	var req struct {
		ExecID string `json:"exec_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	execID, _ := uuid.Parse(req.ExecID)

	if err := h.service.CompleteTour(execID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) AbandonTour(w http.ResponseWriter, r *http.Request) {
	if !isTourist(r) {
		http.Error(w, "forbidden: only tourists can abandon tours", http.StatusForbidden)
		return
	}

	var req struct {
		ExecID string `json:"exec_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	execID, _ := uuid.Parse(req.ExecID)

	if err := h.service.AbandonTour(execID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TourExecutionHandler) CheckKeyPoint(w http.ResponseWriter, r *http.Request) {
	if !isTourist(r) {
		http.Error(w, "forbidden: only tourists can check keypoints", http.StatusForbidden)
		return
	}

	var req struct {
		ExecID     string  `json:"exec_id"`
		KeyPointID string  `json:"key_point_id"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	execID, _ := uuid.Parse(req.ExecID)
	keyPointID, _ := uuid.Parse(req.KeyPointID)

	reached, err := h.service.CheckKeyPoint(execID, keyPointID, req.Latitude, req.Longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"reached": reached})
}

func (h *TourExecutionHandler) GetSimulatedPosition(w http.ResponseWriter, r *http.Request) {
	if !isTourist(r) {
		http.Error(w, "forbidden: only tourists can get simulated position", http.StatusForbidden)
		return
	}

	var req struct {
		ExecID string `json:"exec_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	execID, err := uuid.Parse(req.ExecID)
	if err != nil {
		http.Error(w, "invalid execID", http.StatusBadRequest)
		return
	}

	lat, lng, err := h.service.SimulatedPosition(execID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{
		"latitude":  lat,
		"longitude": lng,
	})
}
