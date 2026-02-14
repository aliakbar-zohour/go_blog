// repository/category_repository: Data access for categories (CRUD).
package repository

import (
	"context"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, c *model.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*model.Category, error) {
	var c model.Category
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) List(ctx context.Context) ([]model.Category, error) {
	var list []model.Category
	err := r.db.WithContext(ctx).Order("name").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) Update(ctx context.Context, c *model.Category) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}
