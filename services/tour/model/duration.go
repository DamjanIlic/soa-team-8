package model

import (
	"time"

	"github.com/google/uuid"
)

type TransportType string

const (
	TransportWalk TransportType = "walk"
	TransportBike TransportType = "bike"
	TransportCar  TransportType = "car"
)

type Duration struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	TourID        uuid.UUID     `json:"tour_id"`
	TransportType TransportType `json:"transport_type"`
	Minutes       int           `json:"minutes"`
	CreatedAt     time.Time     `json:"created_at"`
}
