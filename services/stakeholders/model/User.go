package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role string

const (
	RoleTourist Role = "tourist"
	RoleGuide   Role = "guide"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `gorm:"uniqueIndex" json:"username"`
	Email    string    `gorm:"uniqueIndex" json:"email"`
	Password string    `json:"-"`
	Role     Role      `json:"role"`
	Blocked  bool      `gorm:"default:false" json:"blocked"` 
}


type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Blocked  bool   `json:"blocked"`
}


type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
