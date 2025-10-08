package transaction

import "github.com/google/uuid"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(name string, amount int64, transactionType Type) (*Transaction, error) {
	transaction := Transaction{
		Name:   name,
		Amount: amount,
		Type:   transactionType,
	}
	if err := s.repo.Create(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (s *Service) List(userID uuid.UUID) (*[]Transaction, error) {
	transactions, err := s.repo.List(userID)

	if err != nil {
		return nil, err
	}

	return transactions, nil
}
