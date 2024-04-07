package usecase

import (
	"banner_service/internal/entity"
	"context"
)

type UseCase struct {
	repo repository
}

func New(r repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (s *UseCase) GetUserBanner(ctx context.Context, tagID uint64, featureID uint64, useLastVersion bool) (entity.UserBannerResponse, error) {
	return s.repo.GetUserBanner(ctx, tagID, featureID, useLastVersion)
}

func (s *UseCase) CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error) {
	return s.repo.CreateBanner(ctx, banner)
}

func (s *UseCase) GetBanners(ctx context.Context, tagID uint64, featureID uint64, limit uint64, offset uint64) ([]entity.BannersList, error) {
	return s.repo.GetBanners(ctx, tagID, featureID, limit, offset)
}
