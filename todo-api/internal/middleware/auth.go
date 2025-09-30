package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok || claims.Subject == "" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			userId, err := uuid.Parse(claims.Subject)

			if err != nil {
				http.Error(w, "Invalid user id", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), auth.UserIdKey, userId)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
