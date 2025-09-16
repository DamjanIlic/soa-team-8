package model

import (
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
	Order       int     `json:"order"`
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
	return nil
}

func (KeyPoint) TableName() string {
	return "key_points"
}