package usecase

import (
	"banner_service/internal/entity"
	"context"
)

type repository interface {
	GetUserBanner(ctx context.Context, userBanner *entity.UserBannerRequest) (entity.UserBannerResponse, error)
}
