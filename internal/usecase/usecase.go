package usecase

import (
	"banner_service/internal/entity"
	"context"
)

type service struct {
	repo repository
}

func New(r repository) *service {
	return &service{
		repo: r,
	}
}

func (s *service) GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]entity.Product, error) {
	products, err := s.repo.GetProducts(ctx, page, limit, sortOrder)
	if err != nil {
		return nil, err
	}
	if products == nil {
		return []entity.Product{}, nil
	}
	return products, nil
}
