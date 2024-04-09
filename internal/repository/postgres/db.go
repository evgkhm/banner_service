package postgres

import (
	"banner_service/internal/config"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
)

func NewDB(_ context.Context, cfg config.PGConfig) (*sqlx.DB, error) {
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
}

func CloseDB(db *sqlx.DB) error {
	return db.Close()
}
