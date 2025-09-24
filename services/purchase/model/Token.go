package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourPurchaseToken struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey"`
	TourID     uuid.UUID `json:"tour_id" gorm:"not null"`
	TouristID  uuid.UUID `json:"tourist_id" gorm:"not null"`
	Token      string    `json:"token" gorm:"unique;not null"`
	IsExecuted bool      `json:"is_executed" gorm:"default:false"`
	IsReviewed bool      `json:"is_reviewed" gorm:"default:false"`
	CreatedAt  time.Time `json:"created_at"`
}

type TokenResponse struct {
	ID         string    `json:"id"`
	TourID     string    `json:"tour_id"`
	TouristID  string    `json:"tourist_id"`
	Token      string    `json:"token"`
	IsExecuted bool      `json:"is_executed"`
	IsReviewed bool      `json:"is_reviewed"`
	CreatedAt  time.Time `json:"created_at"`
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

func (t *TourPurchaseToken) ToResponse() TokenResponse {
	return TokenResponse{
		ID:         t.ID.String(),
		TourID:     t.TourID.String(),
		TouristID:  t.TouristID.String(),
		Token:      t.Token,
		IsExecuted: t.IsExecuted,
		IsReviewed: t.IsReviewed,
		CreatedAt:  t.CreatedAt,
	}
}

func (TourPurchaseToken) TableName() string {
	return "tour_purchase_tokens"
}
