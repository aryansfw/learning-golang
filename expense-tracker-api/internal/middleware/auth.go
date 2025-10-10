package middleware

import (
	"context"
	"net/http"
	"spendime/internal/auth"
	"spendime/internal/contextkey"
	"spendime/internal/response"
	"spendime/internal/user"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Auth(jwtSecret string, gormDB *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authorizationHeader, "Bearer ") {
				response.Error(w, "Invalid token", http.StatusUnauthorized, "invalid token")
				return
			}

			tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				response.Error(w, "Invalid token", http.StatusUnauthorized, "invalid token")
				return
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok || claims.UserID == uuid.Nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			_, err = gorm.G[user.User](gormDB).Where("id = ?", claims.UserID).First(ctx)

			if err != nil {
				response.Error(w, "Invalid token", http.StatusUnauthorized, "user does not exist")
				return
			}

			ctx = context.WithValue(r.Context(), contextkey.UserID, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
