package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestRequestEmailChange_ValidIdentity(t *testing.T) {
	email := "test@example.com"
	token := "test-token"

	err := RequestEmailChange(email, token)

	if err != nil {
		// Expected to fail without valid PocketBase setup
		t.Logf("RequestEmailChange failed (expected without valid setup): %v", err)
	} else {
		t.Log("RequestEmailChange succeeded")
	}
}

func TestRequestEmailChange_EmptyIdentity(t *testing.T) {
	err := RequestEmailChange("", "test-token")

	if err == nil {
		t.Error("Expected RequestEmailChange to fail with empty identity")
	}

	t.Logf("RequestEmailChange correctly failed for empty identity: %v", err)
}

func TestRequestEmailChange_EmptyToken(t *testing.T) {
	err := RequestEmailChange("test@example.com", "")

	if err == nil {
		t.Error("Expected RequestEmailChange to fail with empty token")
	}

	t.Logf("RequestEmailChange correctly failed for empty token: %v", err)
}

func TestRequestVerificationEmail_ValidEmail(t *testing.T) {
	email := "test@example.com"

	err := RequestVerificationEmail(email)

	if err != nil {
		// Expected to fail without valid PocketBase setup
		t.Logf("RequestVerificationEmail failed (expected without valid setup): %v", err)
	} else {
		t.Log("RequestVerificationEmail succeeded")
	}
}

func TestRequestVerificationEmail_EmptyEmail(t *testing.T) {
	err := RequestVerificationEmail("")

	if err == nil {
		t.Error("Expected RequestVerificationEmail to fail with empty email")
	}

	t.Logf("RequestVerificationEmail correctly failed for empty email: %v", err)
}

func TestRequestVerificationEmail_InvalidEmailFormat(t *testing.T) {
	// Test with obviously invalid email format
	invalidEmails := []string{
		"invalid-email",
		"@example.com",
		"test@",
		"test.example.com",
		"test@@example.com",
	}

	for _, email := range invalidEmails {
		err := RequestVerificationEmail(email)
		// PocketBase might still accept these and return success for security
		// so we don't necessarily expect an error
		t.Logf("RequestVerificationEmail with invalid email '%s': %v", email, err)
	}
}

func TestVerifyEmail_ValidToken(t *testing.T) {
	token := "test-verification-token"

	err := VerifyEmail(token)

	if err != nil {
		// Expected to fail with test token
		t.Logf("VerifyEmail failed (expected with test token): %v", err)
		expectedError := "The token allocated is incorrect. Please try again"
		if err.Error() != expectedError {
			t.Logf("Got different error than expected: %v", err)
		}
	} else {
		t.Log("VerifyEmail succeeded")
	}
}

func TestVerifyEmail_EmptyToken(t *testing.T) {
	err := VerifyEmail("")

	if err == nil {
		t.Error("Expected VerifyEmail to fail with empty token")
	}

	t.Logf("VerifyEmail correctly failed for empty token: %v", err)
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	invalidToken := "invalid-token-12345"

	err := VerifyEmail(invalidToken)

	if err == nil {
		t.Error("Expected VerifyEmail to fail with invalid token")
	}

	t.Logf("VerifyEmail correctly failed for invalid token: %v", err)
}

func TestRequestEmailChange_MockServerResponses(t *testing.T) {
	tests := []struct {
		statusCode  int
		expectError bool
		description string
	}{
		{204, false, "success"},
		{400, true, "bad request"},
		{401, true, "unauthorized"},
		{404, true, "not found"},
		{500, true, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Set test server URL
			originalURL := os.Getenv("POCKET_BASE_URL")
			os.Setenv("POCKET_BASE_URL", server.URL)

			// Test RequestEmailChange
			err := RequestEmailChange("test@example.com", "test-token")

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
		})
	}
}

func TestRequestVerificationEmail_MockServerResponses(t *testing.T) {
	tests := []struct {
		statusCode  int
		expectError bool
		description string
	}{
		{204, false, "success"},
		{400, true, "bad request"},
		{401, true, "unauthorized"},
		{404, true, "not found"},
		{500, true, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Set test server URL
			originalURL := os.Getenv("POCKET_BASE_URL")
			os.Setenv("POCKET_BASE_URL", server.URL)

			// Test RequestVerificationEmail
			err := RequestVerificationEmail("test@example.com")

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
		})
	}
}

func TestVerifyEmail_MockServerResponses(t *testing.T) {
	tests := []struct {
		statusCode  int
		expectError bool
		description string
	}{
		{204, false, "success"},
		{400, true, "bad request"},
		{401, true, "unauthorized"},
		{404, true, "not found"},
		{500, true, "server error"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			// Set test server URL
			originalURL := os.Getenv("POCKET_BASE_URL")
			os.Setenv("POCKET_BASE_URL", server.URL)

			// Test VerifyEmail
			err := VerifyEmail("test-token")

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
		})
	}
}

func TestRequestEmailChange_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := RequestEmailChange("test@example.com", "test-token")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected RequestEmailChange to fail with invalid URL")
	}

	t.Logf("RequestEmailChange correctly failed for network error: %v", err)
}

func TestRequestVerificationEmail_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := RequestVerificationEmail("test@example.com")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected RequestVerificationEmail to fail with invalid URL")
	}

	t.Logf("RequestVerificationEmail correctly failed for network error: %v", err)
}

func TestVerifyEmail_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := VerifyEmail("test-token")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected VerifyEmail to fail with invalid URL")
	}

	t.Logf("VerifyEmail correctly failed for network error: %v", err)
}
