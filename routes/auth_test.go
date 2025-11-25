package routes

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/alnav3/splitweb/auth"
)

func TestMain(m *testing.M) {
	// Set up test environment - using localhost to avoid external calls
	os.Setenv("POCKET_BASE_URL", "http://localhost:99999") // Non-existent port for safety
	os.Setenv("SESSION_SECRET", "test-secret-key")

	// Run tests
	code := m.Run()

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
	os.Unsetenv("SESSION_SECRET")

	os.Exit(code)
}

func TestAuthLoginHandler_InvalidCredentials(t *testing.T) {
	// Create form data with invalid credentials
	form := url.Values{}
	form.Add("email", "nonexistent@example.com")
	form.Add("password", "wrongpassword")

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthLoginHandler(w, req)

	// Should return 401 Unauthorized for invalid credentials
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Should contain error message
	body := w.Body.String()
	if !strings.Contains(body, "Invalid email or password") {
		t.Error("Expected error message in response body")
	}

	t.Log("Login handler correctly rejected invalid credentials")
}

func TestAuthLoginHandler_EmptyCredentials(t *testing.T) {
	// Create form data with empty credentials
	form := url.Values{}
	form.Add("email", "")
	form.Add("password", "")

	req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthLoginHandler(w, req)

	// Should return 400 Bad Request for empty credentials
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Should contain validation error
	body := w.Body.String()
	if !strings.Contains(body, "Email and password are required") {
		t.Error("Expected validation error in response body")
	}

	t.Log("Login handler correctly rejected empty credentials")
}

func TestAuthRegisterHandler_InvalidData(t *testing.T) {
	// Create form data with password mismatch
	form := url.Values{}
	form.Add("name", "Test User")
	form.Add("email", "test@example.com")
	form.Add("password", "password123")
	form.Add("confirm-password", "password456") // Different password

	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthRegisterHandler(w, req)

	// Should return 400 Bad Request for mismatched passwords
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Should contain validation error
	body := w.Body.String()
	if !strings.Contains(body, "Passwords do not match") {
		t.Error("Expected password mismatch error in response body")
	}

	t.Log("Register handler correctly rejected mismatched passwords")
}

func TestAuthRegisterHandler_MissingFields(t *testing.T) {
	// Create form data with missing fields
	form := url.Values{}
	form.Add("name", "")
	form.Add("email", "test@example.com")
	form.Add("password", "password123")
	form.Add("confirm-password", "password123")

	req := httptest.NewRequest("POST", "/auth/register", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthRegisterHandler(w, req)

	// Should return 400 Bad Request for missing name
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Should contain validation error
	body := w.Body.String()
	if !strings.Contains(body, "All fields are required") {
		t.Error("Expected missing fields error in response body")
	}

	t.Log("Register handler correctly rejected missing fields")
}

func TestAuthForgotPasswordHandler_EmptyEmail(t *testing.T) {
	// Create form data with empty email
	form := url.Values{}
	form.Add("email", "")

	req := httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthForgotPasswordHandler(w, req)

	// Should return 400 Bad Request for empty email
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Should contain validation error
	body := w.Body.String()
	if !strings.Contains(body, "Email is required") {
		t.Error("Expected email required error in response body")
	}

	t.Log("Forgot password handler correctly rejected empty email")
}

func TestAuthForgotPasswordHandler_NonExistentEmail(t *testing.T) {
	// Create a mock server that simulates PocketBase behavior 
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate PocketBase behavior - returns 204 even for non-existent emails (security)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Set the mock server URL
	originalURL := os.Getenv("POCKET_BASE_URL")
	os.Setenv("POCKET_BASE_URL", server.URL)
	defer func() {
		if originalURL != "" {
			os.Setenv("POCKET_BASE_URL", originalURL)
		} else {
			os.Unsetenv("POCKET_BASE_URL")
		}
	}()

	// Create form data with non-existent email
	form := url.Values{}
	form.Add("email", "definitely-not-existing@nonexistent-domain-12345.com")

	req := httptest.NewRequest("POST", "/auth/forgot-password", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	AuthForgotPasswordHandler(w, req)

	// Should return 200 OK (successful form processing)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	t.Logf("Forgot password handler response: %d", w.Code)
}

func TestAuthLogoutHandler(t *testing.T) {
	req := httptest.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()

	// Set up a session first
	auth.SetUserSession(req, w, "test-token", "test-user")

	AuthLogoutHandler(w, req)

	// Should redirect after logout
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Should have HX-Redirect header
	redirect := w.Header().Get("HX-Redirect")
	if redirect != "/login" {
		t.Errorf("Expected HX-Redirect to '/login', got '%s'", redirect)
	}

	t.Log("Logout handler correctly processed logout")
}

func TestMiddlewareProtection(t *testing.T) {
	// Test protected route without authentication
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	// Create a dummy handler
	protectedHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Protected content"))
	}

	// Wrap with auth middleware
	wrappedHandler := auth.RequireAuth(protectedHandler)

	// Call the handler
	wrappedHandler(w, req)

	// Should redirect to login
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
	}

	// Check Location header
	location := w.Header().Get("Location")
	if location != "/login" {
		t.Errorf("Expected redirect to '/login', got '%s'", location)
	}

	t.Log("Middleware correctly protected route without authentication")
}

func TestMiddlewareWithInvalidToken(t *testing.T) {
	// Test protected route with invalid token
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()

	// Set invalid session
	auth.SetUserSession(req, w, "invalid-token", "fake-user")

	// Create a dummy handler that should not be called
	handlerCalled := false
	protectedHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Protected content"))
	}

	// Wrap with auth middleware
	wrappedHandler := auth.RequireAuth(protectedHandler)

	// Call the handler
	wrappedHandler(w, req)

	// Handler should not have been called
	if handlerCalled {
		t.Error("Expected protected handler not to be called with invalid token")
	}

	// Should redirect to login
	if w.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect status %d, got %d", http.StatusSeeOther, w.Code)
	}

	t.Log("Middleware correctly rejected invalid token and redirected to login")
}

func TestRedirectIfAuthenticated_WithoutSession(t *testing.T) {
	// Test auth page without session (should allow access)
	req := httptest.NewRequest("GET", "/login", nil)
	w := httptest.NewRecorder()

	// Create login handler
	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login page"))
	}

	// Wrap with redirect middleware
	wrappedHandler := auth.RedirectIfAuthenticated(loginHandler)

	// Call the handler
	wrappedHandler(w, req)

	// Should allow access to login page
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Login page") {
		t.Error("Expected login page content")
	}

	t.Log("RedirectIfAuthenticated correctly allowed access without session")
}
