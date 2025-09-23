package handler

import (
	"encoding/json"
	"net/http"
	"purchase/middleware"
	"purchase/service"

	"github.com/google/uuid"
)

type TokenHandler struct {
	TokenService *service.TokenService
}

// Checkout – sve stavke iz korpe postaju kupljene ture sa tokenima
func (h *TokenHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	touristID, err := uuid.Parse(userID.(string))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	tokens, err := h.TokenService.Checkout(touristID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func (h *TokenHandler) GetPurchasedTours(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	touristID, err := uuid.Parse(userID.(string))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	tokens, err := h.TokenService.GetTokensForTourist(touristID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
