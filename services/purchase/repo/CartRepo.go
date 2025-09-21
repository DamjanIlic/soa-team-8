package repo

import (
	"purchase/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func (r *CartRepository) Create(cart *model.ShoppingCart) error {
	return r.DB.Create(cart).Error
}

func (r *CartRepository) GetByID(id uuid.UUID) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	if err := r.DB.Preload("Items").First(&cart, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) Update(cart *model.ShoppingCart) error {
	return r.DB.Save(cart).Error
}

func (r *CartRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.ShoppingCart{}, "id = ?", id).Error
}

func (r *CartRepository) GetByTouristID(touristID uuid.UUID) (*model.ShoppingCart, error) {
	var cart model.ShoppingCart
	if err := r.DB.Preload("Items").
		Where("tourist_id = ?", touristID).
		First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}
