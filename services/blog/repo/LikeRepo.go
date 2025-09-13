package repo

import (
	"blog/model"

	"gorm.io/gorm"
)

type LikeRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *LikeRepository) Create(like *model.Like) error {
	return r.DatabaseConnection.Create(like).Error
}

func (r *LikeRepository) Delete(userID, blogID string) error {
	return r.DatabaseConnection.
		Where("user_id = ? AND blog_id = ?", userID, blogID).
		Delete(&model.Like{}).Error
}

func (r *LikeRepository) Exists(userID, blogID string) (bool, error) {
	var count int64
	err := r.DatabaseConnection.
		Model(&model.Like{}).
		Where("user_id = ? AND blog_id = ?", userID, blogID).
		Count(&count).Error
	return count > 0, err
}

func (r *LikeRepository) CountByBlogID(blogID string) (int64, error) {
	var count int64
	err := r.DatabaseConnection.
		Model(&model.Like{}).
		Where("blog_id = ?", blogID).
		Count(&count).Error
	return count, err
}
