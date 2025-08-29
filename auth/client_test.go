package auth

import (
	"os"
	"testing"
)

func TestGetPocketBaseURL_ValidURL(t *testing.T) {
	expectedURL := "http://test.example.com:8090"
	os.Setenv("POCKET_BASE_URL", expectedURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	t.Logf("getPocketBaseURL returned correct URL: %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

func TestGetPocketBaseURL_EmptyURL(t *testing.T) {
	// Ensure the environment variable is not set
	os.Unsetenv("POCKET_BASE_URL")

	url, err := getPocketBaseURL()

	if err == nil {
		t.Error("Expected error when POCKET_BASE_URL is not set")
	}

	if url != "" {
		t.Errorf("Expected empty URL, got %s", url)
	}

	expectedError := "POCKET_BASE_URL environment variable not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}

	t.Logf("getPocketBaseURL correctly failed when environment variable not set: %v", err)
}

func TestGetPocketBaseURL_WhitespaceURL(t *testing.T) {
	// Test with whitespace-only URL
	os.Setenv("POCKET_BASE_URL", "   ")

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != "   " {
		t.Errorf("Expected URL '   ', got '%s'", url)
	}

	t.Logf("getPocketBaseURL returned whitespace URL as-is: '%s'", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

func TestGetPocketBaseURL_URLWithPath(t *testing.T) {
	expectedURL := "http://test.example.com:8090/api/v1"
	os.Setenv("POCKET_BASE_URL", expectedURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	t.Logf("getPocketBaseURL handled URL with path correctly: %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

func TestGetPocketBaseURL_HTTPSUrl(t *testing.T) {
	expectedURL := "https://secure.example.com"
	os.Setenv("POCKET_BASE_URL", expectedURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	t.Logf("getPocketBaseURL handled HTTPS URL correctly: %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

func TestGetPocketBaseURL_LocalhostURL(t *testing.T) {
	expectedURL := "http://localhost:8090"
	os.Setenv("POCKET_BASE_URL", expectedURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	t.Logf("getPocketBaseURL handled localhost URL correctly: %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

func TestGetPocketBaseURL_IPAddressURL(t *testing.T) {
	expectedURL := "http://192.168.1.100:8090"
	os.Setenv("POCKET_BASE_URL", expectedURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, url)
	}

	t.Logf("getPocketBaseURL handled IP address URL correctly: %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}

// Test that the function doesn't validate URL format - it just returns what's set
func TestGetPocketBaseURL_InvalidURLFormat(t *testing.T) {
	invalidURL := "not-a-valid-url"
	os.Setenv("POCKET_BASE_URL", invalidURL)

	url, err := getPocketBaseURL()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if url != invalidURL {
		t.Errorf("Expected URL %s, got %s", invalidURL, url)
	}

	t.Logf("getPocketBaseURL returned invalid URL as-is (validation should be done elsewhere): %s", url)

	// Clean up
	os.Unsetenv("POCKET_BASE_URL")
}
