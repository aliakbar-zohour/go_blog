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

var (
	ErrCommentForbidden = errors.New("you can only edit your own comment")
	ErrDeleteCommentForbidden = errors.New("you can only delete your own comment")
)

type CommentService struct {
	repo     *repository.CommentRepository
	postRepo *repository.PostRepository
}

func NewCommentService(repo *repository.CommentRepository, postRepo *repository.PostRepository) *CommentService {
	return &CommentService{repo: repo, postRepo: postRepo}
}

const maxCommentBodyLen = 2000

func (s *CommentService) Create(ctx context.Context, postID uint, body, authorName string, authorID *uint) (*model.Comment, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.New("body is required")
	}
	if len(body) > maxCommentBodyLen {
		return nil, errors.New("comment body too long")
	}
	if _, err := s.postRepo.GetByID(ctx, postID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	c := &model.Comment{PostID: postID, Body: body, AuthorID: authorID, AuthorName: strings.TrimSpace(authorName)}
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

func (s *CommentService) Update(ctx context.Context, id uint, body string, authorID uint) (*model.Comment, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if c.AuthorID == nil || *c.AuthorID != authorID {
		return nil, ErrCommentForbidden
	}
	if body != "" {
		b := strings.TrimSpace(body)
		if len(b) > maxCommentBodyLen {
			return nil, errors.New("comment body too long")
		}
		c.Body = b
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) Delete(ctx context.Context, id uint, authorID uint) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}
	if c.AuthorID == nil || *c.AuthorID != authorID {
		return ErrDeleteCommentForbidden
	}
	return s.repo.Delete(ctx, id)
}
