package logs

import (
	"net/http"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
)

type statsWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statsWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statsWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sw := &statsWriter{
			ResponseWriter: w,
		}

		reqID := middleware.GetReqID(r.Context())

		l := slog.Default().With(
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
		)

		ctx := into(r.Context(), l)
		r = r.WithContext(ctx)

		start := time.Now()

		next.ServeHTTP(sw, r)

		l.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"bytes", sw.bytes,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}
