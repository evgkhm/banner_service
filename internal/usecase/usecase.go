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

type CacheBanner struct {
	Banner entity.Content
	Timer  time.Time
	mu     sync.Mutex
}

var CacheUserBanner = CacheBanner{Banner: entity.Content{}, Timer: time.Now()}

func New(r repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

type FakeTimeProvider struct {
	currentTime time.Time
}

func (f FakeTimeProvider) Now() time.Time {
	return f.currentTime
}

func (s *UseCase) GetUserBanner(ctx context.Context, tagID, featureID uint64, useLastVersion bool) (entity.Content, error) {
	timeProvider := RealTimeProvider{}

	if useLastVersion || timeProvider.Now().Sub(CacheUserBanner.Timer) >= 5*time.Minute {
		newBanner, err := s.repo.GetUserBanner(ctx, tagID, featureID)
		if err != nil {
			return newBanner, err
		}

		CacheUserBanner.mu.Lock()
		CacheUserBanner.Banner = newBanner
		CacheUserBanner.Timer = timeProvider.Now()
		CacheUserBanner.mu.Unlock()
		return CacheUserBanner.Banner, nil
	}
	return CacheUserBanner.Banner, nil
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
