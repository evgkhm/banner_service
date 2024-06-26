package usecase

import (
	"banner_service/internal/entity"
	"context"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=mocks_test.go -package=$GOPACKAGE --build_flags=--mod=mod
type repository interface {
	GetUserBanner(ctx context.Context, tagID, featureID uint64) (entity.Content, bool, error)
	CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error)
	GetBanners(ctx context.Context, tagID, featureID, limit, offset uint64) ([]entity.BannersList, error)
	UpdateBanner(ctx context.Context, bannerID uint64, banner *entity.Banner) error
	DeleteBanner(ctx context.Context, bannerID uint64) error
}

type TimeProvider interface {
	Now() time.Time
}
