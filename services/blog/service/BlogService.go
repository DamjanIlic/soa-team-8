package service

import (
	"blog/model"
	"blog/repo"
)

type BlogService struct {
	BlogRepo *repo.BlogRepository
}

func (s *BlogService) Create(blog *model.Blog) error {
	return s.BlogRepo.Create(blog)
}

func (s *BlogService) Get(id string) (*model.Blog, error) {
	return s.BlogRepo.Get(id)
}

func (s *BlogService) Like(id string) error {
	blog, err := s.BlogRepo.Get(id)
	if err != nil {
		return err
	}
	blog.Likes++
	return s.BlogRepo.Update(blog)
}

func (s *BlogService) Unlike(id string) error {
	blog, err := s.BlogRepo.Get(id)
	if err != nil {
		return err
	}
	if blog.Likes > 0 {
		blog.Likes--
	}
	return s.BlogRepo.Update(blog)
}
