// middleware/requestid: Injects a unique request ID into context and response header.
package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type requestIDCtxKey string

const RequestIDKey requestIDCtxKey = "request_id"

// RequestID generates a 16-byte hex ID, sets X-Request-ID header and context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), RequestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID returns the request ID from context, or empty string.
func GetRequestID(ctx context.Context) string {
	v := ctx.Value(RequestIDKey)
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback"
	}
	return hex.EncodeToString(b)
}
