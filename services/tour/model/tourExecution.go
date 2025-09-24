package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourExecutionStatus string

const (
	StatusActive    TourExecutionStatus = "active"
	StatusCompleted TourExecutionStatus = "completed"
	StatusAbandoned TourExecutionStatus = "abandoned"
)

type KeyPointExecution struct {
	ID              uuid.UUID  `json:"id" gorm:"primaryKey"`
	TourExecutionID uuid.UUID  `json:"tour_execution_id" gorm:"not null"`
	KeyPointID      uuid.UUID  `json:"key_point_id" gorm:"not null"`
	ReachedAt       *time.Time `json:"reached_at"`
}

type TourExecution struct {
	ID         uuid.UUID           `json:"id" gorm:"primaryKey"`
	UserID     uuid.UUID           `json:"user_id" gorm:"not null"`
	TourID     uuid.UUID           `json:"tour_id" gorm:"not null"`
	Status     TourExecutionStatus `json:"status" gorm:"default:active"`
	StartedAt  time.Time           `json:"started_at"`
	EndedAt    *time.Time          `json:"ended_at"`
	LastActive time.Time           `json:"last_active"`

	KeyPoints []KeyPointExecution `json:"key_points" gorm:"foreignKey:TourExecutionID"`
}

func (exec *TourExecution) BeforeCreate(tx *gorm.DB) error {
	if exec.ID == uuid.Nil {
		exec.ID = uuid.New()
	}
	exec.StartedAt = time.Now()
	exec.LastActive = time.Now()
	return nil
}
func (kpe *KeyPointExecution) BeforeCreate(tx *gorm.DB) error {
	if kpe.ID == uuid.Nil {
		kpe.ID = uuid.New()
	}
	return nil
}

func (TourExecution) TableName() string {
	return "tour_executions"
}

func (KeyPointExecution) TableName() string {
	return "key_point_executions"
}
