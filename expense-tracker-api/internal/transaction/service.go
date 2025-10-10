package transaction

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(ctx context.Context, transaction *Transaction) error {
	if err := s.repo.Create(ctx, transaction); err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context, filters TransactionFilters) (*[]Transaction, error) {
	transactions, err := s.repo.List(ctx, filters)

	if err != nil {
		return nil, err
	}

	for idx := range *transactions {
		(*transactions)[idx].CategoryID = nil
	}

	return transactions, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, transaction *Transaction) error {
	if err := s.repo.Update(ctx, transaction); err != nil {
		return err
	}

	return nil
}
