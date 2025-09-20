package httpx

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

// RequestID injects a simple request ID into the context and header.
func RequestID(next http.Handler) http.Handler {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := fmt.Sprintf("%08x", rnd.Uint32())
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger logs basic request/response details.
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, _ := r.Context().Value(requestIDKey).(string)
			start := time.Now()

			ww := &wrapWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(ww, r)

			logger.Info("http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"remote", remoteIP(r.RemoteAddr),
				"status", ww.status,
				"duration_ms", time.Since(start).Milliseconds(),
				"request_id", id,
			)
		})
	}
}

type wrapWriter struct {
	http.ResponseWriter
	status int
}

func (w *wrapWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func remoteIP(addr string) string {
	if i := strings.LastIndex(addr, ":"); i != -1 {
		return addr[:i]
	}
	return addr
}
