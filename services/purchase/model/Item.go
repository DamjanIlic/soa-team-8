package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	CartID    uuid.UUID `json:"cart_id"`
	TourID    uuid.UUID `json:"tour_id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemRequest struct {
	TourID string  `json:"tour_id"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}

func (item *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return nil
}

func (OrderItem) TableName() string {
	return "order_items"
}
