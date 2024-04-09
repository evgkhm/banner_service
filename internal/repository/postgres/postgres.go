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

var ErrUserBanner = errors.New("баннер не найден")

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetUserBanner(ctx context.Context, tagID, featureID uint64) (entity.Content, error) {
	data := entity.Content{}

	beginx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return entity.Content{}, err
	}

	query := `SELECT b.title, b.text, b.url FROM banner b
			JOIN banner_tag bg ON bg.banner_id = b.id
			WHERE bg.tag_id = $1 AND b.feature_id = $2`

	errRow := beginx.QueryRowContext(ctx, query, tagID, featureID).Scan(&data.Title, &data.Text, &data.URL)
	if errRow != nil {
		if errors.Is(errRow, sql.ErrNoRows) {
			errRollBack := beginx.Rollback()
			if errRollBack != nil {
				return entity.Content{}, errRollBack
			}
			return entity.Content{}, ErrUserBanner
		}
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return entity.Content{}, errRollBack
		}
		return entity.Content{}, errRow
	}

	errCommit := beginx.Commit()
	if errCommit != nil {
		return entity.Content{}, errCommit
	}
	return data, nil
}

func (r *Repo) GetBanners(ctx context.Context, tagID, featureID, limit, offset uint64) ([]entity.BannersList, error) {
	var banners []entity.BannersList
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
	_, errInsertFeature := beginx.ExecContext(ctx, `INSERT INTO "feature" (id, created_at, updated_at)
	VALUES ($1, $2, $3)  ON CONFLICT DO NOTHING`, banner.FeatureID, now, now)
	if errInsertFeature != nil {
		errRollBack := beginx.Rollback()
		if errRollBack != nil {
			return 0, errRollBack
		}
		return 0, errInsertFeature
	}

	// insert into tag slice of tagIDs
	valueStrings := make([]string, 0, len(banner.TagID))
	valueArgs := make([]interface{}, 0, len(banner.TagID))
	for _, tagID := range banner.TagID {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, tagID, now, now)
	}
	justString := strings.Join(valueStrings, ", ")
	query := "INSERT INTO tag (id, created_at, updated_at) VALUES " +
		sqlx.Rebind(sqlx.DOLLAR, justString) +
		" ON CONFLICT DO NOTHING"

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
	valueStringsBannerTags := make([]string, 0, len(banner.TagID))
	valueArgsBannerTags := make([]interface{}, 0, len(banner.TagID))
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

func (r *Repo) UpdateBanner(ctx context.Context, bannerID uint64, banner *entity.Banner) error {
	now := time.Now()

	tr, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	updateFeature := `UPDATE "feature" SET id = $1, created_at = $2, updated_at = $3 WHERE id = 
		(SELECT feature_id FROM banner WHERE id = $4)`
	_, errExec := tr.ExecContext(ctx, updateFeature, banner.FeatureID, now, now, bannerID)
	if errExec != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errExec
	}

	valueStringsTags := make([]string, 0, len(banner.TagID))
	valueArgsTags := make([]interface{}, 0, len(banner.TagID))
	for _, tagID := range banner.TagID {
		valueStringsTags = append(valueStringsTags, "(?, ?, ?)")
		valueArgsTags = append(valueArgsTags, tagID, now, now)
	}
	justString := strings.Join(valueStringsTags, ", ")
	query := "INSERT INTO tag (id, created_at, updated_at) VALUES " +
		sqlx.Rebind(sqlx.DOLLAR, justString) + " ON CONFLICT DO NOTHING"
	_, errInsertBannerTag := tr.ExecContext(ctx, query, valueArgsTags...)
	if errInsertBannerTag != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errInsertBannerTag
	}

	deleteBannerTags := `DELETE from banner_tag WHERE banner_id = $1`
	_, errDeleteBannerTag := tr.ExecContext(ctx, deleteBannerTags, bannerID)
	if errDeleteBannerTag != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errDeleteBannerTag
	}

	// insert into banner_tag
	valueStringsBannerTags := make([]string, 0, len(banner.TagID))
	valueArgsBannerTags := make([]interface{}, 0, len(banner.TagID))
	for _, tagID := range banner.TagID {
		valueStringsBannerTags = append(valueStringsBannerTags, "(?, ?)")
		valueArgsBannerTags = append(valueArgsBannerTags, bannerID, tagID)
	}
	justString = strings.Join(valueStringsBannerTags, ", ")
	query = "INSERT INTO banner_tag (banner_id, tag_id) VALUES " +
		sqlx.Rebind(sqlx.DOLLAR, justString) + " ON CONFLICT DO NOTHING"
	_, errInsertBannerTag = tr.ExecContext(ctx, query, valueArgsBannerTags...)
	if errInsertBannerTag != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errInsertBannerTag
	}

	errCommit := tr.Commit()
	if errCommit != nil {
		return errCommit
	}

	return nil
}

func (r *Repo) DeleteBanner(ctx context.Context, bannerID uint64) error {
	tr, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	deleteBannerTags := `DELETE from banner_tag WHERE banner_id = $1`
	_, errDelete := tr.ExecContext(ctx, deleteBannerTags, bannerID)
	if errDelete != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errDelete
	}

	deleteBanner := `DELETE from banner WHERE id = $1`
	_, errDelete = tr.ExecContext(ctx, deleteBanner, bannerID)
	if errDelete != nil {
		errRollBack := tr.Rollback()
		if errRollBack != nil {
			return errRollBack
		}
		return errDelete
	}

	errCommit := tr.Commit()
	if errCommit != nil {
		return errCommit
	}
	return nil
}
