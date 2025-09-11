package service

import (
	"blog/model"
	"blog/repo"
)

type CommentService struct {
	CommentRepo *repo.CommentRepository
}

// Kreiranje komentara
func (s *CommentService) CreateComment(comment *model.Comment) error {
	return s.CommentRepo.Create(comment)
}

// Dohvatanje komentara za blog
func (s *CommentService) GetComments(blogID string) ([]model.Comment, error) {
	return s.CommentRepo.GetByBlogID(blogID)
}
