// model/post: Post domain model and relation to media.
package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Body        string         `gorm:"type:text" json:"body"`
	BannerPath  string         `gorm:"size:512" json:"banner_path"`
	AuthorID    uint           `gorm:"index" json:"author_id"`
	Author      *Author        `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CategoryID  uint           `gorm:"index" json:"category_id"`
	Category    *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Media       []Media        `gorm:"foreignKey:PostID" json:"media,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
