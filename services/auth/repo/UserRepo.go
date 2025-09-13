package repo

import (
	"auth-service/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepo) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepo) Update(user *model.User) error {
	return r.DB.Save(user).Error
}
