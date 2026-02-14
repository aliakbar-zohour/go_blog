// repository/media_repository: Create and delete media records in the database.
package repository

import (
	"context"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(ctx context.Context, m *model.Media) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MediaRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Media{}, id).Error
}

func (r *MediaRepository) GetByID(ctx context.Context, id uint) (*model.Media, error) {
	var m model.Media
	err := r.db.WithContext(ctx).First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}
