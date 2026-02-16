package service

import (
	"context"
	"testing"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/model"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	// SQLite driver requires CGO on Windows; skip if unavailable.
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Skipf("sqlite (CGO) not available: %v", err)
	}
	if err := db.AutoMigrate(&model.Post{}, &model.Media{}, &model.Author{}, &model.Category{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func TestPostService_List_ReturnsTotalAndItems(t *testing.T) {
	db := setupTestDB(t)
	postRepo := repository.NewPostRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	cfg := &config.Config{UploadDir: "uploads", MaxFileMB: 50}
	svc := NewPostService(postRepo, mediaRepo, cfg)
	ctx := context.Background()

	// Empty list
	result, err := svc.List(ctx, 10, 0, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 0 || len(result.Items) != 0 {
		t.Errorf("empty DB: want total=0, items=0; got total=%d, items=%d", result.Total, len(result.Items))
	}

	// Create a post via repo (no handler)
	post := &model.Post{Title: "Test", Body: "Body", AuthorID: 1, CategoryID: 1}
	if err := postRepo.Create(ctx, post); err != nil {
		t.Fatalf("Create post: %v", err)
	}

	result, err = svc.List(ctx, 10, 0, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 {
		t.Errorf("one post: want total=1, items=1; got total=%d, items=%d", result.Total, len(result.Items))
	}
	if result.Items[0].Title != "Test" {
		t.Errorf("item title want Test, got %s", result.Items[0].Title)
	}
}

func TestPostService_List_RespectsLimitAndOffset(t *testing.T) {
	db := setupTestDB(t)
	postRepo := repository.NewPostRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	cfg := &config.Config{UploadDir: "uploads", MaxFileMB: 50}
	svc := NewPostService(postRepo, mediaRepo, cfg)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		post := &model.Post{Title: "Post", Body: "", AuthorID: 1, CategoryID: 1}
		_ = postRepo.Create(ctx, post)
	}

	result, err := svc.List(ctx, 2, 1, nil)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if result.Total != 5 || len(result.Items) != 2 {
		t.Errorf("want total=5, items=2; got total=%d, items=%d", result.Total, len(result.Items))
	}
}
