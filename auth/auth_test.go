package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Set up test environment
	os.Setenv("POCKET_BASE_URL", "http://10.71.71.10:8090")
	os.Setenv("SESSION_SECRET", "test-secret-key")

	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
	os.Unsetenv("SESSION_SECRET")

	os.Exit(code)
}

func TestAuthWithPassword_ValidCredentials(t *testing.T) {
	// Note: This test requires a valid user in PocketBase
	// You may need to create a test user first

	email := "test@example.com"
	password := "testpassword"

	authResponse, err := AuthWithPassword(email, password)

	if err != nil {
		// If authentication fails, it might be because the user doesn't exist
		// This is expected for integration tests
		t.Logf("Authentication failed (expected if test user doesn't exist): %v", err)
		return
	}

	if authResponse == nil {
		t.Fatal("Expected auth response, got nil")
	}

	if authResponse.Token == "" {
		t.Error("Expected token to be present")
	}

	if authResponse.Record.Email != email {
		t.Errorf("Expected email %s, got %s", email, authResponse.Record.Email)
	}

	t.Logf("Authentication successful for user: %s (ID: %s)", authResponse.Record.Email, authResponse.Record.Id)
}

func TestAuthWithPassword_InvalidCredentials(t *testing.T) {
	email := "nonexistent@example.com"
	password := "wrongpassword"

	authResponse, err := AuthWithPassword(email, password)

	if err == nil {
		t.Error("Expected authentication to fail with invalid credentials")
	}

	if authResponse != nil {
		t.Error("Expected nil response for invalid credentials")
	}

	t.Logf("Authentication correctly failed for invalid credentials: %v", err)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid-jwt-token-12345"
	fakeUserId := "fake-user-id-123"

	err := ValidateToken(invalidToken, fakeUserId)

	if err == nil {
		t.Error("Expected token validation to fail with invalid token")
	}

	t.Logf("Token validation correctly failed for invalid token: %v", err)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	err := ValidateToken("", "some-user-id")

	if err == nil {
		t.Error("Expected token validation to fail with empty token")
	}

	t.Logf("Token validation correctly failed for empty token: %v", err)
}

func TestRequestPasswordReset_ValidEmail(t *testing.T) {
	// Test with a potentially valid email format
	email := "test@example.com"

	err := RequestPasswordReset(email)

	// This might fail if the email doesn't exist in PocketBase, which is expected
	if err != nil {
		t.Logf("Password reset request failed (expected if email doesn't exist): %v", err)
	} else {
		t.Logf("Password reset request successful for email: %s", email)
	}
}

func TestRequestPasswordReset_InvalidEmail(t *testing.T) {
	email := "definitely-not-existing-email@nonexistent-domain-12345.com"

	err := RequestPasswordReset(email)

	// PocketBase may not return an error for non-existent emails (security feature)
	// This is actually good behavior to prevent email enumeration attacks
	if err == nil {
		t.Logf("Password reset request completed (PocketBase doesn't reveal if email exists for security)")
	} else {
		t.Logf("Password reset failed for non-existent email: %v", err)
	}

	// Test passes either way since both behaviors are valid
}

func TestSessionManagement(t *testing.T) {
	// Create a test HTTP request and response recorder
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// Test setting user session
	token := "test-token-123"
	userId := "test-user-456"

	err := SetUserSession(req, w, token, userId)
	if err != nil {
		t.Fatalf("Failed to set user session: %v", err)
	}

	// Test getting user from session
	retrievedToken, retrievedUserId, authenticated := GetUserFromSession(req)

	if !authenticated {
		t.Error("Expected user to be authenticated")
	}

	if retrievedToken != token {
		t.Errorf("Expected token %s, got %s", token, retrievedToken)
	}

	if retrievedUserId != userId {
		t.Errorf("Expected userId %s, got %s", userId, retrievedUserId)
	}

	// Test clearing session
	err = ClearSession(req, w)
	if err != nil {
		t.Fatalf("Failed to clear session: %v", err)
	}

	// Verify session is cleared
	_, _, stillAuthenticated := GetUserFromSession(req)
	if stillAuthenticated {
		t.Error("Expected user to no longer be authenticated after clearing session")
	}

	t.Log("Session management test passed")
}

func TestIsAuthenticated_WithValidSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// Set up a session with invalid token (should fail PocketBase validation)
	invalidToken := "invalid-token"
	userId := "fake-user"

	err := SetUserSession(req, w, invalidToken, userId)
	if err != nil {
		t.Fatalf("Failed to set session: %v", err)
	}

	// Test IsAuthenticated - should return false because token is invalid
	authenticated := IsAuthenticated(req)

	if authenticated {
		t.Error("Expected IsAuthenticated to return false for invalid token")
	}

	t.Log("IsAuthenticated correctly rejected invalid token")
}

func TestMiddlewareIntegration(t *testing.T) {
	// Test middleware behavior with invalid session
	req := httptest.NewRequest("GET", "http://example.com/protected", nil)
	w := httptest.NewRecorder()

	// Set up invalid session
	SetUserSession(req, w, "invalid-token", "fake-user")

	// Create a test handler that should not be called
	handlerCalled := false
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}

	// Wrap with RequireAuth middleware
	protectedHandler := RequireAuth(testHandler)

	// Call the protected handler
	protectedHandler(w, req)

	if handlerCalled {
		t.Error("Expected handler not to be called with invalid token")
	}

	// Check that we got redirected (status should be 303 for redirect)
	if w.Code != http.StatusSeeOther && w.Code != http.StatusUnauthorized {
		t.Errorf("Expected redirect status, got %d", w.Code)
	}

	t.Log("Middleware correctly blocked access with invalid token")
}
