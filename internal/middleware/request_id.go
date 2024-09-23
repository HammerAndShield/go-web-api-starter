package middleware

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
)

const (
	RequestIdLog = "requestId"
)

func GetRequestID(r *http.Request) (string, error) {
	id := r.Header.Get("X-Request-Id")
	if id == "" {
		return "", errors.New("no X-Request-Id header")
	}

	return id, nil
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		r.Header.Set("X-Request-ID", id.String())
		w.Header().Set("X-Request-ID", id.String())

		next.ServeHTTP(w, r)
	})
}
