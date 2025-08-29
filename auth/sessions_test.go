package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/sessions"
)

func TestGetSession_ValidRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	session, err := GetSession(req)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if session == nil {
		t.Error("Expected session, got nil")
	}

	if session.Name() != "splitweb-session" {
		t.Errorf("Expected session name 'splitweb-session', got %s", session.Name())
	}

	t.Log("GetSession successfully created session")
}

func TestSaveSession_ValidSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	session, err := GetSession(req)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	session.Values["test_key"] = "test_value"

	err = SaveSession(req, w, session)
	if err != nil {
		t.Errorf("Expected no error saving session, got %v", err)
	}

	// Check that a Set-Cookie header was added
	cookies := w.Header()["Set-Cookie"]
	if len(cookies) == 0 {
		t.Error("Expected Set-Cookie header to be set")
	}

	t.Log("SaveSession successfully saved session")
}

func TestSetUserSession_ValidData(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	token := "test-token-123"
	userId := "test-user-456"

	err := SetUserSession(req, w, token, userId)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the session was set by getting it back
	retrievedToken, retrievedUserId, authenticated := GetUserFromSession(req)

	if !authenticated {
		t.Error("Expected user to be authenticated after setting session")
	}

	if retrievedToken != token {
		t.Errorf("Expected token %s, got %s", token, retrievedToken)
	}

	if retrievedUserId != userId {
		t.Errorf("Expected userId %s, got %s", userId, retrievedUserId)
	}

	t.Log("SetUserSession successfully set user session")
}

func TestSetUserSession_EmptyToken(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	err := SetUserSession(req, w, "", "test-user")
	if err != nil {
		t.Errorf("Expected no error for empty token, got %v", err)
	}

	// Should still be able to retrieve empty token
	token, _, authenticated := GetUserFromSession(req)
	if !authenticated {
		t.Error("Expected user to be authenticated even with empty token")
	}

	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}

	t.Log("SetUserSession handled empty token correctly")
}

func TestSetUserSession_EmptyUserId(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	err := SetUserSession(req, w, "test-token", "")
	if err != nil {
		t.Errorf("Expected no error for empty user ID, got %v", err)
	}

	// Should still be able to retrieve empty user ID
	_, userId, authenticated := GetUserFromSession(req)
	if !authenticated {
		t.Error("Expected user to be authenticated even with empty user ID")
	}

	if userId != "" {
		t.Errorf("Expected empty user ID, got %s", userId)
	}

	t.Log("SetUserSession handled empty user ID correctly")
}

func TestClearSession_ValidSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// First set a session
	err := SetUserSession(req, w, "test-token", "test-user")
	if err != nil {
		t.Fatalf("Failed to set session: %v", err)
	}

	// Verify session exists
	_, _, authenticated := GetUserFromSession(req)
	if !authenticated {
		t.Fatal("Expected user to be authenticated before clearing")
	}

	// Clear the session
	err = ClearSession(req, w)
	if err != nil {
		t.Errorf("Expected no error clearing session, got %v", err)
	}

	// Verify session is cleared
	_, _, stillAuthenticated := GetUserFromSession(req)
	if stillAuthenticated {
		t.Error("Expected user to no longer be authenticated after clearing session")
	}

	t.Log("ClearSession successfully cleared session")
}

func TestGetUserFromSession_NoSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	token, userId, authenticated := GetUserFromSession(req)

	if authenticated {
		t.Error("Expected user not to be authenticated without session")
	}

	if token != "" {
		t.Errorf("Expected empty token, got %s", token)
	}

	if userId != "" {
		t.Errorf("Expected empty userId, got %s", userId)
	}

	t.Log("GetUserFromSession correctly handled no session")
}

func TestGetUserFromSession_ValidSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	expectedToken := "test-token-789"
	expectedUserId := "test-user-101112"

	// Set up session
	err := SetUserSession(req, w, expectedToken, expectedUserId)
	if err != nil {
		t.Fatalf("Failed to set session: %v", err)
	}

	// Get session data
	token, userId, authenticated := GetUserFromSession(req)

	if !authenticated {
		t.Error("Expected user to be authenticated with valid session")
	}

	if token != expectedToken {
		t.Errorf("Expected token %s, got %s", expectedToken, token)
	}

	if userId != expectedUserId {
		t.Errorf("Expected userId %s, got %s", expectedUserId, userId)
	}

	t.Log("GetUserFromSession correctly retrieved session data")
}

func TestIsAuthenticated_NoSession(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	authenticated := IsAuthenticated(req)

	if authenticated {
		t.Error("Expected IsAuthenticated to return false without session")
	}

	t.Log("IsAuthenticated correctly returned false for no session")
}

