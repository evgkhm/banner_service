package api

import (
	"banner_service/internal/entity"
	"context"
)

type useCase interface {
	GetUserBanner(ctx context.Context, tagID, featureID uint64, useLastVersion bool, isUserRequest bool) (entity.Content, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset uint64) ([]entity.BannersList, error)
	UpdateBanner(ctx context.Context, bannerID uint64, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, bannerID uint64) error
}

type logger interface {
	Info(text ...any)
	Error(text ...any)
	Errorf(format string, args ...any)
}
