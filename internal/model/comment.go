// model/comment: Comment model for post comments.
package model

import "time"

type Comment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PostID     uint      `gorm:"not null;index" json:"post_id"`
	Body       string    `gorm:"type:text;not null" json:"body"`
	AuthorID   *uint     `gorm:"index" json:"author_id,omitempty"`
	AuthorName string    `gorm:"size:255;not null" json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
