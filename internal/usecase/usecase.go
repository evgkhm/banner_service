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
	//userBannerResponse, err := s.repo.GetUserBanner(ctx, userBanner)
	//if err != nil {
	//	return entity.UserBannerResponse{}, err
	//}
	//return userBannerResponse, nil
	return s.repo.GetUserBanner(ctx, tagID, featureID, useLastVersion)
}

func (s *UseCase) CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error) {
	return s.repo.CreateBanner(ctx, banner)
}
