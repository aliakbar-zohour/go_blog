// repository/email_verification_repository: Store and lookup email verification codes.
package repository

import (
	"context"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/gorm"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(ctx context.Context, ev *model.EmailVerification) error {
	return r.db.WithContext(ctx).Create(ev).Error
}

func (r *EmailVerificationRepository) FindValid(ctx context.Context, email, code string) (*model.EmailVerification, error) {
	var ev model.EmailVerification
	err := r.db.WithContext(ctx).Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).First(&ev).Error
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *EmailVerificationRepository) DeleteByEmail(ctx context.Context, email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(&model.EmailVerification{}).Error
}
