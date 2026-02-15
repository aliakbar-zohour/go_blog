// model/author: Author model (name and avatar image).
package model

import "time"

type Author struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	AvatarPath string    `gorm:"size:512" json:"avatar_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
