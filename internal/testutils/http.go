package testutils

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func CheckJSONResponseError(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	t.Helper()
	if rec.Code != expectedStatus {
		t.Errorf("Expected status code %d, got %d", expectedStatus, rec.Code)
	}

	var response map[string]interface{}
	var err = json.NewDecoder(rec.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	errorMsg, ok := response["error"]
	if !ok {
		t.Fatalf("Expected 'error' key in response, but it was not found")
	}

	if errorMsg != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, errorMsg)
	}
}
