package repo

import (
	"purchase/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemRepository struct {
	DB *gorm.DB
}

func (r *ItemRepository) Create(item *model.OrderItem) error {
	return r.DB.Create(item).Error
}

func (r *ItemRepository) GetByID(id uuid.UUID) (*model.OrderItem, error) {
	var item model.OrderItem
	if err := r.DB.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Update(item *model.OrderItem) error {
	return r.DB.Save(item).Error
}

func (r *ItemRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&model.OrderItem{}, "id = ?", id).Error
}
