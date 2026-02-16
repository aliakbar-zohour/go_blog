// Integration test: requires PostgreSQL (env or .env). Skips if DB unavailable.
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/database"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"github.com/aliakbar-zohour/go_blog/internal/router"
	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/joho/godotenv"
)

func TestIntegration_HealthAndPosts(t *testing.T) {
	_ = godotenv.Load()
	if os.Getenv("DB_HOST") == "" && os.Getenv("DB_NAME") == "" {
		t.Skip("skip integration test when DB env not set")
	}
	cfg := config.Load()
	db, err := database.New(cfg)
	if err != nil {
		t.Skipf("database not available (set env for integration): %v", err)
	}
	sqlDB, _ := db.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("database ping failed: %v", err)
	}

	postRepo := repository.NewPostRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	authorRepo := repository.NewAuthorRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	evRepo := repository.NewEmailVerificationRepository(db)
	postSvc := service.NewPostService(postRepo, mediaRepo, cfg)
	authorSvc := service.NewAuthorService(authorRepo, cfg)
	categorySvc := service.NewCategoryService(categoryRepo)
	commentSvc := service.NewCommentService(commentRepo, postRepo)
	authSvc := service.NewAuthService(authorRepo, evRepo, cfg)
	r := router.New(db, postSvc, authorSvc, categorySvc, commentSvc, authSvc, cfg)

	// GET /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("GET /health: status want 200, got %d", rr.Code)
	}

	// GET /api/posts
	req2 := httptest.NewRequest(http.MethodGet, "/api/posts?limit=5", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusOK {
		t.Errorf("GET /api/posts: status want 200, got %d", rr2.Code)
	}
}
