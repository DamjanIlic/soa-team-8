package service

import (
	"tour/model"
	"tour/repo"

	"github.com/google/uuid"
)

type DurationService struct {
	DurationRepo *repo.DurationRepository
	TourRepo     *repo.TourRepository
}

func (s *DurationService) AddDuration(tourID string, transport string, minutes int) (*model.Duration, error) {
	tid, err := uuid.Parse(tourID)
	if err != nil {
		return nil, err
	}

	// Provera da li tura postoji
	_, err = s.TourRepo.GetByID(tourID)
	if err != nil {
		return nil, err
	}

	duration := &model.Duration{
		TourID:        tid,
		TransportType: model.TransportType(transport),
		Minutes:       minutes,
	}

	if err := s.DurationRepo.Create(duration); err != nil {
		return nil, err
	}

	return duration, nil
}

func (s *DurationService) GetDurationsByTour(tourID string) ([]model.Duration, error) {
	return s.DurationRepo.GetByTourID(tourID)
}
