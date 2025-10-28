package auth

import (
	"ecommerce/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func JWTGenerateToken(user domain.User) (string, error) {
	claims := &domain.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_KEY")))
}