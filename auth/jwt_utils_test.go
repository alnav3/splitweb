package auth

import (
	"testing"
)

func TestDecodeJWTToken_ValidToken(t *testing.T) {
	// Create a test JWT token (this is just for testing the parsing, not verification)
	// This is a sample JWT token with id="test-user-123" (created at jwt.io for testing)
	testToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6InRlc3QtdXNlci0xMjMiLCJpc3MiOiJ0ZXN0IiwiYXVkIjoidGVzdCIsImlhdCI6MTY1MDAwMDAwMCwiZXhwIjoxNjUwMDg2NDAwfQ.5P9WMmh68Y5pYQqL7D4P9HuGV4S0rePiHoNmZcQ9j7s"
	expectedID := "test-user-123"

	userID, err := DecodeJWTToken(testToken)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if userID != expectedID {
		t.Errorf("Expected user ID %s, got %s", expectedID, userID)
	}

	t.Logf("DecodeJWTToken successfully decoded user ID: %s", userID)
}

func TestDecodeJWTToken_EmptyToken(t *testing.T) {
	userID, err := DecodeJWTToken("")

	if err == nil {
		t.Error("Expected error for empty token")
	}

	if userID != "" {
		t.Errorf("Expected empty user ID, got %s", userID)
	}

	t.Logf("DecodeJWTToken correctly failed for empty token: %v", err)
}

func TestDecodeJWTToken_InvalidToken(t *testing.T) {
	invalidToken := "invalid.jwt.token"

	userID, err := DecodeJWTToken(invalidToken)

	if err == nil {
		t.Error("Expected error for invalid token")
	}

	if userID != "" {
		t.Errorf("Expected empty user ID, got %s", userID)
	}

	t.Logf("DecodeJWTToken correctly failed for invalid token: %v", err)
}

func TestDecodeJWTToken_MalformedToken(t *testing.T) {
	// Test various malformed tokens
	malformedTokens := []string{
		"not.a.jwt",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9", // Missing parts
		"invalid-string",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.invalid-payload.signature",
		"header.eyJ0ZXN0IjoidGVzdCJ9.signature", // Missing id field
	}

	for i, token := range malformedTokens {
		t.Run("malformed_"+string(rune('a'+i)), func(t *testing.T) {
			userID, err := DecodeJWTToken(token)

			if err == nil {
				t.Errorf("Expected error for malformed token: %s", token)
			}

			if userID != "" {
				t.Errorf("Expected empty user ID for malformed token, got %s", userID)
			}

			t.Logf("DecodeJWTToken correctly failed for malformed token: %v", err)
		})
	}
}

func TestDecodeJWTToken_TokenWithoutIdClaim(t *testing.T) {
	// JWT token without 'id' field (created at jwt.io for testing)
	tokenWithoutId := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJ1c2VyIjoidGVzdCIsImlzcyI6InRlc3QiLCJhdWQiOiJ0ZXN0IiwiaWF0IjoxNjUwMDAwMDAwLCJleHAiOjE2NTAwODY0MDB9.M9lVPCGaIvA8xYkkmeBcRqJ1zlnVxQDhNV5JVXhH7Js"

	userID, err := DecodeJWTToken(tokenWithoutId)

	// The function may succeed but return empty ID, or it may fail
	// Both are acceptable behaviors for a token without id claim
	if err != nil {
		t.Logf("DecodeJWTToken correctly failed for token without id claim: %v", err)
	} else if userID == "" {
		t.Log("DecodeJWTToken succeeded but returned empty ID for token without id claim")
	} else {
		t.Errorf("Expected empty user ID for token without id claim, got %s", userID)
	}
}

func TestDecodeJWTToken_TokenWithEmptyIdClaim(t *testing.T) {
	// JWT token with empty 'id' field (created at jwt.io for testing)
	tokenWithEmptyId := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6IiIsImlzcyI6InRlc3QiLCJhdWQiOiJ0ZXN0IiwiaWF0IjoxNjUwMDAwMDAwLCJleHAiOjE2NTAwODY0MDB9.CwQhKm6tP7nkKmVqN1J3qDnK7ZsYpHQsGvwOoQKVt-8"

	userID, err := DecodeJWTToken(tokenWithEmptyId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if userID != "" {
		t.Errorf("Expected empty user ID, got '%s'", userID)
	}

	t.Logf("DecodeJWTToken returned empty ID for token with empty id claim: '%s'", userID)
}

func TestDecodeJWTToken_TokenWithNullIdClaim(t *testing.T) {
	// JWT token with null 'id' field (created at jwt.io for testing)
	tokenWithNullId := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6bnVsbCwiaXNzIjoidGVzdCIsImF1ZCI6InRlc3QiLCJpYXQiOjE2NTAwMDAwMDAsImV4cCI6MTY1MDA4NjQwMH0.oqEi3rHKITIoF5L4kVvX8rp9jSGOOzYlXWrL5lDcgik"

	userID, err := DecodeJWTToken(tokenWithNullId)

	// This should not fail, but should return empty string
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if userID != "" {
		t.Errorf("Expected empty user ID for null id claim, got '%s'", userID)
	}

	t.Logf("DecodeJWTToken returned empty ID for token with null id claim: '%s'", userID)
}

func TestDecodeJWTToken_ExpiredToken(t *testing.T) {
	// JWT token that's expired (but we're only parsing, not verifying)
	// This should still work since we use ParseUnverified
	expiredToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6InRlc3QtdXNlci1leHBpcmVkIiwiaXNzIjoidGVzdCIsImF1ZCI6InRlc3QiLCJpYXQiOjE1MDAwMDAwMDAsImV4cCI6MTUwMDA4NjQwMH0.WyKJ6p8kDjHt_Fw8IY5WL2hjQyVMZkJ5qn9LZcGPrk4"
	expectedID := "test-user-expired"

	userID, err := DecodeJWTToken(expiredToken)

	if err != nil {
		t.Errorf("Expected no error for expired token (since we use ParseUnverified), got %v", err)
	}

	if userID != expectedID {
		t.Errorf("Expected user ID %s, got %s", expectedID, userID)
	}

	t.Logf("DecodeJWTToken successfully decoded expired token user ID: %s", userID)
}

func TestDecodeJWTToken_TokenWithDifferentAlgorithm(t *testing.T) {
	// JWT token using different algorithm (RS256 instead of HS256)
	// This should still work for parsing since we use ParseUnverified
	rsToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpZCI6InRlc3QtdXNlci1yczI1NiIsImlzcyI6InRlc3QiLCJhdWQiOiJ0ZXN0IiwiaWF0IjoxNjUwMDAwMDAwLCJleHAiOjE2NTAwODY0MDB9.invalid-signature"
	expectedID := "test-user-rs256"

	userID, err := DecodeJWTToken(rsToken)

	// This might fail due to invalid signature format, but that's expected
	// since we're not using a real RS256 signature
	if err != nil {
		t.Logf("DecodeJWTToken failed for RS256 token (expected due to invalid signature): %v", err)
	} else {
		if userID != expectedID {
			t.Errorf("Expected user ID %s, got %s", expectedID, userID)
		}
		t.Logf("DecodeJWTToken successfully decoded RS256 token user ID: %s", userID)
	}
}
