package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShoppingCart struct {
	ID        uuid.UUID   `json:"id" gorm:"primaryKey"`
	TouristID uuid.UUID   `json:"tourist_id" gorm:"not null"`
	Items     []OrderItem `json:"items" gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	Total     float64     `json:"total" gorm:"default:0"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (cart *ShoppingCart) BeforeCreate(tx *gorm.DB) error {
	if cart.ID == uuid.Nil {
		cart.ID = uuid.New()
	}
	return nil
}

func (ShoppingCart) TableName() string {
	return "shopping_carts"
}
