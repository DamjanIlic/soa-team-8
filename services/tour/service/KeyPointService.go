package service

import (
	"math"
	"tour/model"
	"tour/repo"

	"github.com/google/uuid"
)

type KeyPointService struct {
	KeyPointRepo *repo.KeyPointRepository
	TourRepo     *repo.TourRepository
}

func (s *KeyPointService) CreateKeyPoint(tourID string, req *model.KeyPointRequest) (*model.KeyPointResponse, error) {
	tid, err := uuid.Parse(tourID)
	if err != nil {
		return nil, err
	}

	// Proveri da li tura postoji
	tour, err := s.TourRepo.GetByID(tourID)
	if err != nil {
		return nil, err
	}

	keyPoint := &model.KeyPoint{
		TourID:      tid,
		Name:        req.Name,
		Description: req.Description,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		ImageURL:    req.ImageURL,
		Order:       req.Order,
	}

	if err := s.KeyPointRepo.Create(keyPoint); err != nil {
		return nil, err
	}

	// Update distance ture ako postoji prethodna ključna tačka
	keyPoints, err := s.KeyPointRepo.GetByTourID(tid)
	if err == nil && len(keyPoints) > 1 {
		last := keyPoints[len(keyPoints)-2] // pređašnja tačka
		addedDistance := haversine(last.Latitude, last.Longitude, keyPoint.Latitude, keyPoint.Longitude)
		tour.DistanceKm += addedDistance
		_ = s.TourRepo.Update(tour)
	}

	response := &model.KeyPointResponse{
		ID:          keyPoint.ID.String(),
		TourID:      keyPoint.TourID.String(),
		Name:        keyPoint.Name,
		Description: keyPoint.Description,
		Latitude:    keyPoint.Latitude,
		Longitude:   keyPoint.Longitude,
		ImageURL:    keyPoint.ImageURL,
		Order:       keyPoint.Order,
		CreatedAt:   keyPoint.CreatedAt,
	}

	return response, nil
}

func (s *KeyPointService) GetKeyPointsByTour(tourID string) ([]model.KeyPointResponse, error) {
	tid, err := uuid.Parse(tourID)
	if err != nil {
		return nil, err
	}

	keyPoints, err := s.KeyPointRepo.GetByTourID(tid)
	if err != nil {
		return nil, err
	}

	var responses []model.KeyPointResponse
	for _, kp := range keyPoints {
		responses = append(responses, model.KeyPointResponse{
			ID:          kp.ID.String(),
			TourID:      kp.TourID.String(),
			Name:        kp.Name,
			Description: kp.Description,
			Latitude:    kp.Latitude,
			Longitude:   kp.Longitude,
			ImageURL:    kp.ImageURL,
			Order:       kp.Order,
			CreatedAt:   kp.CreatedAt,
		})
	}

	return responses, nil
}

func (s *KeyPointService) GetKeyPoint(id string) (*model.KeyPointResponse, error) {
	keyPoint, err := s.KeyPointRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &model.KeyPointResponse{
		ID:          keyPoint.ID.String(),
		TourID:      keyPoint.TourID.String(),
		Name:        keyPoint.Name,
		Description: keyPoint.Description,
		Latitude:    keyPoint.Latitude,
		Longitude:   keyPoint.Longitude,
		ImageURL:    keyPoint.ImageURL,
		Order:       keyPoint.Order,
		CreatedAt:   keyPoint.CreatedAt,
	}

	return response, nil
}

func (s *KeyPointService) UpdateKeyPoint(id string, updates map[string]interface{}) error {
	keyPoint, err := s.KeyPointRepo.GetByID(id)
	if err != nil {
		return err
	}

	// update polja
	if name, ok := updates["name"].(string); ok {
		keyPoint.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		keyPoint.Description = description
	}
	if latitude, ok := updates["latitude"].(float64); ok {
		keyPoint.Latitude = latitude
	}
	if longitude, ok := updates["longitude"].(float64); ok {
		keyPoint.Longitude = longitude
	}
	if imageURL, ok := updates["image_url"].(string); ok {
		keyPoint.ImageURL = &imageURL
	}
	if order, ok := updates["order"].(float64); ok {
		keyPoint.Order = int(order)
	}

	return s.KeyPointRepo.Update(keyPoint)
}

func (s *KeyPointService) DeleteKeyPoint(id string) error {
	return s.KeyPointRepo.Delete(id)
}

// Haversine formula za udaljenost između dve geografske tačke
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
