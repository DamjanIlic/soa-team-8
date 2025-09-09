package repo

import (
	"stakeholder/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StakeholderRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *StakeholderRepository) Create(stakeholder *model.Stakeholder) error {
	return r.DatabaseConnection.Create(stakeholder).Error
}

func (r *StakeholderRepository) Get(id string) (*model.Stakeholder, error) {
	var s model.Stakeholder

	// konvertuj string u UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	// GORM query
	if err := r.DatabaseConnection.First(&s, "id = ?", uid).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
