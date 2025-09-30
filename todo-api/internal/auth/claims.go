package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}
