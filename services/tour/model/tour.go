package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourStatus string

const (
	StatusDraft     TourStatus = "draft"
	StatusPublished TourStatus = "published"
	StatusArchived  TourStatus = "archived"
)

type Tour struct {
	ID          uuid.UUID  `json:"id" gorm:"primaryKey"`
	AuthorID    uuid.UUID  `json:"author_id" gorm:"not null"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Difficulty  string     `json:"difficulty"`
	Tags        string     `json:"tags"`
	Status      TourStatus `json:"status" gorm:"default:draft"`
	Price       float64    `json:"price" gorm:"default:0"`
	DistanceKm  float64    `json:"distance_km" gorm:"default:0"` // nova polja
	Durations   []Duration `json:"durations" gorm:"foreignKey:TourID"`
	PublishedAt *time.Time `json:"published_at"`
	ArchivedAt  *time.Time `json:"archived_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TourRequest struct {
	AuthorID    string `json:"author_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
}

type TourResponse struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Difficulty  string    `json:"difficulty"`
	Tags        string    `json:"tags"`
	Status      string    `json:"status"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func (tour *Tour) BeforeCreate(tx *gorm.DB) error {
	if tour.ID == uuid.Nil {
		tour.ID = uuid.New()
	}
	return nil
}

func (Tour) TableName() string {
	return "tours"
}
