// router: API route definitions and handler/middleware wiring.
package router

import (
	"net/http"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/handler"
	"github.com/aliakbar-zohour/go_blog/internal/middleware"
	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(postSvc *service.PostService, authorSvc *service.AuthorService, categorySvc *service.CategoryService, commentSvc *service.CommentService, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recover, middleware.SecureHeaders, middleware.Log)
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))
	r.Route("/api", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			ph := handler.NewPostHandler(postSvc)
			r.Get("/", ph.List)
			r.Post("/", ph.Create)
			r.Route("/{postId}/comments", func(r chi.Router) {
				ch := handler.NewCommentHandler(commentSvc)
				r.Get("/", ch.ListByPostID)
				r.Post("/", ch.Create)
			})
			r.Get("/{id}", ph.GetByID)
			r.Put("/{id}", ph.Update)
			r.Delete("/{id}", ph.Delete)
		})
		r.Route("/authors", func(r chi.Router) {
			ah := handler.NewAuthorHandler(authorSvc)
			r.Get("/", ah.List)
			r.Post("/", ah.Create)
			r.Get("/{id}", ah.GetByID)
			r.Put("/{id}", ah.Update)
			r.Delete("/{id}", ah.Delete)
		})
		r.Route("/categories", func(r chi.Router) {
			ch := handler.NewCategoryHandler(categorySvc)
			r.Get("/", ch.List)
			r.Post("/", ch.Create)
			r.Get("/{id}", ch.GetByID)
			r.Put("/{id}", ch.Update)
			r.Delete("/{id}", ch.Delete)
		})
		r.Route("/comments", func(r chi.Router) {
			ch := handler.NewCommentHandler(commentSvc)
			r.Put("/{id}", ch.Update)
			r.Delete("/{id}", ch.Delete)
		})
	})
	return r
}
