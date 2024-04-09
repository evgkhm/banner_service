package usecase

import (
	"banner_service/internal/entity"
	"context"
	"sync"
	"time"
)

type UseCase struct {
	repo repository
}

var cache struct {
	banner entity.Content
	time.Time
	mu sync.Mutex
}

func New(r repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (s *UseCase) GetUserBanner(ctx context.Context, tagID, featureID uint64, useLastVersion bool) (entity.Content, error) {
	if useLastVersion || time.Since(cache.Time) >= 5*time.Minute {
		newBanner, err := s.repo.GetUserBanner(ctx, tagID, featureID)
		if err != nil {
			return newBanner, err
		}

		cache.mu.Lock()
		cache.banner = newBanner
		cache.Time = time.Now()
		cache.mu.Unlock()
		return cache.banner, nil
	}
	return cache.banner, nil
}

func (s *UseCase) CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error) {
	return s.repo.CreateBanner(ctx, banner)
}

func (s *UseCase) GetBanners(ctx context.Context, tagID, featureID, limit, offset uint64) ([]entity.BannersList, error) {
	return s.repo.GetBanners(ctx, tagID, featureID, limit, offset)
}

func (s *UseCase) UpdateBanner(ctx context.Context, bannerID uint64, banner *entity.Banner) error {
	return s.repo.UpdateBanner(ctx, bannerID, banner)
}

func (s *UseCase) DeleteBanner(ctx context.Context, bannerID uint64) error {
	return s.repo.DeleteBanner(ctx, bannerID)
}
