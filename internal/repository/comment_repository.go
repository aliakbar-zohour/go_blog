// repository/comment_repository: Data access for comments.
package repository

import (
	"context"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CommentRepository) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var c model.Comment
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CommentRepository) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	var list []model.Comment
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).Order("created_at ASC").Find(&list).Error
	return list, err
}

func (r *CommentRepository) Update(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}
