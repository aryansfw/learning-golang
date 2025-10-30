package api

import (
	"context"
	"ecommerce/internal/auth"
	"ecommerce/internal/domain"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Name     string
	Email    string
	Password string
}

type authResponse struct {
	Token string `json:"token"`
}

func (a *api) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var rr registerRequest

	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(rr.Password), 12)
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user := domain.User{
		Name:     rr.Name,
		Email:    rr.Email,
		Password: string(encryptedPassword),
	}

	if err := a.userRepo.Create(ctx, &user); err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	token, err := auth.JWTGenerateToken(user)
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	a.successResponse(w, authResponse{Token: token})
}

type loginRequest struct {
	Email    string
	Password string
}

func (a *api) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var lr loginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	user, err := a.userRepo.GetByEmail(ctx, lr.Email)
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(lr.Password))
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	token, err := auth.JWTGenerateToken(user)
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	a.successResponse(w, authResponse{Token: token})
}

func (a *api) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			a.errorResponse(w, http.StatusUnauthorized, "no token provided")
			return
		}

		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &domain.UserClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil || !token.Valid {
			a.errorResponse(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(*domain.UserClaims)
		if !ok {
			a.errorResponse(w, http.StatusUnauthorized, "failed to parse")
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
