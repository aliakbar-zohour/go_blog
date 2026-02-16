// handler/health_handler: Health check for load balancers and orchestration.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/aliakbar-zohour/go_blog/pkg/response"
	"gorm.io/gorm"
)

// HealthHandler handles health and readiness checks.
type HealthHandler struct {
	db *gorm.DB
}

// NewHealthHandler returns a HealthHandler. db may be nil to skip DB check.
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health godoc
//
//	@Summary		Health check
//	@Description	Returns 200 if the API is up. Checks DB connectivity when database is configured.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	response.Body{data=object}
//	@Failure		503	{object}	response.Body
//	@Router			/health [get]
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			response.Err(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if err := sqlDB.PingContext(ctx); err != nil {
			response.Err(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
	}
	response.OK(w, map[string]string{"status": "ok"})
}
