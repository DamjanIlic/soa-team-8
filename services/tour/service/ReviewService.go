package service

import (
	"tour/model"
	"tour/repo"
)

type ReviewService struct {
	ReviewRepo *repo.ReviewRepository
}

func (s *ReviewService) CreateReview(tourID, touristID string, req *model.ReviewRequest) (*model.ReviewResponse, error) {
	review, err := model.FromRequest(tourID, touristID, req)
	if err != nil {
		return nil, err
	}
	if err := s.ReviewRepo.Create(review); err != nil {
		return nil, err
	}
	resp := review.ToResponse()
	return &resp, nil
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
