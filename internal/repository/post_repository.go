// repository/post_repository: Data access layer for posts (CRUD).
package repository

import (
	"context"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) GetByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.WithContext(ctx).Preload("Media").Preload("Author").Preload("Category").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) List(ctx context.Context, limit, offset int, categoryID *uint) ([]model.Post, error) {
	var posts []model.Post
	q := r.db.WithContext(ctx).Preload("Media").Preload("Author").Preload("Category").Limit(limit).Offset(offset).Order("created_at DESC")
	if categoryID != nil && *categoryID > 0 {
		q = q.Where("category_id = ?", *categoryID)
	}
	err := q.Find(&posts).Error
	return posts, err
}

func (r *PostRepository) Update(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Post{}, id).Error
}
