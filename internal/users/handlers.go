package users

import (
	"go-web-api-starter/internal/apiutils"
	"go-web-api-starter/internal/middleware"
	"log/slog"
	"net/http"
)

func GetCurrentUserHandler(
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ContextGetUser(r)
		responseData := apiutils.Envelope{"user": user}

		err := apiutils.WriteJson(w, http.StatusOK, responseData, http.Header{})
		if err != nil {
			requestId, _ := middleware.GetRequestID(r)
			logger.Error("GetCurrentUserHandler write failed", middleware.RequestIdLog, requestId, "error", err)
		}
	})
}
