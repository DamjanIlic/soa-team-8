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

type CommentResponse struct {
	ID        string    `json:"id"`
	BlogID    string    `json:"blog_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`    // novo
	Surname   string    `json:"surname"` // novo
	AvatarUrl string    `json:"avatar_url"`
	Motto     string    `json:"motto"` // novo
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
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
