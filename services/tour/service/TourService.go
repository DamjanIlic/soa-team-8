package service

import (
	"tour/model"
	"tour/repo"
	"github.com/google/uuid"
)

type TourService struct {
	TourRepo *repo.TourRepository
}

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

	response := &model.TourResponse{
		ID:          tour.ID.String(),
		AuthorID:    tour.AuthorID.String(),
		Name:        tour.Name,
		Description: tour.Description,
		Difficulty:  tour.Difficulty,
		Tags:        tour.Tags,
		Status:      string(tour.Status),
		Price:       tour.Price,
		CreatedAt:   tour.CreatedAt,
	}

	return response, nil
}

func (s *TourService) GetTour(id string) (*model.TourResponse, error) {
	tour, err := s.TourRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &model.TourResponse{
		ID:          tour.ID.String(),
		AuthorID:    tour.AuthorID.String(),
		Name:        tour.Name,
		Description: tour.Description,
		Difficulty:  tour.Difficulty,
		Tags:        tour.Tags,
		Status:      string(tour.Status),
		Price:       tour.Price,
		CreatedAt:   tour.CreatedAt,
	}

	return response, nil
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

	var responses []model.TourResponse
	for _, tour := range tours {
		responses = append(responses, model.TourResponse{
			ID:          tour.ID.String(),
			AuthorID:    tour.AuthorID.String(),
			Name:        tour.Name,
			Description: tour.Description,
			Difficulty:  tour.Difficulty,
			Tags:        tour.Tags,
			Status:      string(tour.Status),
			Price:       tour.Price,
			CreatedAt:   tour.CreatedAt,
		})
	}

	return responses, nil
}

func (s *TourService) GetAllTours() ([]model.TourResponse, error) {
	tours, err := s.TourRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var responses []model.TourResponse
	for _, tour := range tours {
		responses = append(responses, model.TourResponse{
			ID:          tour.ID.String(),
			AuthorID:    tour.AuthorID.String(),
			Name:        tour.Name,
			Description: tour.Description,
			Difficulty:  tour.Difficulty,
			Tags:        tour.Tags,
			Status:      string(tour.Status),
			Price:       tour.Price,
			CreatedAt:   tour.CreatedAt,
		})
	}

	return responses, nil
}