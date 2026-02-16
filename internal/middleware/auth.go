// middleware/auth: Extracts JWT and sets author ID in request context.
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aliakbar-zohour/go_blog/pkg/auth"
	"github.com/aliakbar-zohour/go_blog/pkg/response"
)

type contextKey string

const AuthorIDKey contextKey = "author_id"

// RequireAuth validates the Bearer token and sets author_id in context. Returns 401 if missing or invalid.
func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.UnauthorizedWithCode(w, "auth_required", "authorization required")
				return
			}
			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.UnauthorizedWithCode(w, "invalid_header", "invalid authorization header")
				return
			}
			claims, err := auth.ParseToken(parts[1], secret)
			if err != nil {
				response.UnauthorizedWithCode(w, "invalid_token", "invalid or expired token")
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
