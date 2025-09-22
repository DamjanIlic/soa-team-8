package repo

import (
	"tour/model"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *ReviewRepository) Create(review *model.Review) error {
	return r.DatabaseConnection.Create(review).Error
}

func (r *ReviewRepository) GetByTour(tourID string) ([]model.Review, error) {
	var reviews []model.Review
	err := r.DatabaseConnection.Where("tour_id = ?", tourID).Find(&reviews).Error
	return reviews, err
}
