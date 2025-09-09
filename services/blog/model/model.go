package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content"` // markdown na frontu
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int       `json:"likes"`
}

func (blog *Blog) BeforeCreate(tx *gorm.DB) (err error) {
	blog.ID = uuid.New()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Likes = 0
	return
}
