package handler

import (
	"encoding/json"
	"net/http"
	"tour/model"
	"tour/service"

	"github.com/gorilla/mux"
)

type KeyPointHandler struct {
	KeyPointService *service.KeyPointService
}

func (h *KeyPointHandler) CreateKeyPoint(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["tourId"]

	var req model.KeyPointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keyPoint, err := h.KeyPointService.CreateKeyPoint(tourID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(keyPoint)
}

func (h *KeyPointHandler) GetKeyPointsByTour(w http.ResponseWriter, r *http.Request) {
	tourID := mux.Vars(r)["tourId"]

	keyPoints, err := h.KeyPointService.GetKeyPointsByTour(tourID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyPoints)
}

func (h *KeyPointHandler) GetKeyPoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	keyPoint, err := h.KeyPointService.GetKeyPoint(id)
	if err != nil {
		http.Error(w, "KeyPoint not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyPoint)
}

func (h *KeyPointHandler) UpdateKeyPoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.KeyPointService.UpdateKeyPoint(id, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "KeyPoint updated successfully"})
}

func (h *KeyPointHandler) DeleteKeyPoint(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.KeyPointService.DeleteKeyPoint(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "KeyPoint deleted successfully"})
}