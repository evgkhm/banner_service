package postgres

import (
	"banner_service/internal/config"
	"context"
	"fmt"
	//"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jmoiron/sqlx"
	//"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type db struct {
	//Pool *pgxpool.Pool
	db *sqlx.DB
}

func NewDB(ctx context.Context, cfg config.PGConfig) (*sqlx.DB, error) {
	//pool, err := pgxpool.New(ctx, cfg.URL)
	dbURL := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.DB, cfg.Password, cfg.SSLMode)

	dbConn, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}
	errPing := dbConn.Ping()
	if errPing != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", errPing)
	}
	return dbConn, nil
	//return &DB{Pool: pool}, nil
}

func (db *db) Close() {
	//db.Pool.Close()
	err := db.db.Close()
	if err != nil {
		return
	}
}
