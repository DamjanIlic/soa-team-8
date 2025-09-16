package handler

import (
	"encoding/json"
	"net/http"
	"stakeholder/model"
	"stakeholder/service"

	"github.com/gorilla/mux"
)

type StakeholderHandler struct {
	StakeholderService *service.StakeholderService
}

func (h *StakeholderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var stakeholder model.Stakeholder
	json.NewDecoder(r.Body).Decode(&stakeholder)
	h.StakeholderService.Create(&stakeholder)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	stakeholder, _ := h.StakeholderService.Get(id)
	if stakeholder == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	
	profile, err := h.StakeholderService.GetProfile(userID)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *StakeholderHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if err := h.StakeholderService.UpdateProfile(userID, updates); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}
