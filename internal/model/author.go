// model/author: Author model (name and avatar image).
package model

import "time"

type Author struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Name            string     `gorm:"size:255;not null" json:"name"`
	AvatarPath      string     `gorm:"size:512" json:"avatar_path"`
	Email           *string    `gorm:"size:255;uniqueIndex" json:"email,omitempty"`
	PasswordHash    string     `gorm:"size:255" json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
