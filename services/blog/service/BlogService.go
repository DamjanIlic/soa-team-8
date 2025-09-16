package service

import (
	"blog/model"
	"blog/repo"
	"time"

	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepo *repo.BlogRepository
	LikeRepo *repo.LikeRepository
}

// Kreiranje novog bloga
func (s *BlogService) Create(blog *model.Blog) error {
	blog.ID = uuid.New().String()
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.Likes = 0
	return s.BlogRepo.Create(blog)
}

// Dohvatanje svih blogova
func (s *BlogService) GetAll() ([]model.Blog, error) {
	return s.BlogRepo.GetAll()
}

// Dohvatanje bloga po ID-u
func (s *BlogService) Get(id string) (*model.Blog, error) {
	return s.BlogRepo.Get(id)
}

// Lajkovanje bloga
func (s *BlogService) Like(blogID, userID string) (int, error) {
	exists, err := s.LikeRepo.Exists(userID, blogID)
	if err != nil {
		return 0, err
	}
	if exists {
		count, _ := s.LikeRepo.CountByBlogID(blogID)
		return int(count), nil
	}

	like := &model.Like{
		ID:        uuid.New().String(),
		UserID:    userID,
		BlogID:    blogID,
		CreatedAt: time.Now(),
	}

	if err := s.LikeRepo.Create(like); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}

// Uklanjanje lajka
func (s *BlogService) Unlike(blogID, userID string) (int, error) {
	if err := s.LikeRepo.Delete(userID, blogID); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}
