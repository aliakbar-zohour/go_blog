// middleware/compress: Gzip compression for JSON responses to reduce bandwidth and improve speed.
package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	w *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.w.Write(b)
}

// Gzip compresses response with gzip when client accepts it.
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz := gzip.NewWriter(w)
		defer gz.Close()
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, w: gz}, r)
	})
}
