// middleware/ratelimit: In-memory rate limiter by client IP to mitigate brute-force and abuse.
package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aliakbar-zohour/go_blog/pkg/response"
)

// RateLimit allows at most maxRequests requests per window per IP. Window is fixed (e.g. 1 minute).
type RateLimit struct {
	mu       sync.Mutex
	m        map[string][]time.Time
	maxReq   int
	window   time.Duration
	cleanup  time.Time
}

// NewRateLimit returns a middleware that limits to maxRequests per window per IP.
func NewRateLimit(maxRequests int, window time.Duration) *RateLimit {
	return &RateLimit{
		m:      make(map[string][]time.Time),
		maxReq: maxRequests,
		window: window,
	}
}

func (rl *RateLimit) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		rl.mu.Lock()
		now := time.Now()
		if now.After(rl.cleanup) {
			rl.cleanup = now.Add(rl.window)
			rl.m = make(map[string][]time.Time)
		}
		cutoff := now.Add(-rl.window)
		times := rl.m[ip]
		n := 0
		for _, t := range times {
			if t.After(cutoff) {
				times[n] = t
				n++
			}
		}
		times = times[:n]
		if len(times) >= rl.maxReq {
			rl.mu.Unlock()
			response.Err(w, http.StatusTooManyRequests, "too many requests")
			return
		}
		rl.m[ip] = append(times, now)
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	if x := r.Header.Get("X-Forwarded-For"); x != "" {
		if first, _, ok := strings.Cut(x, ","); ok {
			return strings.TrimSpace(first)
		}
		return strings.TrimSpace(x)
	}
	if x := r.Header.Get("X-Real-IP"); x != "" {
		return strings.TrimSpace(x)
	}
	return r.RemoteAddr
}