func TestIsAuthenticated_InvalidToken(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	// Set up session with invalid token
	err := SetUserSession(req, w, "invalid-token", "test-user")
	if err != nil {
		t.Fatalf("Failed to set session: %v", err)
	}

	// IsAuthenticated should return false because ValidateToken will fail
	authenticated := IsAuthenticated(req)

	if authenticated {
		t.Error("Expected IsAuthenticated to return false for invalid token")
	}

	t.Log("IsAuthenticated correctly returned false for invalid token")
}

func TestSessionStore_Configuration(t *testing.T) {
	// Test session store configuration
	if store == nil {
		t.Fatal("Expected store to be initialized")
	}

	options := store.Options
	if options.Path != "/" {
		t.Errorf("Expected Path '/', got '%s'", options.Path)
	}

	if options.MaxAge != 86400*7 {
		t.Errorf("Expected MaxAge %d, got %d", 86400*7, options.MaxAge)
	}

	if !options.HttpOnly {
		t.Error("Expected HttpOnly to be true")
	}

	if options.SameSite != http.SameSiteLaxMode {
		t.Errorf("Expected SameSite %d, got %d", http.SameSiteLaxMode, options.SameSite)
	}

	t.Log("Session store configuration is correct")
}

func TestSessionStore_WithCustomSecret(t *testing.T) {
	// Test with custom session secret
	originalSecret := os.Getenv("SESSION_SECRET")
	customSecret := "test-custom-secret-key"
	os.Setenv("SESSION_SECRET", customSecret)

	// Re-initialize store (this is a bit hacky but necessary for testing)
	// In a real scenario, the store would be initialized once at startup
	oldStore := store
	store = sessions.NewCookieStore([]byte(customSecret))

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	err := SetUserSession(req, w, "test-token", "test-user")
	if err != nil {
		t.Errorf("Expected no error with custom secret, got %v", err)
	}

	token, userId, authenticated := GetUserFromSession(req)
	if !authenticated {
		t.Error("Expected authentication with custom secret")
	}

	if token != "test-token" {
		t.Errorf("Expected token 'test-token', got %s", token)
	}

	if userId != "test-user" {
		t.Errorf("Expected userId 'test-user', got %s", userId)
	}

	// Restore original store and environment
	store = oldStore
	if originalSecret != "" {
		os.Setenv("SESSION_SECRET", originalSecret)
	} else {
		os.Unsetenv("SESSION_SECRET")
	}

	t.Log("Session store worked correctly with custom secret")
}

func TestSession_MultipleConcurrentRequests(t *testing.T) {
	// Test that sessions work independently for different requests
	req1 := httptest.NewRequest("GET", "http://example.com", nil)
	req2 := httptest.NewRequest("GET", "http://example.com", nil)
	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	// Set different sessions for each request
	err1 := SetUserSession(req1, w1, "token1", "user1")
	err2 := SetUserSession(req2, w2, "token2", "user2")

	if err1 != nil {
		t.Errorf("Error setting session 1: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Error setting session 2: %v", err2)
	}

	// Verify each request has its own session data
	token1, user1, auth1 := GetUserFromSession(req1)
	token2, user2, auth2 := GetUserFromSession(req2)

	if !auth1 || !auth2 {
		t.Error("Expected both requests to be authenticated")
	}

	if token1 != "token1" {
		t.Errorf("Expected token1 'token1', got %s", token1)
	}

	if token2 != "token2" {
		t.Errorf("Expected token2 'token2', got %s", token2)
	}

	if user1 != "user1" {
		t.Errorf("Expected user1 'user1', got %s", user1)
	}

	if user2 != "user2" {
		t.Errorf("Expected user2 'user2', got %s", user2)
	}

	t.Log("Multiple concurrent requests handled correctly")
}

func TestSession_CorruptedData(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com", nil)

	// Manually create a session with corrupted data
	session, err := GetSession(req)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	// Set invalid data types
	session.Values["token"] = 12345            // Should be string
	session.Values["user_id"] = []byte("user") // Should be string
	session.Values["authenticated"] = "true"   // Should be boolean

	// Try to get user from corrupted session
	token, userId, authenticated := GetUserFromSession(req)

	if authenticated {
		t.Error("Expected authentication to fail with corrupted session data")
	}

	if token != "" {
		t.Errorf("Expected empty token with corrupted data, got %s", token)
	}

	if userId != "" {
		t.Errorf("Expected empty userId with corrupted data, got %s", userId)
	}

	t.Log("GetUserFromSession correctly handled corrupted session data")
}
