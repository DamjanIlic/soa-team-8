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
			Role:     string(u.Role),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req model.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     model.Role(req.Role), // cast na tip Role
	}

	if err := h.UserService.RegisterUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := model.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
