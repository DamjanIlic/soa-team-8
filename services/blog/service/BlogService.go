package service

import (
	"blog/model"
	"blog/repo"

	"github.com/google/uuid"
)

type BlogService struct {
	BlogRepo *repo.BlogRepository
	LikeRepo *repo.LikeRepository
}


func (s *BlogService) Create(blog *model.Blog) error {
	return s.BlogRepo.Create(blog)
}


func (s *BlogService) GetAll() ([]model.Blog, error) {
	return s.BlogRepo.GetAll()
}


func (s *BlogService) Get(id string) (*model.Blog, error) {
	return s.BlogRepo.Get(id)
}

func (s *BlogService) Like(blogID, userID string) (int, error) {
	bid, err := uuid.Parse(blogID)
	if err != nil {
		return 0, err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	exists, err := s.LikeRepo.Exists(uid.String(), bid.String())
	if err != nil {
		return 0, err
	}
	if exists {
		count, _ := s.LikeRepo.CountByBlogID(bid.String())
		return int(count), nil
	}

	like := &model.Like{
		UserID: uid,
		BlogID: bid,
	}
	if err := s.LikeRepo.Create(like); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(bid.String())
	return int(count), nil
}

func (s *BlogService) Unlike(blogID, userID string) (int, error) {
	bid, err := uuid.Parse(blogID)
	if err != nil {
		return 0, err
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	if err := s.LikeRepo.Delete(uid.String(), bid.String()); err != nil {
		return 0, err
	}

	count, _ := s.LikeRepo.CountByBlogID(bid.String())
	return int(count), nil
}
