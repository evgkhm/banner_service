package postgres

import (
	"banner_service/internal/entity"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	maxPaginationLimit = 10

	sortAscending  string = "ASC"
	sortDescending string = "DESC"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	return &Repo{
		pool: db,
	}
}

func (r *Repo) GetUserBanner(ctx context.Context, userBanner *entity.UserBannerRequest) (entity.UserBannerResponse, error) {
	//data := entity.UserBannerResponse{}

	//_, err := r.pool.QueryRow(ctx, `SELECT `)
	return entity.UserBannerResponse{}, nil
}

func (r *Repo) CreateBanner(ctx context.Context, banner *entity.Banner) error {
	now := time.Now()
	var newID uint64
	// pool insert with returning id //QueryRow?
	err := r.pool.Exec(ctx, `INSERT INTO "banner" (feature_id, title, text, url, is_active, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		banner.FeatureID, banner.Title, banner.Text, banner.URL, banner.IsActive, now, now).Scan(&newID)
	if err != nil {
		return err
	}

	_, errExecBannerTags := r.pool.Exec(ctx, `INSERT INTO "banner_tag" (banner_id, tag_id) VALUES ($1, $2)`, newID, banner.TagID)
	if errExecBannerTags != nil {
		return errExecBannerTags
	}

	_, errExecFeatures := r.pool.Exec(ctx, `INSERT INTO "feature" (feature_id, created_at, updated_at) 
VALUES ($1, $2, $3)`, banner.FeatureID, now, now)
	if errExecFeatures != nil {
		return errExecFeatures
	}

	// slice of tags!!!!!
	_, errExecTags := r.pool.Exec(ctx, `INSERT INTO "tag" (tag_id, created_at, updated_at)
VALUES ($1, $2, $3)`, banner.TagID, now, now)
	if errExecTags != nil {
		return errExecTags
	}

	return nil
}
