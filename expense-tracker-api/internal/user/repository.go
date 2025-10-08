package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	ctx := context.Background()

	user, err := gorm.G[User](r.db).Where("email = ?", email).First(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(user *User) error {
	ctx := context.Background()
	if err := gorm.G[User](r.db).Create(ctx, user); err != nil {
		return err
	}
	return nil
}
