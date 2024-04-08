package postgres

import (
	"banner_service/internal/entity"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

const (
	maxPaginationLimit = 10

	sortAscending  string = "ASC"
	sortDescending string = "DESC"
)

var ErrorUserBannerNotFound = errors.New("баннер не найден")

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetUserBanner(ctx context.Context, tagID uint64, featureID uint64, useLastVersion bool) (entity.UserBannerResponse, error) {
	data := entity.UserBannerResponse{}

	beginx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.UserBannerResponse{}, err
	}

	query := `SELECT b.title, b.text, b.url FROM banner b
			JOIN banner_tag bg ON bg.banner_id = b.id
			WHERE bg.tag_id = $1 AND b.feature_id = $2`

	errRow := beginx.QueryRowContext(ctx, query, tagID, featureID).Scan(&data.Title, &data.Text, &data.URL)
	if errRow != nil {
		if errors.Is(errRow, sql.ErrNoRows) {
			errRollBack := beginx.Rollback()
			if errRollBack != nil {
				return entity.UserBannerResponse{}, errRollBack
			}
			return entity.UserBannerResponse{}, ErrorUserBannerNotFound
		}
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return entity.UserBannerResponse{}, errRollBack
		}
		return entity.UserBannerResponse{}, errRow
	}

	errCommit := beginx.Commit()
	if errCommit != nil {
		return entity.UserBannerResponse{}, errCommit
	}
	return data, nil
}

func (r *Repo) GetBanners(ctx context.Context, tagID []uint64, featureID uint64, limit uint64, offset uint64) ([]entity.BannersList, error) {
	banners := []entity.BannersList{}
	tr, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	query := `WITH aggregated_tags AS (
	SELECT banner_id, string_agg(tag_id::text, ',') AS all_tag_ids
	FROM banner_tag
	GROUP BY banner_id
	)
	
	SELECT b.id, at.all_tag_ids, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at FROM banner b
	LEFT JOIN aggregated_tags at ON b.id = at.banner_id
	WHERE b.feature_id = (?)
	 OR b.id IN (
		SELECT bt.banner_id
		FROM banner_tag bt
		WHERE bt.tag_id IN (?)
		LIMIT (?) OFFSET (?)
	)`
	qry, args, errQuery := sqlx.In(query, featureID, tagID, limit, offset)
	if errQuery != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return nil, errRollBack
		}
		return nil, errQuery
	}
	qry = r.db.Rebind(qry)

	rows, errRows := tr.QueryContext(ctx, qry, args...)
	if errRows != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return nil, errRollBack
		}
		return nil, errRows
	}

	for rows.Next() {
		var row entity.BannersList
		errScan := rows.Scan(&row.BannerID, &row.TagIDs, &row.FeatureID, &row.Content.Title, &row.Content.Text,
			&row.Content.URL, &row.IsActive, &row.CreatedAt, &row.UpdatedAt)
		if errScan != nil {
			errRollBack := tr.Rollback()
			if errRollBack != nil {
				return nil, errRollBack
			}
			return nil, errScan
		}
		banners = append(banners, row)
	}
	errCommit := tr.Commit()
	if errCommit != nil {
		return nil, errCommit
	}
	return banners, nil
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
