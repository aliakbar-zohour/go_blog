// cmd/api: Application entry point; loads config, DB, services and starts the HTTP server.
//
//	@title			Go Blog API
//	@version		1.0
//	@description	REST API for a blog with auth (register with email code, login), CRUD posts (create/update/delete require JWT), authors, categories, comments. Errors are returned in the response body with an "error" field.
//	@host			localhost:8080
//	@BasePath		/api
//	@schemes		http
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
package main

import (
	"log"
	"net/http"

	_ "github.com/aliakbar-zohour/go_blog/docs"
	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/database"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"github.com/aliakbar-zohour/go_blog/internal/router"
	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
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
	r := router.New(postSvc, authorSvc, categorySvc, commentSvc, authSvc, cfg)
	addr := ":" + cfg.ServerPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
