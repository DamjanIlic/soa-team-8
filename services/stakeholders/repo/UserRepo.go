package repo

import (
	"stakeholder/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	DatabaseConnection *gorm.DB
}

func (r *UserRepo) FindAll() []model.User {
	var users []model.User
	r.DatabaseConnection.Find(&users)
	return users
}
