package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourPurchaseToken struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	TourID    uuid.UUID `json:"tour_id" gorm:"not null"`
	TouristID uuid.UUID `json:"tourist_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
}

func (t *TourPurchaseToken) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.Token == "" {
		t.Token = uuid.NewString()
	}
	return nil
}

func (TourPurchaseToken) TableName() string {
	return "tour_purchase_tokens"
}
