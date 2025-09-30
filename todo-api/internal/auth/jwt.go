package auth

import (
	"time"
	"todo/internal/config"
	"todo/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User, cfg config.Config) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Name:  user.Name,
		Email: user.Email,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(cfg.JWTSecret))
}
