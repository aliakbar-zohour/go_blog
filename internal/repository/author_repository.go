// repository/author_repository: Data access for authors (CRUD).
package repository

import (
	"context"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Create(ctx context.Context, a *model.Author) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *AuthorRepository) GetByID(ctx context.Context, id uint) (*model.Author, error) {
	var a model.Author
	err := r.db.WithContext(ctx).First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AuthorRepository) GetByEmail(ctx context.Context, email string) (*model.Author, error) {
	var a model.Author
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AuthorRepository) List(ctx context.Context) ([]model.Author, error) {
	var list []model.Author
	err := r.db.WithContext(ctx).Order("name").Find(&list).Error
	return list, err
}

func (r *AuthorRepository) Update(ctx context.Context, a *model.Author) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *AuthorRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Author{}, id).Error
}
