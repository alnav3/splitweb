package encryption

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(password string) (string, error) {
    bcryptHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(bcryptHash), nil
}

func Compare(password, encryptedPass string) error {
    err := bcrypt.CompareHashAndPassword([]byte(encryptedPass), []byte(password))
    if err != nil {
        return err
    }
    return nil
}

func GenerateJWT(username, secret string) string {
    // Create the Claims

    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["user"] = username
    claims["authorized"] = true
    claims["exp"] = time.Now().Add(time.Hour * 10 * 24).Unix()
    tokenString, _ := token.SignedString([]byte(secret))

    return tokenString
}

func ValidateToken(tokenString, secret string) (string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte(secret), nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        if claims["authorized"] == true {
            return claims["user"].(string), nil
        }
    }

    return "", errors.New("invalid token")
}

