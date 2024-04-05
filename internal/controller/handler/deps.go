package handler

import "context"

type usecase interface {
	GetProducts(ctx context.Context, page int, limit int, sortOrder string) ([]models.Product, error)
}

type logger interface {
	Info(text ...any)
	Error(text ...any)
	Errorf(format string, args ...any)
}
