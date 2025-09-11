package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type CommentWithUser struct {
//     Comment
//     User User `json:"user"`
// }

type Comment struct {
	ID        uuid.UUID `json:"id"`
	BlogID    uuid.UUID `json:"blog_id"`
	UserID    uuid.UUID `json:"user_id"`
	Text      string    `json:"text" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook za automatsko generisanje UUID-a
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
