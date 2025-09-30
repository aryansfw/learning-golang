package services

import (
	"todo/internal/models"
	"todo/internal/repository"

	"github.com/google/uuid"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *AuthService) Register(name string, email string, password string) (*models.User, error) {
	user := models.User{
		Id:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	return s.repo.CreateUser(user)
}
