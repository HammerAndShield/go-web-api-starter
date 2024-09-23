package users

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-web-api-starter/internal/database"
	"go-web-api-starter/internal/jwtauth"
	"go-web-api-starter/internal/testutils"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

//region auth tests

// Mock implementations
type MockJWTReader struct {
	ReadFunc           func(tokenString string) (jwt.MapClaims, error)
	ValidateClaimsFunc func(claims jwt.MapClaims) error
}

func (m *MockJWTReader) Read(tokenString string) (jwt.MapClaims, error) {
	return m.ReadFunc(tokenString)
}

func (m *MockJWTReader) ValidateClaims(claims jwt.MapClaims) error {
	return m.ValidateClaimsFunc(claims)
}

type MockUserGetterInserter struct {
	GetByIdFunc           func(id uuid.UUID) (*User, error)
	InsertDefaultUserFunc func(email string, id uuid.UUID) error
}

func (m *MockUserGetterInserter) GetById(id uuid.UUID) (*User, error) {
	return m.GetByIdFunc(id)
}

func (m *MockUserGetterInserter) InsertDefaultUser(email string, id uuid.UUID) error {
	return m.InsertDefaultUserFunc(email, id)
}

// Helper functions
func createAuthTestRequest(method, target string, header string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	if header != "" {
		req.Header.Set("Authorization", header)
	}
	req.Header.Set("X-Request-ID", "test-request-id")
	return req
}

func createAuthTestHandler(jr JWTReader, ugi userGetterInserter) http.Handler {
	return Authenticate(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jr,
		ugi,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func createAuthTestHandlerWithMethod(
	jr JWTReader,
	ugi userGetterInserter,
	handle http.Handler,
) http.Handler {
	return Authenticate(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		jr,
		ugi,
	)(handle)
}

// Test cases
func TestAuthenticateMissingHeader(t *testing.T) {
	req := createAuthTestRequest("GET", "/", "")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(&MockJWTReader{}, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusBadRequest, "missing authorization header")
}

func TestAuthenticateInvalidHeaderFormat(t *testing.T) {
	req := createAuthTestRequest("GET", "/", "InvalidFormat")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(&MockJWTReader{}, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusBadRequest, "invalid authorization header")
}

func TestAuthenticateInvalidToken(t *testing.T) {
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return nil, errors.New("invalid token")
		},
	}

	req := createAuthTestRequest("GET", "/", "Bearer invalid_token")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusUnauthorized, "invalid token")
}

func TestAuthenticateExpiredToken(t *testing.T) {
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return jwtauth.ErrExpiredToken
		},
	}

	req := createAuthTestRequest("GET", "/", "Bearer expired_token")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusUnauthorized, "token is expired")
}

func TestAuthenticateInvalidIssuer(t *testing.T) {
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return jwtauth.ErrInvalidIssuer
		},
	}

	req := createAuthTestRequest("GET", "/", "Bearer invalid_issuer_token")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusUnauthorized, "invalid issuer on jwt")
}

func TestAuthenticateEmptySubject(t *testing.T) {
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return jwtauth.ErrEmptySubject
		},
	}

	req := createAuthTestRequest("GET", "/", "Bearer empty_subject_token")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusUnauthorized, "the jwt has no subject")
}

