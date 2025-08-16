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
