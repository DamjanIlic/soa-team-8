package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"purchase/model"
	"purchase/repo"

	"github.com/google/uuid"
)

type CartService struct {
	CartRepo *repo.CartRepository
	ItemRepo *repo.ItemRepository
}

// kreira novu praznu korpu za turistu
func (s *CartService) CreateCart(touristID uuid.UUID) (*model.ShoppingCart, error) {
	cart := &model.ShoppingCart{
		TouristID: touristID,
	}
	if err := s.CartRepo.Create(cart); err != nil {
		return nil, err
	}
	return cart, nil
}

// dodaje item u korpu (auto-kreira korpu ako ne postoji)
func (s *CartService) AddItemToUserCart(touristID uuid.UUID, tourID uuid.UUID, name string, price float64) (*model.OrderItem, error) {
	tourStatus, err := s.checkTourStatus(tourID)
	if err != nil {
		return nil, errors.New("failed to verify tour status")
	}
	if tourStatus != "available" {
		return nil, errors.New("tour not available for purchase")
	}

	cart, err := s.CartRepo.GetByTouristID(touristID)
	if err != nil {
		// ako korpa ne postoji, kreiraj novu
		cart, err = s.CreateCart(touristID)
		if err != nil {
			return nil, errors.New("failed to create cart: " + err.Error())
		}
	}

	// sada dodaj item u korpu
	item := &model.OrderItem{
		CartID: cart.ID,
		TourID: tourID,
		Name:   name,
		Price:  price,
	}
	if err := s.ItemRepo.Create(item); err != nil {
		return nil, err
	}

	// osvezi korpu
	cart.Items = append(cart.Items, *item)
	if err := s.CartRepo.Update(cart); err != nil {
		return nil, err
	}

	if err := s.updateCartTotal(cart); err != nil {
		return nil, err
	}

	return item, nil
}

// uklanja item iz korpe
func (s *CartService) RemoveItem(cartID, itemID uuid.UUID) error {
	cart, err := s.CartRepo.GetByID(cartID)
	if err != nil {
		return err
	}

	found := false
	for i, item := range cart.Items {
		if item.ID == itemID {
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return errors.New("item not found in cart")
	}

	if err := s.ItemRepo.Delete(itemID); err != nil {
		return err
	}

	return s.updateCartTotal(cart)
}

// vraca ukupnu cenu
func (s *CartService) GetTotal(cartID uuid.UUID) (float64, error) {
	cart, err := s.CartRepo.GetByID(cartID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range cart.Items {
		total += item.Price
	}
	return total, nil
}

func (s *CartService) GetByTouristID(touristID uuid.UUID) (*model.ShoppingCart, error) {
	return s.CartRepo.GetByTouristID(touristID)
}

func (s *CartService) updateCartTotal(cart *model.ShoppingCart) error {
	var total float64
	for _, item := range cart.Items {
		total += item.Price
	}
	cart.Total = total
	return s.CartRepo.Update(cart)
}

func (s *CartService) checkTourStatus(tourID uuid.UUID) (string, error) {
	url := fmt.Sprintf("http://tour-service:8080/api/tours/%s/status", tourID)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
