// service/comment_service: Business logic for comments.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"gorm.io/gorm"
)

type CommentService struct {
	repo     *repository.CommentRepository
	postRepo *repository.PostRepository
}

func NewCommentService(repo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{repo: repo, postRepo: postRepo}
}

func (s *CommentService) Create(ctx context.Context, postID uint, body, authorName string) (*model.Comment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.New("body is required")
	}
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	c := &model.Comment{PostID: postID, Body: body, AuthorName: strings.TrimSpace(authorName)}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, c.ID)
}

func (s *CommentService) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	return s.repo.ListByPostID(ctx, postID)
}

func (s *CommentService) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) Update(ctx context.Context, id uint, body string) (*model.Comment, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if body != "" {
		c.Body = strings.TrimSpace(body)
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
