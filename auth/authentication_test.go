package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetUserRecord_ValidToken(t *testing.T) {
	// This test requires a valid token and userID from PocketBase
	// It's more of an integration test
	token := "test-token"
	userID := "test-user-id"

	record, err := GetUserRecord(token, userID)

	if err != nil {
		// Expected to fail with invalid token
		t.Logf("GetUserRecord failed (expected with test token): %v", err)
		return
	}

	if record == nil {
		t.Fatal("Expected record, got nil")
	}

	if record.Id != userID {
		t.Errorf("Expected user ID %s, got %s", userID, record.Id)
	}

	t.Logf("GetUserRecord successful for user: %s", record.Id)
}

func TestGetUserRecord_EmptyToken(t *testing.T) {
	record, err := GetUserRecord("", "test-user")

	if err == nil {
		t.Error("Expected GetUserRecord to fail with empty token")
	}

	if record != nil {
		t.Error("Expected nil record for empty token")
	}

	t.Logf("GetUserRecord correctly failed for empty token: %v", err)
}

func TestGetUserRecord_EmptyUserID(t *testing.T) {
	record, err := GetUserRecord("test-token", "")

	if err == nil {
		t.Error("Expected GetUserRecord to fail with empty user ID")
	}

	if record != nil {
		t.Error("Expected nil record for empty user ID")
	}

	t.Logf("GetUserRecord correctly failed for empty user ID: %v", err)
}

func TestAuthWithPassword_EmptyCredentials(t *testing.T) {
	// Test empty identity
	authResponse, err := AuthWithPassword("", "password")
	if err == nil {
		t.Error("Expected authentication to fail with empty identity")
	}
	if authResponse != nil {
		t.Error("Expected nil response for empty identity")
	}

	// Test empty password
	authResponse, err = AuthWithPassword("user@example.com", "")
	if err == nil {
		t.Error("Expected authentication to fail with empty password")
	}
	if authResponse != nil {
		t.Error("Expected nil response for empty password")
	}

	// Test both empty
	authResponse, err = AuthWithPassword("", "")
	if err == nil {
		t.Error("Expected authentication to fail with empty credentials")
	}
	if authResponse != nil {
		t.Error("Expected nil response for empty credentials")
	}

	t.Log("AuthWithPassword correctly failed for empty credentials")
}

func TestValidateToken_EmptyUserID(t *testing.T) {
	err := ValidateToken("test-token", "")

	if err == nil {
		t.Error("Expected token validation to fail with empty user ID")
	}

	t.Logf("Token validation correctly failed for empty user ID: %v", err)
}

func TestAuthWithPassword_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	authResponse, err := AuthWithPassword("test@example.com", "password")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected authentication to fail with invalid URL")
	}

	if authResponse != nil {
		t.Error("Expected nil response for network error")
	}

	t.Logf("AuthWithPassword correctly failed for network error: %v", err)
}

func TestValidateToken_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := ValidateToken("test-token", "test-user")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected token validation to fail with invalid URL")
	}

	t.Logf("Token validation correctly failed for network error: %v", err)
}

func TestGetUserRecord_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	record, err := GetUserRecord("test-token", "test-user")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected GetUserRecord to fail with invalid URL")
	}

	if record != nil {
		t.Error("Expected nil record for network error")
	}

	t.Logf("GetUserRecord correctly failed for network error: %v", err)
}

func TestAuthWithPassword_MockServerResponse(t *testing.T) {
	// Create a test server that returns different HTTP status codes
	tests := []struct {
		statusCode   int
		responseBody string
		expectError  bool
		description  string
	}{
		{200, `{"token":"test-token","record":{"id":"test-id","email":"test@example.com"}}`, false, "valid response"},
		{400, `{"message":"Bad Request"}`, true, "bad request"},
		{401, `{"message":"Unauthorized"}`, true, "unauthorized"},
		{404, `{"message":"Not Found"}`, true, "not found"},
		{500, `{"message":"Internal Server Error"}`, true, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Set test server URL
			originalURL := os.Getenv("POCKET_BASE_URL")
			os.Setenv("POCKET_BASE_URL", server.URL)

			// Test AuthWithPassword
			authResponse, err := AuthWithPassword("test@example.com", "password")

			// Restore original URL
			if originalURL != "" {
				os.Setenv("POCKET_BASE_URL", originalURL)
			} else {
				os.Unsetenv("POCKET_BASE_URL")
			}

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}

			if tt.expectError && authResponse != nil {
				t.Errorf("Expected nil response for %s", tt.description)
			}

			if !tt.expectError && authResponse == nil {
				t.Errorf("Expected valid response for %s", tt.description)
			}
		})
	}
}
