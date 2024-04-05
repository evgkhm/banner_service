package usecase

import (
	"banner_service/internal/entity"
	"context"
)

type repository interface {
	GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]entity.Product, error)
}
