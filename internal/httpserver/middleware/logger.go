package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/stratflow-labs/stratflow/internal/authkit"
	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(p []byte) (int, error) {
	if lrw.status == 0 {
		lrw.status = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(p)
	lrw.size += int64(n)
	return n, err
}

// Logger writes structured request logs via foundation slog logger.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		status := lrw.status
		if status == 0 {
			status = http.StatusOK
		}

		reqID, _ := RequestIDFromContext(r.Context())
		userID, _ := authkit.UserIDFromContext(r.Context())

		logger.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration", time.Since(start).Seconds(),
			"size", strconv.FormatInt(lrw.size, 10),
			"request_id", reqID,
			"user_id", userID,
		)
	})
}
