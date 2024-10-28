package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go-web-api-starter/internal/apiutils"
	"go-web-api-starter/internal/database"
	"go-web-api-starter/internal/jwt"
	"log/slog"
	"net/http"
	"strings"
)

//region user auth and context setting middleware

type userGetter interface {
	GetById(id uuid.UUID) (*User, error)
}

type defaultUserInserter interface {
	InsertDefaultUser(email string, id uuid.UUID) error
}

type userGetterInserter interface {
	userGetter
	defaultUserInserter
}

type JWTReader interface {
	Read(tokenString string) (jwt.MapClaims, error)
	ValidateClaims(claims jwt.MapClaims) error
}

func Authenticate(
	logger *slog.Logger,
	reader JWTReader,
	userGetterInserter userGetterInserter,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the auth header from the request and extract/verify it
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apiutils.ErrorResponse(w, r, logger, http.StatusBadRequest, "missing authorization header")
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				apiutils.ErrorResponse(w, r, logger, http.StatusBadRequest, "invalid authorization header")
				return
			}

			tokenString := headerParts[1]

			// Read the claims and validate them
			claims, err := reader.Read(tokenString)
			if err != nil {
				logger.Error("error reading jwt", "error", err, "token", tokenString)
				apiutils.ErrorResponse(w, r, logger, http.StatusUnauthorized, "invalid token")
				return
			}

			err = reader.ValidateClaims(claims)
			if err != nil {
				message := ""
				switch {
				case errors.Is(err, jwt.ErrExpiredToken):
					message = "token is expired"
				case errors.Is(err, jwt.ErrInvalidIssuer):
					message = "invalid issuer on jwt"
				case errors.Is(err, jwt.ErrEmptySubject):
					message = "the jwt has no subject"
				}

				apiutils.ErrorResponse(w, r, logger, http.StatusUnauthorized, message)
				return
			}

			// Retrieve the user from the database and add it to the context
			userId, err := claims.GetSubject()
			if err != nil {
				apiutils.ErrorResponse(w, r, logger, http.StatusBadRequest, "userId not in JWT subject")
				return
			}

			userUuid, err := uuid.Parse(userId)
			if err != nil {
				apiutils.ErrorResponse(w, r, logger, http.StatusBadRequest, "userId not a valid UUID")
				return
			}

			user, err := userGetterInserter.GetById(userUuid)
			if err != nil {
				switch {
				case errors.Is(err, database.ErrRecordNotFound):
					// If the user does not exist, but the JWT is valid, we can insert them into the database
					email := claims["email"].(string)
					user, err = insertAndRetrieveUnknownUser(userGetterInserter, email, userUuid)
					if err != nil {
						apiutils.ServerErrorResponse(w, r, logger, err)
						return
					}
				default:
					apiutils.ServerErrorResponse(w, r, logger, err)
					return
				}
			}

			ur := contextSetUser(r, user)

			// Log the auth so we can associate with a request_id
			requestId := r.Header.Get("X-Request-ID")
			if requestId == "" {
				apiutils.ServerErrorResponse(w, r, logger, fmt.Errorf("missing X-Request-Id header"))
				return
			}
			logger.Info("user authenticated", "request id", requestId, "user id", userId)

			next.ServeHTTP(w, ur)
		})
	}
}

func insertAndRetrieveUnknownUser(uGetterInserter userGetterInserter, email string, id uuid.UUID) (*User, error) {
	// Insert the new user
	err := uGetterInserter.InsertDefaultUser(email, id)
	if err != nil {
		return nil, err
	}

	// Retrieve the user and assign them to the empty user
	user, err := uGetterInserter.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

//endregion

//region require permissions middleware

// RequirePermissions creates a middleware that checks if the authenticated user's role
// has the required permissions available in the ...string argument. These permissions are
// checked against a set of permissions associated with this user's role. The permissions set
// is an implementation-dependent feature of the role.
//
// The middleware uses the contextGetUser() method to retrieve the user object from the request's context.
// If the user is authenticated as the AnonymousUser or if the user's role does not include
// one or more of the required permissions, it responds with an "Unauthorized" status.
// If the user's role includes all the required permissions, the request handling continues
// with the next middleware in the chain.
//
// This middleware must be called after getting the User, or it will panic.
func RequirePermissions(logger *slog.Logger, permissions ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := ContextGetUser(r)

			if !user.Role.Permissions.Includes(permissions...) {
				apiutils.ForbiddenResponse(w, r, logger)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

//endregion
