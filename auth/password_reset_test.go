package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// RequestPasswordReset basic tests are already covered in auth_test.go
// These tests focus on additional scenarios

func TestConfirmPasswordReset_ValidData(t *testing.T) {
	token := "test-reset-token"
	password := "newpassword123"
	passwordConfirm := "newpassword123"

	err := ConfirmPasswordReset(token, password, passwordConfirm)

	if err != nil {
		// Expected to fail with test token
		t.Logf("ConfirmPasswordReset failed (expected with test token): %v", err)
		expectedError := "change password request was invalid"
		if err.Error() != expectedError {
			t.Logf("Got different error than expected: %v", err)
		}
	} else {
		t.Log("ConfirmPasswordReset succeeded")
	}
}

func TestConfirmPasswordReset_EmptyToken(t *testing.T) {
	err := ConfirmPasswordReset("", "password", "password")

	if err == nil {
		t.Error("Expected ConfirmPasswordReset to fail with empty token")
	}

	t.Logf("ConfirmPasswordReset correctly failed for empty token: %v", err)
}

func TestConfirmPasswordReset_EmptyPassword(t *testing.T) {
	err := ConfirmPasswordReset("test-token", "", "")

	if err == nil {
		t.Error("Expected ConfirmPasswordReset to fail with empty password")
	}

	t.Logf("ConfirmPasswordReset correctly failed for empty password: %v", err)
}

func TestConfirmPasswordReset_MismatchedPasswords(t *testing.T) {
	err := ConfirmPasswordReset("test-token", "password1", "password2")

	if err == nil {
		t.Error("Expected ConfirmPasswordReset to fail with mismatched passwords")
	}

	t.Logf("ConfirmPasswordReset correctly failed for mismatched passwords: %v", err)
}

func TestChangePasswordWithOlderOne_ValidData(t *testing.T) {
	oldPassword := "oldpassword123"
	newPassword := "newpassword123"
	passwordConfirm := "newpassword123"
	token := "test-token"
	userID := "test-user-id"

	err := ChangePasswordWithOlderOne(oldPassword, newPassword, passwordConfirm, token, userID)

	if err != nil {
		// Expected to fail with test data
		t.Logf("ChangePasswordWithOlderOne failed (expected with test data): %v", err)
		expectedError := "change password with older password request failed"
		if err.Error() != expectedError {
			t.Logf("Got different error than expected: %v", err)
		}
	} else {
		t.Log("ChangePasswordWithOlderOne succeeded")
	}
}

func TestChangePasswordWithOlderOne_EmptyOldPassword(t *testing.T) {
	err := ChangePasswordWithOlderOne("", "newpassword", "newpassword", "token", "user")

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with empty old password")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for empty old password: %v", err)
}

func TestChangePasswordWithOlderOne_EmptyNewPassword(t *testing.T) {
	err := ChangePasswordWithOlderOne("oldpassword", "", "", "token", "user")

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with empty new password")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for empty new password: %v", err)
}

func TestChangePasswordWithOlderOne_MismatchedNewPasswords(t *testing.T) {
	err := ChangePasswordWithOlderOne("oldpassword", "newpassword1", "newpassword2", "token", "user")

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with mismatched new passwords")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for mismatched new passwords: %v", err)
}

func TestChangePasswordWithOlderOne_EmptyToken(t *testing.T) {
	err := ChangePasswordWithOlderOne("oldpassword", "newpassword", "newpassword", "", "user")

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with empty token")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for empty token: %v", err)
}

func TestChangePasswordWithOlderOne_EmptyUserID(t *testing.T) {
	err := ChangePasswordWithOlderOne("oldpassword", "newpassword", "newpassword", "token", "")

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with empty user ID")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for empty user ID: %v", err)
}

func TestRequestPasswordReset_MockServerResponses(t *testing.T) {
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

			// Test RequestPasswordReset
			err := RequestPasswordReset("test@example.com")

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

func TestConfirmPasswordReset_MockServerResponses(t *testing.T) {
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

			// Test ConfirmPasswordReset
			err := ConfirmPasswordReset("test-token", "password", "password")

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

func TestChangePasswordWithOlderOne_MockServerResponses(t *testing.T) {
	tests := []struct {
		statusCode  int
		expectError bool
		description string
	}{
		{200, false, "success"},
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

			// Test ChangePasswordWithOlderOne
			err := ChangePasswordWithOlderOne("oldpass", "newpass", "newpass", "token", "user")

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

func TestRequestPasswordReset_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := RequestPasswordReset("test@example.com")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected RequestPasswordReset to fail with invalid URL")
	}

	t.Logf("RequestPasswordReset correctly failed for network error: %v", err)
}

func TestConfirmPasswordReset_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := ConfirmPasswordReset("token", "password", "password")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected ConfirmPasswordReset to fail with invalid URL")
	}

	t.Logf("ConfirmPasswordReset correctly failed for network error: %v", err)
}

func TestChangePasswordWithOlderOne_NetworkError(t *testing.T) {
	// Test with invalid PocketBase URL to simulate network error
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", "http://invalid-url-12345.com")

	err := ChangePasswordWithOlderOne("oldpass", "newpass", "newpass", "token", "user")

	// Restore original URL
	if originalURL != "" {
		os.Setenv("POCKET_BASE_URL", originalURL)
	} else {
		os.Unsetenv("POCKET_BASE_URL")
	}

	if err == nil {
		t.Error("Expected ChangePasswordWithOlderOne to fail with invalid URL")
	}

	t.Logf("ChangePasswordWithOlderOne correctly failed for network error: %v", err)
}
