package service

import (
	"stakeholder/model"
	"stakeholder/repo"
)

type UserService struct {
	UserRepo *repo.UserRepo
}

func (s *UserService) GetAllUsers() []model.User {
	return s.UserRepo.FindAll()
}
