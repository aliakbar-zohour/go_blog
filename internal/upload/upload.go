// upload: اعتبارسنجی و ذخیرهٔ امن فایل‌های تصویر و ویدئو.
package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/model"
)

var (
	allowedImages = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	allowedVideos = map[string]bool{".mp4": true, ".webm": true, ".mov": true}
)

func SaveFile(file *multipart.FileHeader, uploadDir string, postID uint, maxBytes int64) (*model.Media, string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	mediaType := model.MediaTypeImage
	if allowedVideos[ext] {
		mediaType = model.MediaTypeVideo
	} else if !allowedImages[ext] {
		return nil, "", fmt.Errorf("نوع فایل مجاز نیست")
	}
	if file.Size > maxBytes {
		return nil, "", fmt.Errorf("حجم فایل بیش از حد مجاز است")
	}
	dir := filepath.Join(uploadDir, "posts", fmt.Sprintf("%d", postID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, "", err
	}
	newName := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), postID, ext)
	dstPath := filepath.Join(dir, newName)
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return nil, "", err
	}
	relPath := filepath.Join("posts", fmt.Sprintf("%d", postID), newName)
	return &model.Media{PostID: postID, Type: mediaType, Path: relPath, Filename: file.Filename}, relPath, nil
}
