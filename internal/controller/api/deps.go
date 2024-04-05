package api

import (
	"banner_service/internal/entity"
	"context"
)

type useCase interface {
	GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]entity.Product, error)
}

type logger interface {
	Info(text ...any)
	Error(text ...any)
	Errorf(format string, args ...any)
}
