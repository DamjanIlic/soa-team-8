package service

import (
	"errors"
	"stakeholder/model"
	"stakeholder/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

func (s *UserService) GetAllUsers() []model.User {
	return s.UserRepo.FindAll()
}

func (s *UserService) RegisterUser(user *model.User) error {
	if user.Role != model.RoleTourist && user.Role != model.RoleGuide {
		return errors.New("invalid role")
	}

	if _, err := s.UserRepo.FindByEmail(user.Email); err == nil {
		return errors.New("email already exists")
	}
	if _, err := s.UserRepo.FindByUsername(user.Username); err == nil {
		return errors.New("username already exists")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)

	return s.UserRepo.Create(user)
}

func (s *UserService) BlockUser(userID string) error {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Blocked = true
	return s.UserRepo.Update(user)
}

