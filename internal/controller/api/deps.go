package api

import (
	"banner_service/internal/entity"
	"context"
)

type useCase interface {
	GetUserBanner(ctx context.Context, tagID uint64, featureID uint64, useLastVersion bool) (entity.UserBannerResponse, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error)
	GetBanners(ctx context.Context, tagID []uint64, featureID uint64, limit uint64, offset uint64) ([]entity.BannersList, error)
}

type logger interface {
	Info(text ...any)
	Error(text ...any)
	Errorf(format string, args ...any)
}
