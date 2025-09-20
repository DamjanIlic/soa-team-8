package handler

import (
	"encoding/json"
	"net/http"
	"purchase/middleware"
	"purchase/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	CartService *service.CartService
}

func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value(middleware.ContextRole).(string)
	if !ok || role != "tourist" {
		http.Error(w, "Only tourists can create carts", http.StatusForbidden)
		return
	}

	userID := r.Context().Value(middleware.ContextUserID)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	touristID, err := uuid.Parse(userID.(string))
	if err != nil {
		http.Error(w, "Invalid tourist ID", http.StatusBadRequest)
		return
	}

	cart, err := h.CartService.CreateCart(touristID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	cartIDStr := mux.Vars(r)["cartId"]
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	var req struct {
		TourID string  `json:"tour_id"`
		Name   string  `json:"name"`
		Price  float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tourID, err := uuid.Parse(req.TourID)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	item, err := h.CartService.AddItem(cartID, tourID, req.Name, req.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID, err := uuid.Parse(vars["cartId"])
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}
	itemID, err := uuid.Parse(vars["itemId"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	if err := h.CartService.RemoveItem(cartID, itemID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed successfully"})
}

func (h *CartHandler) GetTotal(w http.ResponseWriter, r *http.Request) {
	cartIDStr := mux.Vars(r)["cartId"]
	cartID, err := uuid.Parse(cartIDStr)
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	total, err := h.CartService.GetTotal(cartID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]float64{"total": total}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
