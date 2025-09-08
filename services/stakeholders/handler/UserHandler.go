package handler

import (
	"encoding/json"
	"net/http"
	"stakeholder/model"
	"stakeholder/service"
)

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()

	var response []model.UserResponse
	for _, u := range users {
		response = append(response, model.UserResponse{
			ID:       u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
