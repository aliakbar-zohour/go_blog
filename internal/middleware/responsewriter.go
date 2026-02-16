// middleware/responsewriter: Wraps ResponseWriter to capture status code and size.
package middleware

import (
	"net/http"
)

// responseWriter wraps http.ResponseWriter to record status and bytes written.
type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func (w *responseWriter) unwrap() http.ResponseWriter { return w.ResponseWriter }
