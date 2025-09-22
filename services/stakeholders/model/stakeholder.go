package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stakeholder struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID       uuid.UUID `json:"user_id" gorm:"not null;uniqueIndex"`
	Name         string    `json:"name" gorm:"not null;type:string"`
	Surname      string    `json:"surname"`
	ProfileImage *string   `json:"profile_image,omitempty"`
	Biography    *string   `json:"biography,omitempty"`
	Motto        *string   `json:"motto,omitempty"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

type StakeholderInput struct {
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	Biography    *string `json:"biography,omitempty"`
	Motto        *string `json:"motto,omitempty"`
	ProfileImage *string `json:"profile_image,omitempty"`
}

type ProfileResponse struct {
	ID           string  `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	ProfileImage *string `json:"profile_image,omitempty"`
	Biography    *string `json:"biography,omitempty"`
	Motto        *string `json:"motto,omitempty"`
}

func (stakeholder *Stakeholder) BeforeCreate(scope *gorm.DB) error {
	if stakeholder.ID == uuid.Nil {
		stakeholder.ID = uuid.New()
	}
	return nil
}
