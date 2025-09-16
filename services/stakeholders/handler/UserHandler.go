package handler

import (
	"encoding/json"
	"net/http"
	"stakeholder/model"
	"stakeholder/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserService *service.UserService
}

// Dohvati sve korisnike
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users := h.UserService.GetAllUsers()

	var response []model.UserResponse
	for _, u := range users {
		response = append(response, model.UserResponse{
			ID:       u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
			Blocked:  u.Blocked, // ako si dodao polje u model.UserResponse
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Registracija korisnika
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
		Role:     model.Role(req.Role),
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
		Blocked:  user.Blocked, // isto ovde ako koristiš
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Blokiranje korisnika
func (h *UserHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	if err := h.UserService.BlockUser(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User blocked successfully",
	})
}
