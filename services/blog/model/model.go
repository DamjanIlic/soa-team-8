package model

import (
	"time"

	"github.com/google/uuid"
)

type Blog struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	ImageURL  string    `json:"image_url,omitempty" bson:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Likes     int       `json:"likes" bson:"likes"`
}

func NewBlog(title, content, imageURL string) *Blog {
	return &Blog{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
	}
}

type Like struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id"`
	BlogID    string    `json:"blog_id" bson:"blog_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

func NewLike(userID, blogID string) *Like {
	return &Like{
		ID:        uuid.New().String(),
		UserID:    userID,
		BlogID:    blogID,
		CreatedAt: time.Now(),
	}
}
