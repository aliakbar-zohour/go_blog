// service/post_service: Business logic for posts and media (validation, orchestration).
package service

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/model"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"github.com/aliakbar-zohour/go_blog/internal/upload"
	"gorm.io/gorm"
)

type PostService struct {
	postRepo  *repository.PostRepository
	mediaRepo *repository.MediaRepository
	cfg       *config.Config
}

func NewPostService(postRepo *repository.PostRepository, mediaRepo *repository.MediaRepository, cfg *config.Config) *PostService {
	return &PostService{postRepo: postRepo, mediaRepo: mediaRepo, cfg: cfg}
}

func (s *PostService) Create(ctx context.Context, title, body string, files []*multipart.FileHeader) (*model.Post, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title is required")
	}
	post := &model.Post{Title: title, Body: trim(body)}
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}
	maxBytes := int64(s.cfg.MaxFileMB * 1024 * 1024)
	for _, f := range files {
		m, _, err := upload.SaveFile(f, s.cfg.UploadDir, post.ID, maxBytes)
		if err != nil {
			continue
		}
		_ = s.mediaRepo.Create(ctx, m)
	}
	return s.postRepo.GetByID(ctx, post.ID)
}

func (s *PostService) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	return s.postRepo.GetByID(ctx, id)
}

func (s *PostService) List(ctx context.Context, limit, offset int) ([]model.Post, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.postRepo.List(ctx, limit, offset)
}

func (s *PostService) Update(ctx context.Context, id uint, title, body string, files []*multipart.FileHeader) (*model.Post, error) {
	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if title != "" {
		post.Title = strings.TrimSpace(title)
	}
	if body != "" {
		post.Body = strings.TrimSpace(body)
	}
	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}
	maxBytes := int64(s.cfg.MaxFileMB * 1024 * 1024)
	for _, f := range files {
		m, _, err := upload.SaveFile(f, s.cfg.UploadDir, post.ID, maxBytes)
		if err != nil {
			continue
		}
		_ = s.mediaRepo.Create(ctx, m)
	}
	return s.postRepo.GetByID(ctx, id)
}

func (s *PostService) Delete(ctx context.Context, id uint) error {
	return s.postRepo.Delete(ctx, id)
}

func trim(s string) string {
	const max = 10000
	if len(s) > max {
		return s[:max]
	}
	return strings.TrimSpace(s)
}
