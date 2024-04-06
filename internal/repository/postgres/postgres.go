package postgres

import (
	"banner_service/internal/entity"
	"context"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	maxPaginationLimit = 10

	sortAscending  string = "ASC"
	sortDescending string = "DESC"
)

type Repo struct {
	//pool *pgxpool.Pool
	db *sqlx.DB
	//db *db
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
		//pool: db,
	}
}

func (r *Repo) GetUserBanner(ctx context.Context, userBanner *entity.UserBannerRequest) (entity.UserBannerResponse, error) {
	//data := entity.UserBannerResponse{}

	//_, err := r.pool.QueryRow(ctx, `SELECT `)
	return entity.UserBannerResponse{}, nil
}

func (r *Repo) CreateBanner(ctx context.Context, banner *entity.Banner) (uint64, error) {
	now := time.Now()
	var newID uint64

	beginx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// insert into feature
	errInsertFeature := beginx.QueryRowContext(ctx, `INSERT INTO "feature" (id, created_at, updated_at)
	VALUES ($1, $2, $3) RETURNING id`, banner.FeatureID, now, now).Scan(&newID)
	if errInsertFeature != nil {
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return 0, errRollBack
		}
		return 0, errInsertFeature
	}

	// insert into tag
	// TODO: TadID is slice!!!!!!!!!!!!!!!!!!!!!!!!!!!
	errInsertTag := beginx.QueryRowContext(ctx, `INSERT INTO "tag" (id, created_at, updated_at)
	VALUES ($1, $2, $3) RETURNING id`, banner.TagID[0], now, now).Scan(&newID)
	if errInsertTag != nil {
		err := beginx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	errInsertBanner := beginx.QueryRowContext(ctx, `INSERT INTO "banner" (feature_id, title, text, url, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive, now, now).Scan(&newID)
	if errInsertBanner != nil {
		err := beginx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	// insert into banner_tag
	// TODO: TadID is slice!!!!!!!!!!!!!!!!!!!!!!!!!!!
	_, errInsertBannerTag := beginx.ExecContext(ctx, `INSERT INTO "banner_tag" (banner_id, tag_id)
    VALUES ($1, $2)`, newID, banner.TagID[0])
	if errInsertBannerTag != nil {
		err := beginx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	errCommit := beginx.Commit()
	if errCommit != nil {
		return 0, errCommit
	}
	return newID, nil
}
