package repo

import (
	"blog/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

// Kreiranje novog komentara
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.DB.Create(comment).Error
}

// Dohvatanje svih komentara za jedan blog post
func (r *CommentRepository) GetByBlogID(blogID string) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.DB.Where("blog_id = ?", blogID).Order("created_at asc").Find(&comments).Error
	return comments, err
}
