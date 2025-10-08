package middleware

import (
	"context"
	"net/http"
	"spendime/internal/auth"
	"spendime/internal/contextkey"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authorizationHeader, "Bearer ") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (any, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok || claims.UserID == uuid.Nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), contextkey.UserID, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
