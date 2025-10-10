package category

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) List(ctx context.Context) (*[]Category, error) {
	categories, err := s.repo.List(ctx)

	if err != nil {
		return nil, err
	}

	return categories, nil
}
