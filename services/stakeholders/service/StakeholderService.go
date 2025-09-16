package service

import (
	"stakeholder/model"
	"stakeholder/repo"

	"github.com/google/uuid"
)

type StakeholderService struct {
	StakeholderRepo *repo.StakeholderRepository
}

func (s *StakeholderService) Create(stakeholder *model.Stakeholder) error {
	//
	return s.StakeholderRepo.Create(stakeholder)
}

func (s *StakeholderService) Get(id string) (*model.Stakeholder, error) {
	return s.StakeholderRepo.Get(id)
}

func (s *StakeholderService) GetProfile(userID string) (*model.ProfileResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	stakeholder, err := s.StakeholderRepo.GetByUserID(uid)
	if err != nil {
		return nil, err
	}

	profile := &model.ProfileResponse{
		ID:           stakeholder.ID.String(),
		Username:     stakeholder.User.Username,
		Email:        stakeholder.User.Email,
		Role:         string(stakeholder.User.Role),
		Name:         stakeholder.Name,
		Surname:      stakeholder.Surname,
		ProfileImage: stakeholder.ProfileImage,
		Biography:    stakeholder.Biography,
		Motto:        stakeholder.Motto,
	}

	return profile, nil
}

func (s *StakeholderService) UpdateProfile(userID string, updates map[string]interface{}) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	stakeholder, err := s.StakeholderRepo.GetByUserID(uid)
	if err != nil {
		return err
	}

	// update polja
	if name, ok := updates["name"].(string); ok {
		stakeholder.Name = name
	}
	if surname, ok := updates["surname"].(string); ok {
		stakeholder.Surname = surname
	}
	if biography, ok := updates["biography"].(string); ok {
		stakeholder.Biography = &biography
	}
	if motto, ok := updates["motto"].(string); ok {
		stakeholder.Motto = &motto
	}
	if profileImage, ok := updates["profile_image"].(string); ok {
		stakeholder.ProfileImage = &profileImage
	}

	return s.StakeholderRepo.Update(stakeholder)
}
