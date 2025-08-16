package service

import (
	"stakeholder/model"
	"stakeholder/repo"
)

type StakeholderService struct {
	StakeholderRepo *repo.StakeholderRepository
}

func (s *StakeholderService) Create(stakeholder *model.Stakeholder) error {
	return s.StakeholderRepo.Create(stakeholder)
}

func (s *StakeholderService) Get(id string) (*model.Stakeholder, error) {
	return s.StakeholderRepo.Get(id)
}
