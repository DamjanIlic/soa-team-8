package handler

import (
	"encoding/json"
	"net/http"
	"purchase/middleware"
	"purchase/service"

	"github.com/gorilla/mux"

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

// purchase/handler/token.go
func (h *TokenHandler) MarkAsExecuted(w http.ResponseWriter, r *http.Request) {
	tokenID := mux.Vars(r)["tokenId"]

	id, err := uuid.Parse(tokenID)
	if err != nil {
		http.Error(w, "Invalid token ID", http.StatusBadRequest)
		return
	}

	if err := h.TokenService.MarkAsExecuted(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token marked as executed"})
}

func (h *TokenHandler) MarkAsReviewed(w http.ResponseWriter, r *http.Request) {
	tokenID := mux.Vars(r)["tokenId"]

	id, err := uuid.Parse(tokenID)
	if err != nil {
		http.Error(w, "Invalid token ID", http.StatusBadRequest)
		return
	}

	if err := h.TokenService.MarkAsReviewed(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token marked as reviewed"})
}

func (h *TokenHandler) GetTokenByTouristAndTour(w http.ResponseWriter, r *http.Request) {
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

	tourID := mux.Vars(r)["tourId"]
	tourUUID, err := uuid.Parse(tourID)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	token, err := h.TokenService.GetByTouristAndTour(touristID, tourUUID)
	if err != nil {
		http.Error(w, "Token not found", http.StatusNotFound)
		return
	}

	response := token.ToResponse()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
