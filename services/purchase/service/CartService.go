package service

import (
	"errors"
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

// dodaje item u korpu
func (s *CartService) AddItem(cartID uuid.UUID, tourID uuid.UUID, name string, price float64) (*model.OrderItem, error) {
	cart, err := s.CartRepo.GetByID(cartID)
	if err != nil {
		return nil, err
	}

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
	return s.CartRepo.Update(cart)
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
