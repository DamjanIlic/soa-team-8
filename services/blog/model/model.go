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
	UserID    string    `bson:"user_id"`
}
type BlogResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Likes     int       `json:"likes"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewBlog(title, content, imageURL string, userID string) *Blog {
	return &Blog{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Likes:     0,
		UserID:    userID,
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
