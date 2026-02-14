// middleware/limit: Request body size limit to prevent large payload attacks.
package middleware

import (
	"net/http"

	"github.com/aliakbar-zohour/go_blog/pkg/response"
)

// MaxBytes limits the request body size. Returns 413 with JSON if Content-Length exceeds max or read exceeds max.
func MaxBytes(max int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > max {
				response.Err(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, max)
			}
			next.ServeHTTP(w, r)
		})
	}
}
