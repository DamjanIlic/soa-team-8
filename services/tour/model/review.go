package model

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	TourID    uuid.UUID `json:"tour_id" gorm:"type:uuid;not null"`
	TouristID uuid.UUID `json:"tourist_id" gorm:"type:uuid;not null"`
	Rating    int       `json:"rating" gorm:"not null"`
	Comment   string    `json:"comment"`
	VisitedAt time.Time `json:"visited_at"`
	CreatedAt time.Time `json:"created_at"`
	Images    string    `json:"images"` // csv format: "slika1.jpg,slika2.jpg"
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now()
	}
	return nil
}

func (Review) TableName() string {
	return "reviews"
}

type ReviewRequest struct {
	Rating    int      `json:"rating"`
	Comment   string   `json:"comment"`
	VisitedAt string   `json:"visited_at"` // prima kao string
	Images    []string `json:"images"`
}

type ReviewResponse struct {
	ID        string    `json:"id"`
	TourID    string    `json:"tour_id"`
	TouristID string    `json:"tourist_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	VisitedAt time.Time `json:"visited_at"`
	CreatedAt time.Time `json:"created_at"` // datum komentara
	Images    []string  `json:"images"`
}

func (r *Review) ToResponse() ReviewResponse {
	images := []string{}
	if r.Images != "" {
		images = strings.Split(r.Images, ",")
	}
	return ReviewResponse{
		ID:        r.ID.String(),
		TourID:    r.TourID.String(),
		TouristID: r.TouristID.String(),
		Rating:    r.Rating,
		Comment:   r.Comment,
		VisitedAt: r.VisitedAt,
		CreatedAt: r.CreatedAt,
		Images:    images,
	}
}

func FromRequest(tourID, touristID string, req *ReviewRequest) (*Review, error) {
	visitedAt, err := time.Parse("2006-01-02", req.VisitedAt)
	if err != nil {
		return nil, err
	}
	return &Review{
		TourID:    uuid.MustParse(tourID),
		TouristID: uuid.MustParse(touristID),
		Rating:    req.Rating,
		Comment:   req.Comment,
		VisitedAt: visitedAt,
		Images:    strings.Join(req.Images, ","),
	}, nil
}
