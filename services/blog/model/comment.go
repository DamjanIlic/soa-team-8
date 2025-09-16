package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Konstruktor za novi komentar
func NewComment(blogID, userID, text string) *Comment {
	return &Comment{
		ID:        uuid.New().String(),
		BlogID:    blogID,
		UserID:    userID,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
