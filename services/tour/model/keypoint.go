package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeyPoint struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey"`
	TourID      uuid.UUID `json:"tour_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude" gorm:"not null"`
	Longitude   float64   `json:"longitude" gorm:"not null"`
	ImageURL    *string   `json:"image_url,omitempty"`
	Order       int       `json:"order" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Tour Tour `json:"-" gorm:"foreignKey:TourID"`
}

type KeyPointRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ImageURL    *string `json:"image_url,omitempty"`
}

type KeyPointResponse struct {
	ID          string    `json:"id"`
	TourID      string    `json:"tour_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	ImageURL    *string   `json:"image_url,omitempty"`
	Order       int       `json:"order"`
	CreatedAt   time.Time `json:"created_at"`
}

func (keyPoint *KeyPoint) BeforeCreate(tx *gorm.DB) error {
	if keyPoint.ID == uuid.Nil {
		keyPoint.ID = uuid.New()
	}
	if keyPoint.Latitude < -90 || keyPoint.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if keyPoint.Longitude < -180 || keyPoint.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

func (KeyPoint) TableName() string {
	return "key_points"
}
