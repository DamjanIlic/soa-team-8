package service

import (
	"blog/model"
	"blog/repo"
	"fmt"
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

// Dohvatanje svih blogova sa ažuriranim brojem lajkova
func (s *BlogService) GetAll() ([]model.Blog, error) {
	blogs, err := s.BlogRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range blogs {
		count, _ := s.LikeRepo.CountByBlogID(blogs[i].ID)
		blogs[i].Likes = int(count)
	}

	return blogs, nil
}

// Dohvatanje bloga po ID-u sa ažuriranim brojem lajkova
func (s *BlogService) Get(id string) (*model.Blog, error) {
	blog, err := s.BlogRepo.Get(id)
	if err != nil {
		return nil, err
	}

	count, _ := s.LikeRepo.CountByBlogID(id)
	blog.Likes = int(count)

	return blog, nil
}

func (s *BlogService) Like(blogID, userID string) (int, error) {
	// Proveri da li blog postoji
	_, err := s.BlogRepo.Get(blogID)
	if err != nil {
		return 0, fmt.Errorf("blog not found")
	}

	// Proveri da li korisnik već lajkovao
	exists, err := s.LikeRepo.Exists(userID, blogID)
	if err != nil {
		return 0, err
	}
	if exists {
		count, _ := s.LikeRepo.CountByBlogID(blogID)
		return int(count), nil
	}

	like := model.NewLike(userID, blogID)
	if err := s.LikeRepo.Create(like); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}

func (s *BlogService) Unlike(blogID, userID string) (int, error) {
	// Proveri da li blog postoji
	if _, err := s.BlogRepo.Get(blogID); err != nil {
		return 0, fmt.Errorf("blog not found")
	}

	if err := s.LikeRepo.Delete(userID, blogID); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(blogID)
	return int(count), nil
}
