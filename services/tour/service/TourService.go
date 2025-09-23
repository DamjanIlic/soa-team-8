package service

import (
	"fmt"
	"time"
	"tour/model"
	"tour/repo"

	"github.com/google/uuid"
)

type TourService struct {
	TourRepo *repo.TourRepository
}

// --- CREATE & GET ---
func (s *TourService) CreateTour(authorID string, req *model.TourRequest) (*model.TourResponse, error) {
	uid, err := uuid.Parse(authorID)
	if err != nil {
		return nil, err
	}

	tour := &model.Tour{
		AuthorID:    uid,
		Name:        req.Name,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
		Status:      model.StatusDraft,
		Price:       0,
	}

	if err := s.TourRepo.Create(tour); err != nil {
		return nil, err
	}

	return s.toTourResponse(tour), nil
}

func (s *TourService) GetTour(id string) (*model.TourResponse, error) {
	tour, err := s.TourRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return s.toTourResponse(tour), nil
}

func (s *TourService) GetToursByAuthor(authorID string) ([]model.TourResponse, error) {
	uid, err := uuid.Parse(authorID)
	if err != nil {
		return nil, err
	}

	tours, err := s.TourRepo.GetByAuthorID(uid)
	if err != nil {
		return nil, err
	}

	return s.toTourResponses(tours), nil
}

func (s *TourService) GetAllTours() ([]model.TourResponse, error) {
	tours, err := s.TourRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.toTourResponses(tours), nil
}

// --- TOUR STATUS ---
func (s *TourService) PublishTour(tourID, authorID string) (*model.TourResponse, error) {
	tour, err := s.getTourAndCheckAuthor(tourID, authorID)
	if err != nil {
		return nil, err
	}

	// Preload keypoints i durations
	if err := s.TourRepo.PreloadCheckpointsAndDurations(tour); err != nil {
		return nil, err
	}

	// 1. Osnovna polja
	if tour.Name == "" || tour.Description == "" || tour.Difficulty == "" || tour.Tags == "" {
		return nil, fmt.Errorf("tour missing required fields")
	}

	// 2. Cena
	if tour.Price <= 0 {
		return nil, fmt.Errorf("tour price must be greater than 0")
	}

	// 3. Minimum 2 ključne tačke
	if len(tour.Checkpoints) < 2 {
		return nil, fmt.Errorf("tour must have at least 2 key points")
	}

	// 4. Bar jedno trajanje po prevozu
	if len(tour.Durations) == 0 {
		return nil, fmt.Errorf("tour must have at least one duration for a transport type")
	}

	now := time.Now()
	tour.Status = model.StatusPublished
	tour.PublishedAt = &now

	if err := s.TourRepo.Update(tour); err != nil {
		return nil, err
	}

	return s.toTourResponse(tour), nil
}

func (s *TourService) ArchiveTour(tourID, authorID string) (*model.TourResponse, error) {
	tour, err := s.getTourAndCheckAuthor(tourID, authorID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	tour.Status = model.StatusArchived
	tour.ArchivedAt = &now

	if err := s.TourRepo.Update(tour); err != nil {
		return nil, err
	}

	return s.toTourResponse(tour), nil
}

func (s *TourService) ReactivateTour(tourID, authorID string) (*model.TourResponse, error) {
	tour, err := s.getTourAndCheckAuthor(tourID, authorID)
	if err != nil {
		return nil, err
	}

	tour.Status = model.StatusDraft
	tour.ArchivedAt = nil

	if err := s.TourRepo.Update(tour); err != nil {
		return nil, err
	}

	return s.toTourResponse(tour), nil
}

// --- HELPERS ---
func (s *TourService) toTourResponse(tour *model.Tour) *model.TourResponse {
	// Mapiranje keypoints
	keypoints := make([]model.KeyPointResponse, len(tour.Checkpoints))
	for i, kp := range tour.Checkpoints {
		keypoints[i] = model.KeyPointResponse{
			ID:          kp.ID.String(),
			TourID:      kp.TourID.String(),
			Name:        kp.Name,
			Description: kp.Description,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			ImageURL:    kp.ImageURL,
			Order:       kp.Order,
			CreatedAt:   kp.CreatedAt,
		}
	}

	// Mapiranje durations
	durations := make([]model.DurationResponse, len(tour.Durations))
	for i, d := range tour.Durations {
		durations[i] = model.DurationResponse{
			ID:            d.ID.String(),
			TourID:        d.TourID.String(),
			TransportType: string(d.TransportType),
			Minutes:       d.Minutes,
			CreatedAt:     d.CreatedAt,
		}
	}

	return &model.TourResponse{
		ID:          tour.ID.String(),
		AuthorID:    tour.AuthorID.String(),
		Name:        tour.Name,
		Description: tour.Description,
		Difficulty:  tour.Difficulty,
		Tags:        tour.Tags,
		Status:      string(tour.Status),
		Price:       tour.Price,
		DistanceKm:  tour.DistanceKm,
		Checkpoints: keypoints,
		Durations:   durations,
		CreatedAt:   tour.CreatedAt,
	}
}

func (s *TourService) toTourResponses(tours []model.Tour) []model.TourResponse {
	responses := make([]model.TourResponse, len(tours))
	for i, t := range tours {
		responses[i] = *s.toTourResponse(&t)
	}
	return responses
}

func (s *TourService) getTourAndCheckAuthor(tourID, authorID string) (*model.Tour, error) {
	tour, err := s.TourRepo.GetByID(tourID)
	if err != nil {
		return nil, err
	}
	if tour.AuthorID.String() != authorID {
		return nil, fmt.Errorf("forbidden: you are not the author")
	}
	return tour, nil
}

func (s *TourService) UpdateDistance(tourID string, distance float64) (*model.TourResponse, error) {
	tour, err := s.TourRepo.GetByID(tourID)
	if err != nil {
		return nil, err
	}
	tour.DistanceKm = distance
	if err := s.TourRepo.Update(tour); err != nil {
		return nil, err
	}
	return s.toTourResponse(tour), nil
}

func (s *TourService) UpdatePrice(tourID string, price float64) (*model.TourResponse, error) {
	tour, err := s.TourRepo.GetByID(tourID)
	if err != nil {
		return nil, err
	}
	tour.Price = price
	if err := s.TourRepo.Update(tour); err != nil {
		return nil, err
	}
	return s.toTourResponse(tour), nil
}
