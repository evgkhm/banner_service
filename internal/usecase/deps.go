package usecase

import "context"

type repository interface {
	GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]models.Product, error)
}
