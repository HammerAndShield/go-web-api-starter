package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.responseSize += len(b)
	return n, err
}

// Logger returns a middleware that logs incoming HTTP requests.
// The middleware logs details such as duration, request ID, user information, URI, method,
// status code, IP address, user agent, request size, and response size.
func Logger(logger *slog.Logger, quietRoutes []string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(quietRoutes, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			currTime := time.Now()
			ww := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusTeapot,
				responseSize:   0,
			}

			defer func(start time.Time) {
				dur := time.Since(start)
				requestId := r.Header.Get("X-Request-Id")
				if requestId == "" {
					requestId = "unknown"
				}

				logger.Info(
					fmt.Sprintf("%s: %s", RequestIdLog, requestId),
					"duration_ms", dur.Milliseconds(),
					"uri", r.RequestURI,
					"method", r.Method,
					"status code", ww.statusCode,
					"requester ip", strings.Split(r.RemoteAddr, ":")[0],
					"user agent", r.UserAgent(),
					"request size", r.ContentLength,
					"response size", ww.responseSize,
				)
			}(currTime)

			next.ServeHTTP(ww, r)
		})
	}
}
