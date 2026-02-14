// middleware/auth: Extracts JWT and sets author ID in request context.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aliakbar-zohour/go_blog/pkg/auth"
)

type contextKey string

const AuthorIDKey contextKey = "author_id"

// RequireAuth validates the Bearer token and sets author_id in context. Returns 401 if missing or invalid.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":"authorization required"}`))
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":"invalid authorization header"}`))
				return
			}
			claims, err := auth.ParseToken(parts[1], secret)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"success":false,"error":"invalid or expired token"}`))
				return
			}
			ctx := context.WithValue(r.Context(), AuthorIDKey, claims.AuthorID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetAuthorID returns the author ID from context, or 0 if not set.
func GetAuthorID(ctx context.Context) uint {
	v := ctx.Value(AuthorIDKey)
	if v == nil {
		return 0
	}
	id, _ := v.(uint)
	return id
}
