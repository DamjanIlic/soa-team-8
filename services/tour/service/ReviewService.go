package service

import (
	"fmt"
	"tour/model"
	"tour/repo"

	"github.com/google/uuid"
)

type ReviewService struct {
	ReviewRepo *repo.ReviewRepository
	TourRepo   *repo.TourRepository
}

func (s *ReviewService) CreateReview(tourID, touristID string, req *model.ReviewRequest) (*model.ReviewResponse, error) {
	review, err := model.FromRequest(tourID, touristID, req)
	if err != nil {
		return nil, err
	}

	if err := s.ReviewRepo.Create(review); err != nil {
		return nil, err
	}

	response := review.ToResponse()
	return &response, nil
}

func (s *ReviewService) GetReviewsByTour(tourID string) ([]model.ReviewResponse, error) {
	reviews, err := s.ReviewRepo.GetByTour(tourID)
	if err != nil {
		return nil, err
	}
	var responses []model.ReviewResponse
	for _, r := range reviews {
		responses = append(responses, r.ToResponse())
	}
	return responses, nil
}

func (s *ReviewService) GetAllReviews() ([]model.ReviewResponse, error) {
	reviews, err := s.ReviewRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}

	return responses, nil
}

func (s *ReviewService) GetReviewsForGuide(guideID string) ([]model.ReviewResponse, error) {
	parsedID, err := uuid.Parse(guideID)
	if err != nil {
		return nil, fmt.Errorf("invalid guide ID: %w", err)
	}

	reviews, err := s.ReviewRepo.GetReviewsForGuide(parsedID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}
	return responses, nil
}

func (s *ReviewService) GetReviewsByTourist(touristID string) ([]model.ReviewResponse, error) {
	reviews, err := s.ReviewRepo.GetByTouristID(touristID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ReviewResponse, len(reviews))
	for i, review := range reviews {
		responses[i] = review.ToResponse()
	}

	return responses, nil
}
