// model/author: Author entity (name, avatar, email auth). Used as the writer account.
package model

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	AvatarPath      string         `gorm:"size:512" json:"avatar_path,omitempty"`
	Email           *string        `gorm:"size:255;uniqueIndex" json:"email,omitempty"`
	PasswordHash    string         `gorm:"size:255" json:"-"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
