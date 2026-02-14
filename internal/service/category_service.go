// service/category_service: Business logic for categories.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"gorm.io/gorm"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

const maxCategoryNameLen = 200

func (s *CategoryService) Create(ctx context.Context, name string) (*model.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if len(name) > maxCategoryNameLen {
		return nil, errors.New("name too long")
	}
	c := &model.Category{Name: name}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, c.ID)
}

func (s *CategoryService) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) List(ctx context.Context) ([]model.Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) Update(ctx context.Context, id uint, name string) (*model.Category, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if name != "" {
		n := strings.TrimSpace(name)
		if len(n) > maxCategoryNameLen {
			return nil, errors.New("name too long")
		}
		c.Name = n
	}
	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
