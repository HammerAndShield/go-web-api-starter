package middleware

import (
	"log/slog"
	"net/http"
)

func RecoverPanic(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestId, _ := GetRequestID(r)
					logger.Error("panic encountered", RequestIdLog, requestId, "error", err)
					w.Header().Set("Connection", "close")
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
