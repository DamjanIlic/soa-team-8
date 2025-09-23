package service

import (
	"errors"
	"time"

	"tour/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourExecutionService struct {
	db *gorm.DB
}

func NewTourExecutionService(db *gorm.DB) *TourExecutionService {
	return &TourExecutionService{db: db}
}

func (s *TourExecutionService) StartTour(userID, tourID uuid.UUID) (*model.TourExecution, error) {
	// 1. Kreiraj TourExecution
	exec := &model.TourExecution{
		UserID: userID,
		TourID: tourID,
		Status: model.StatusActive,
	}

	if err := s.db.Create(exec).Error; err != nil {
		return nil, err
	}

	// 2. Dohvati sve KeyPoints za turu
	var keyPoints []model.KeyPoint
	if err := s.db.Where("tour_id = ?", tourID).Find(&keyPoints).Error; err != nil {
		return nil, err
	}

	// 3. Kreiraj KeyPointExecution za svaki KeyPoint
	for _, kp := range keyPoints {
		kpe := &model.KeyPointExecution{
			TourExecutionID: exec.ID,
			KeyPointID:      kp.ID,
		}
		if err := s.db.Create(kpe).Error; err != nil {
			return nil, err
		}
	}

	// 4. Ponovo učitaj TourExecution sa KeyPoints
	if err := s.db.Preload("KeyPoints").First(&exec, "id = ?", exec.ID).Error; err != nil {
		return nil, err
	}

	return exec, nil
}

func (s *TourExecutionService) CompleteTour(execID uuid.UUID) error {
	var exec model.TourExecution
	if err := s.db.First(&exec, "id = ?", execID).Error; err != nil {
		return err
	}
	if exec.Status != model.StatusActive {
		return errors.New("tour not active")
	}
	now := time.Now()
	exec.Status = model.StatusCompleted
	exec.EndedAt = &now
	return s.db.Save(&exec).Error
}

func (s *TourExecutionService) AbandonTour(execID uuid.UUID) error {
	var exec model.TourExecution
	if err := s.db.First(&exec, "id = ?", execID).Error; err != nil {
		return err
	}
	if exec.Status != model.StatusActive {
		return errors.New("tour not active")
	}
	now := time.Now()
	exec.Status = model.StatusAbandoned
	exec.EndedAt = &now
	return s.db.Save(&exec).Error
}

func (s *TourExecutionService) CheckKeyPoint(execID, keyPointID uuid.UUID, lat, lng float64) (bool, error) {
	var kp model.KeyPoint
	if err := s.db.First(&kp, "id = ?", keyPointID).Error; err != nil {
		return false, err
	}

	// Provera da li je korisnik unutar 0.0001° (~10m)
	if (lat < kp.Latitude+0.0001 && lat > kp.Latitude-0.0001) &&
		(lng < kp.Longitude+0.0001 && lng > kp.Longitude-0.0001) {

		// Update postojećeg KeyPointExecution
		now := time.Now()
		if err := s.db.Model(&model.KeyPointExecution{}).
			Where("tour_execution_id = ? AND key_point_id = ?", execID, keyPointID).
			Update("reached_at", now).Error; err != nil {
			return false, err
		}

		// Update last_activity
		s.db.Model(&model.TourExecution{}).
			Where("id = ?", execID).
			Update("last_active", now)

		return true, nil
	}

	// Update last activity i ako nije stigao
	s.db.Model(&model.TourExecution{}).
		Where("id = ?", execID).
		Update("last_active", time.Now())

	return false, nil
}

func ptr[T any](v T) *T {
	return &v
}

func (s *TourExecutionService) SimulatedPosition(execID uuid.UUID) (float64, float64, error) {
	var exec model.TourExecution
	if err := s.db.First(&exec, "id = ?", execID).Error; err != nil {
		return 0, 0, err
	}

	startLat := 44.8176
	startLng := 20.4569

	// Koliko sekundi je prošlo od starta ture
	secondsPassed := time.Since(exec.StartedAt).Seconds()

	// Male promene po sekundi da pozicija stalno raste
	lat := startLat + secondsPassed*0.00003
	lng := startLng + secondsPassed*0.00003

	// Update last activity
	s.db.Model(&model.TourExecution{}).Where("id = ?", execID).
		Update("last_activity_at", time.Now())

	return lat, lng, nil
}
