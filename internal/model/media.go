// model/media: مدل فایل‌های آپلودشده (عکس/ویدیو) برای هر پست.
package model

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type Media struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PostID    uint           `gorm:"not null;index" json:"post_id"`
	Type      MediaType      `gorm:"size:20;not null" json:"type"`
	Path      string         `gorm:"size:512;not null" json:"path"`
	Filename  string         `gorm:"size:255" json:"filename"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
