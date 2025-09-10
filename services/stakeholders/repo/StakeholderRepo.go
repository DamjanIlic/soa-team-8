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

func (r *StakeholderRepository) GetByUserID(userID uuid.UUID) (*model.Stakeholder, error) {
	var s model.Stakeholder
	
	if err := r.DatabaseConnection.Preload("User").First(&s, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StakeholderRepository) Update(stakeholder *model.Stakeholder) error {
	return r.DatabaseConnection.Save(stakeholder).Error
}

func (r *StakeholderRepository) CreateForUser(userID uuid.UUID, name, surname string) (*model.Stakeholder, error) {
	stakeholder := &model.Stakeholder{
		UserID:  userID,
		Name:    name,
		Surname: surname,
	}
	
	if err := r.DatabaseConnection.Create(stakeholder).Error; err != nil {
		return nil, err
	}
	return stakeholder, nil
}
