package repo

import (
	"tour/model"

	"gorm.io/gorm"
)

type DurationRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *DurationRepository) Create(duration *model.Duration) error {
	return r.DatabaseConnection.Create(duration).Error
}

func (r *DurationRepository) GetByTourID(tourID string) ([]model.Duration, error) {
	var durations []model.Duration
	if err := r.DatabaseConnection.Where("tour_id = ?", tourID).Find(&durations).Error; err != nil {
		return nil, err
	}
	return durations, nil
}
