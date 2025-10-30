package domain

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

const UserContextKey = "user_id"

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)

	Create(ctx context.Context, user *User) error
}
