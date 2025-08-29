package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func DecodeJWTToken(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &JWTClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok {
		return claims.Id, nil
	}

	return "", fmt.Errorf("invalid token claims")
}