func TestAuthenticateValidTokenExistingUser(t *testing.T) {
	userID := uuid.New()
	testUser := &User{ID: userID, Email: "test@example.com"}
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": userID.String(), "email": "test@example.com"}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return nil
		},
	}
	mockUserGetterInserter := &MockUserGetterInserter{
		GetByIdFunc: func(id uuid.UUID) (*User, error) {
			return testUser, nil
		},
	}

	req := createAuthTestRequest("GET", "/", fmt.Sprintf("Bearer %s", userID))
	rec := httptest.NewRecorder()

	testUserHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contextUser := ContextGetUser(r)
			if contextUser == nil {
				t.Error("Expected user to be added to context, but it was not")
			} else if contextUser.ID != userID {
				t.Errorf("Expected user ID %s, got %s", userID, contextUser.ID)
			}
		})
	}

	handler := createAuthTestHandlerWithMethod(mockJWTReader, mockUserGetterInserter, testUserHandler())
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthenticateValidTokenNewUser(t *testing.T) {
	requestCount := 0
	userID := uuid.New()
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": userID.String(), "email": "newuser@example.com"}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return nil
		},
	}
	mockUserGetterInserter := &MockUserGetterInserter{
		GetByIdFunc: func(id uuid.UUID) (*User, error) {
			if requestCount == 0 {
				requestCount++
				return nil, database.ErrRecordNotFound
			} else {
				return &User{ID: userID, Email: "newuser@example.com"}, nil
			}
		},
		InsertDefaultUserFunc: func(email string, id uuid.UUID) error {
			return nil
		},
	}

	testUserHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contextUser := ContextGetUser(r)
			if contextUser == nil {
				t.Error("Expected user to be added to context, but it was not")
			} else if contextUser.ID != userID {
				t.Errorf("Expected user ID %s, got %s", userID, contextUser.ID)
			}
			w.WriteHeader(http.StatusOK)
		})
	}

	req := createAuthTestRequest("GET", "/", fmt.Sprintf("Bearer %s", userID))
	rec := httptest.NewRecorder()

	handler := createAuthTestHandlerWithMethod(mockJWTReader, mockUserGetterInserter, testUserHandler())
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestAuthenticateDatabaseError(t *testing.T) {
	userID := uuid.New()
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": userID.String(), "email": "test@example.com"}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return nil
		},
	}
	mockUserGetterInserter := &MockUserGetterInserter{
		GetByIdFunc: func(id uuid.UUID) (*User, error) {
			return nil, errors.New("database error")
		},
	}

	req := createAuthTestRequest("GET", "/", fmt.Sprintf("Bearer %s", userID))
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, mockUserGetterInserter)
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func TestAuthenticateInvalidUUID(t *testing.T) {
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": "not-a-valid-uuid", "email": "test@example.com"}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return nil
		},
	}

	req := createAuthTestRequest("GET", "/", "Bearer invalid-uuid-token")
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, &MockUserGetterInserter{})
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusBadRequest, "userId not a valid UUID")
}

func TestAuthenticateMissingRequestID(t *testing.T) {
	userID := uuid.New()
	mockJWTReader := &MockJWTReader{
		ReadFunc: func(tokenString string) (jwt.MapClaims, error) {
			return jwt.MapClaims{"sub": userID.String(), "email": "test@example.com"}, nil
		},
		ValidateClaimsFunc: func(claims jwt.MapClaims) error {
			return nil
		},
	}
	mockUserGetterInserter := &MockUserGetterInserter{
		GetByIdFunc: func(id uuid.UUID) (*User, error) {
			return &User{ID: id, Email: "test@example.com"}, nil
		},
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", userID))
	// Intentionally not setting X-Request-ID
	rec := httptest.NewRecorder()

	handler := createAuthTestHandler(mockJWTReader, mockUserGetterInserter)
	handler.ServeHTTP(rec, req)

	testutils.CheckJSONResponseError(t, rec, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

//endregion

//region permissions tests

// helpers
func addUserHandler(user *User, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nr := contextSetUser(r, user)
		next.ServeHTTP(w, nr)
	})
}

func createPermTestHandler(user *User, permissions ...string) http.Handler {
	return addUserHandler(user, RequirePermissions(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		permissions...,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
}

// tests

func TestPermissionsSuccess(t *testing.T) {
	user := &User{
		Role: Role{
			Name:        "admin",
			Permissions: []string{"read:sounds", "write:sounds"},
		},
	}

	testCases := [][]string{
		{"read:sounds"},
		{"read:sounds", "write:sounds"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("permissions %v", tc), func(t *testing.T) {
			rec := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/", nil)

			handler := createPermTestHandler(user, tc...)
			handler.ServeHTTP(rec, request)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
			}
		})
	}
}

func TestPermissionsFailure(t *testing.T) {
	user := &User{
		Role: Role{
			Name:        "user",
			Permissions: []string{"read:sounds", "write:sounds"},
		},
	}

	testCases := [][]string{
		{"delete:sounds"},
		{"read:sounds", "delete:sounds"},
		{"read:sounds", "modify:sounds", "write:sounds"},
	}

	for _, tc := range testCases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		handler := createPermTestHandler(user, tc...)
		handler.ServeHTTP(rec, req)

		testutils.CheckJSONResponseError(t, rec, http.StatusForbidden, "you are not authorized to access this resource")
	}
}

//endregion
