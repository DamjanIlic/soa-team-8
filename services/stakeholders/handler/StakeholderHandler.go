package handler

import (
	"encoding/json"
	"net/http"
	"stakeholder/middleware"
	"stakeholder/model"
	"stakeholder/service"

	"github.com/gorilla/mux"
)

type StakeholderHandler struct {
	StakeholderService *service.StakeholderService
}

func (h *StakeholderHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	var input model.StakeholderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// kreiraj stakeholder uz userID i input polja
	stakeholder, err := h.StakeholderService.Create(userID, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stakeholders := h.StakeholderService.GetAll()

	var response []model.ProfileResponse
	for _, s := range stakeholders {
		response = append(response, model.ProfileResponse{
			ID:           s.ID.String(),
			Username:     s.User.Username,
			Email:        s.User.Email,
			Role:         string(s.User.Role),
			Name:         s.Name,
			Surname:      s.Surname,
			ProfileImage: s.ProfileImage,
			Biography:    s.Biography,
			Motto:        s.Motto,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *StakeholderHandler) Get(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	id := mux.Vars(r)["id"]
	stakeholder, _ := h.StakeholderService.Get(id)
	if stakeholder == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(stakeholder)
}

func (h *StakeholderHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	profile, err := h.StakeholderService.GetProfile(userID)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *StakeholderHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

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
