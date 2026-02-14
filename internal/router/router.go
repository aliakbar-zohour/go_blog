// router: تعریف مسیرهای API و اتصال هندلرها و میدلورها.
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

func New(svc *service.PostService, cfg *config.Config) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recover, middleware.SecureHeaders, middleware.Log)
	r.Get("/docs/*", httpSwagger.WrapHandler)
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))
	r.Route("/api", func(r chi.Router) {
		r.Route("/posts", func(r chi.Router) {
			h := handler.NewPostHandler(svc)
			r.Get("/", h.List)
			r.Post("/", h.Create)
			r.Get("/{id}", h.GetByID)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
		})
	})
	return r
}
