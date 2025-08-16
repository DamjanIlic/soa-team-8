package repo

import (
	"stakeholder/model"

	"gorm.io/gorm"
)

type StakeholderRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *StakeholderRepository) Create(stakeholder *model.Stakeholder) error {
	return r.DatabaseConnection.Create(stakeholder).Error
}

func (r *StakeholderRepository) Get(id string) (*model.Stakeholder, error) {
	var stakeholder model.Stakeholder
	result := r.DatabaseConnection.First(&stakeholder, "id = ?", id)
	return &stakeholder, result.Error
}
