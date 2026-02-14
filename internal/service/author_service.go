// service/author_service: Business logic for authors.
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

type AuthorService struct {
	repo *repository.AuthorRepository
	cfg  *config.Config
}

func NewAuthorService(repo *repository.AuthorRepository, cfg *config.Config) *AuthorService {
	return &AuthorService{repo: repo, cfg: cfg}
}

func (s *AuthorService) Create(ctx context.Context, name string, avatar *multipart.FileHeader) (*model.Author, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	a := &model.Author{Name: name}
	if avatar != nil {
		maxBytes := int64(s.cfg.MaxFileMB * 1024 * 1024)
		path, err := upload.SaveSingleImage(avatar, s.cfg.UploadDir, "avatars", maxBytes)
		if err == nil {
			a.AvatarPath = path
		}
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, a.ID)
}

func (s *AuthorService) GetByID(ctx context.Context, id uint) (*model.Author, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *AuthorService) List(ctx context.Context) ([]model.Author, error) {
	return s.repo.List(ctx)
}

func (s *AuthorService) Update(ctx context.Context, id uint, name string, avatar *multipart.FileHeader) (*model.Author, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if name != "" {
		a.Name = strings.TrimSpace(name)
	}
	if avatar != nil {
		maxBytes := int64(s.cfg.MaxFileMB * 1024 * 1024)
		path, err := upload.SaveSingleImage(avatar, s.cfg.UploadDir, "avatars", maxBytes)
		if err == nil {
			a.AvatarPath = path
		}
	}
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *AuthorService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
