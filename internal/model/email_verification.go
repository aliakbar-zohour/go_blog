// model/email_verification: Temporary code sent to email for sign-up verification.
package model

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerification struct {
	ID        uint           `gorm:"primaryKey" json:"-"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"-"`
	Code      string         `gorm:"size:10;not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
