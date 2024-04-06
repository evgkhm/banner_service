package postgres

import (
	"banner_service/internal/entity"
	"context"
	"github.com/jmoiron/sqlx"
	"strings"
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

	// insert into tag slice of tagIDs
	var valueStrings []string
	var valueArgs []interface{}
	for _, tagID := range banner.TagID {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, tagID, now, now)
	}
	justString := strings.Join(valueStrings, ", ")
	query := "INSERT INTO tag (id, created_at, updated_at) VALUES " +
		sqlx.Rebind(sqlx.DOLLAR, justString)
	_, errInsertTag := beginx.ExecContext(ctx, query, valueArgs...)
	if errInsertTag != nil {
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return 0, errRollBack
		}
		return 0, errInsertTag
	}

	// insert into banner
	errInsertBanner := beginx.QueryRowContext(ctx, `INSERT INTO "banner" (feature_id, title, text, url, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, banner.FeatureID, banner.Content.Title, banner.Content.Text, banner.Content.URL, banner.IsActive, now, now).Scan(&newID)
	if errInsertBanner != nil {
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return 0, errRollBack
		}
		return 0, errInsertBanner
	}

	// insert into banner_tag
	var valueStringsBannerTags []string
	var valueArgsBannerTags []interface{}
	for _, tagID := range banner.TagID {
		valueStringsBannerTags = append(valueStringsBannerTags, "(?, ?)")
		valueArgsBannerTags = append(valueArgsBannerTags, newID, tagID)
	}
	justString = strings.Join(valueStringsBannerTags, ", ")
	query = "INSERT INTO banner_tag (banner_id, tag_id) VALUES " +
		sqlx.Rebind(sqlx.DOLLAR, justString)
	_, errInsertBannerTag := beginx.ExecContext(ctx, query, valueArgsBannerTags...)
	if errInsertBannerTag != nil {
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return 0, errRollBack
		}
		return 0, errInsertBannerTag
	}

	errCommit := beginx.Commit()
	if errCommit != nil {
		return 0, errCommit
	}
	return newID, nil
}
