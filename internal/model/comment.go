// model/comment: Comment on a post (body + optional author name).
package model

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PostID     uint           `gorm:"not null;index" json:"post_id"`
	Body       string         `gorm:"type:text;not null" json:"body"`
	AuthorID   *uint          `gorm:"index" json:"author_id,omitempty"`   // set when user is logged in
	AuthorName string         `gorm:"size:255" json:"author_name"`        // optional display name (guest or override)
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
