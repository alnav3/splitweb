package auth

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

type Record struct {
	Avatar          string `json:"avatar"`
	CollectionId    string `json:"collectionId"`
	CollectionName  string `json:"collectionName"`
	Created         string `json:"created"`
	Email           string `json:"email"`
	EmailVisibility bool   `json:"emailVisibility"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Updated         string `json:"updated"`
	Verified        bool   `json:"verified"`
}

type AuthResponse struct {
	Record Record `json:"record"`
	Token  string `json:"token"`
}
