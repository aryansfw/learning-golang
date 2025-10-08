package auth

import (
	"errors"
	"spendime/internal/password"
	"spendime/internal/user"

	"github.com/google/uuid"
)

type Service struct {
	repo *user.Repository
}

func NewService(repo *user.Repository) *Service {
	return &Service{repo}
}

func (s *Service) Login(email, pw, jwtSecret string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if password.Compare(pw, user.Password) {
		return "", errors.New("invalid password")
	}

	token, err := GenerateJWT(user.ID, jwtSecret)

	if err != nil {
		return "", errors.New("error generating token")
	}

	return token, nil
}

func (s *Service) Register(name, email, pw, jwtSecret string) (string, uuid.UUID, error) {
	hash, err := password.Encrypt(pw)
	if err != nil {
		return "", uuid.Nil, err
	}

	user := &user.User{
		Name:     name,
		Email:    email,
		Password: hash,
	}

	if err := s.repo.Create(user); err != nil {
		return "", uuid.Nil, err
	}

	token, err := GenerateJWT(user.ID, jwtSecret)
	if err != nil {
		return "", uuid.Nil, err
	}

	return token, user.ID, nil
}
