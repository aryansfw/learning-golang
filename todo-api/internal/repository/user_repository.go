package repository

import (
	"database/sql"
	"todo/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.QueryRow(
		"SELECT * FROM users WHERE email = $1",
		email).Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(user models.User) (*models.User, error) {
	var resultUser models.User
	if err := r.db.QueryRow(
		"INSERT INTO users VALUES ($1, $2, $3, $4) RETURNING id, name, email, password",
		user.Id, user.Name, user.Email, user.Password,
	).Scan(&resultUser.Id, &resultUser.Name, &resultUser.Email, &resultUser.Password); err != nil {
		return nil, err
	}

	return &resultUser, nil
}
