package repo

import (
	"tour/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *TourRepository) Create(tour *model.Tour) error {
	return r.DatabaseConnection.Create(tour).Error
}

func (r *TourRepository) GetByID(id string) (*model.Tour, error) {
	var tour model.Tour
	
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if err := r.DatabaseConnection.First(&tour, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

func (r *TourRepository) GetByAuthorID(authorID uuid.UUID) ([]model.Tour, error) {
	var tours []model.Tour
	
	if err := r.DatabaseConnection.Where("author_id = ?", authorID).Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *TourRepository) GetAll() ([]model.Tour, error) {
	var tours []model.Tour
	
	if err := r.DatabaseConnection.Find(&tours).Error; err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *TourRepository) Update(tour *model.Tour) error {
	return r.DatabaseConnection.Save(tour).Error
}

func (r *TourRepository) Delete(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	
	return r.DatabaseConnection.Delete(&model.Tour{}, "id = ?", uid).Error
}