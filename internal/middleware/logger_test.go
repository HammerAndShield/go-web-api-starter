package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer
	writer := io.Writer(&buf)
	logger := slog.New(slog.NewJSONHandler(writer, nil))

	logMiddleware := Logger(logger, []string{})

	// Create a new handler to test the middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello, World!"))
	})

	// Wrap the handler with the middleware
	handler := logMiddleware(testHandler)

	// Create a new HTTP request
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Header.Set("X-Request-Id", "test-request-id")
	req = req.WithContext(context.WithValue(req.Context(), "user", "test-user"))

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler with the ResponseRecorder and request
	handler.ServeHTTP(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedBody := "Hello, World!"
	body, _ := io.ReadAll(rr.Body)
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}

	// Check if log output contains expected content
	logOutput := buf.String()
	if !strings.Contains(logOutput, fmt.Sprintf("%s: test-request-id", RequestIdLog)) {
		t.Errorf("Expected log entry to contain request id %q, but it wasn't found", "test-request-id")
	}
}
