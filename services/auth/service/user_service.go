package service

import (
	"auth-service/model"
	"auth-service/repo"
	"errors"

	"auth-service/util"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

// Registracija korisnika
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

// Login korisnika i vraćanje JWT tokena
func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := util.GenerateJWT(user.ID.String(), string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Blokiranje korisnika
func (s *UserService) BlockUser(userID string) error {
	user, err := s.UserRepo.FindByID(userID)
	if err != nil {
		return err
	}
	user.Blocked = true
	return s.UserRepo.Update(user)
}
