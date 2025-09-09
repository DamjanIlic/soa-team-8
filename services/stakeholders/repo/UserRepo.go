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

func (r *UserRepo) Create(user *model.User) error {
	return r.DatabaseConnection.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.DatabaseConnection.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.DatabaseConnection.Where("username = ?", username).First(&user)
	return &user, result.Error
}
