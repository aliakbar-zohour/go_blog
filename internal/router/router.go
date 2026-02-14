// router: API route definitions and handler/middleware wiring.
package router

import (
	"net/http"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/handler"
	"github.com/aliakbar-zohour/go_blog/internal/middleware"
	"github.com/aliakbar-zohour/go_blog/internal/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(postSvc *service.PostService, authorSvc *service.AuthorService, categorySvc *service.CategoryService, commentSvc *service.CommentService, authSvc *service.AuthService, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recover, middleware.SecureHeaders, middleware.CORS(cfg.CORSOrigins), middleware.Gzip, middleware.Log)
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.MaxBytes(cfg.BodyLimitBytes))
		authRateLimit := middleware.NewRateLimit(cfg.AuthRatePerMin, time.Minute)
		r.Route("/auth", func(r chi.Router) {
			r.Use(authRateLimit.Middleware)
			authH := handler.NewAuthHandler(authSvc)
			r.Post("/register/request", authH.RequestVerification)
			r.Post("/register/verify", authH.VerifyAndRegister)
			r.Post("/login", authH.Login)
		})
		authMW := middleware.RequireAuth(cfg.JWTSecret)
		r.Route("/posts", func(r chi.Router) {
			ph := handler.NewPostHandler(postSvc)
			r.Get("/", ph.List)
			r.Route("/{postId}/comments", func(r chi.Router) {
				ch := handler.NewCommentHandler(commentSvc)
				r.Get("/", ch.ListByPostID)
				r.With(authMW).Post("/", ch.Create)
			})
			r.Get("/{id}", ph.GetByID)
			r.With(authMW).Post("/", ph.Create)
			r.With(authMW).Put("/{id}", ph.Update)
			r.With(authMW).Delete("/{id}", ph.Delete)
		})
		r.Route("/authors", func(r chi.Router) {
			ah := handler.NewAuthorHandler(authorSvc)
			r.Get("/", ah.List)
			r.Post("/", ah.Create)
			r.Get("/{id}", ah.GetByID)
			r.With(authMW).Put("/{id}", ah.Update)
			r.With(authMW).Delete("/{id}", ah.Delete)
		})
		r.Route("/categories", func(r chi.Router) {
			ch := handler.NewCategoryHandler(categorySvc)
			r.Get("/", ch.List)
			r.With(authMW).Post("/", ch.Create)
			r.Get("/{id}", ch.GetByID)
			r.With(authMW).Put("/{id}", ch.Update)
			r.With(authMW).Delete("/{id}", ch.Delete)
		})
		r.Route("/comments", func(r chi.Router) {
			ch := handler.NewCommentHandler(commentSvc)
			r.With(authMW).Put("/{id}", ch.Update)
			r.With(authMW).Delete("/{id}", ch.Delete)
		})
	})
	return r
}
