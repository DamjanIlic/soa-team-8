package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Stakeholder struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name" gorm:"not null;type:string"`
	Surname string    `json:"surname"`
}

func (stakeholder *Stakeholder) BeforeCreate(scope *gorm.DB) error {
	stakeholder.ID = uuid.New()
	return nil
}
