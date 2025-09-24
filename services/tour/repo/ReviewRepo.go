package repo

import (
	"tour/model"

	"github.com/google/uuid"
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

func (r *ReviewRepository) GetAll() ([]model.Review, error) {
	var reviews []model.Review
	if err := r.DatabaseConnection.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) GetByTouristID(touristID string) ([]model.Review, error) {
	var reviews []model.Review
	if err := r.DatabaseConnection.Where("tourist_id = ?", touristID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewRepository) GetReviewsForGuide(guideID uuid.UUID) ([]model.Review, error) {
	var reviews []model.Review
	err := r.DatabaseConnection.
		Joins("JOIN tours ON reviews.tour_id = tours.id").
		Where("tours.author_id = ?", guideID). // ovde je sad uuid, nema konflikta
		Find(&reviews).Error
	return reviews, err
}
