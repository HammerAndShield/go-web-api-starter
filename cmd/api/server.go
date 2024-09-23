package main

import (
	"encoding/json"
	"go-web-api-starter/internal/middleware"
	"log/slog"
	"net/http"
)

func newServer(
	logger *slog.Logger,
) http.Handler {
	v1Mux := http.NewServeMux()

	addRoutesV1(v1Mux)

	mux := http.NewServeMux()
	mux.Handle("/v1/", v1Mux)
	mux.Handle("/ping", ping(logger))

	recoverM := middleware.RecoverPanic(logger)
	loggerM := middleware.Logger(logger, []string{"/ping"})

	var server http.Handler = mux
	server = loggerM(server)
	server = middleware.RealIP(server)
	server = recoverM(server)

	return server
}

func ping(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status": "okay",
		}

		body, err := json.Marshal(response)
		if err != nil {
			logger.Error("failed to marshal ping response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if n, err := w.Write(body); err != nil {
			logger.Error("failed to write ping response", "error", err, "bytes_written", n)
		}
	})
}
