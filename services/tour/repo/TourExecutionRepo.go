package repo

import (
	"tour/model"

	"gorm.io/gorm"
)

type TourExecutionRepository struct {
	DatabaseConnection *gorm.DB
}

func (r *TourExecutionRepository) Create(exec *model.TourExecution) error {
	return r.DatabaseConnection.Create(exec).Error
}
