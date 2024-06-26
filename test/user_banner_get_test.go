package test

import (
	"banner_service/internal/config"
	"banner_service/internal/controller/api"
	"banner_service/internal/entity"
	"banner_service/internal/repository/postgres"
	"banner_service/internal/usecase"
	"banner_service/pkg/httpserver"
	"banner_service/pkg/logging"
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"regexp"
	"testing"
)

func TestGetUserBanner(t *testing.T) {
	t.Parallel()

	// Load .env file
	re := regexp.MustCompile(`^(.*` + "banner_service" + `)`)
	cwd, errGetWd := os.Getwd()
	if errGetWd != nil {
		t.Errorf("failed to get current working directory: %v", errGetWd)
	}
	rootPath := re.Find([]byte(cwd))
	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		t.Errorf("failed to load .env file: %v", err)
		os.Exit(-1)
	}

	// Parse config
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		t.Errorf("failed to parse config from environment variables: %v", err)
	}
	dbURL := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s sslmode=%s",
		"localhost", 5433, cfg.PG.User, cfg.PG.DB, cfg.PG.Password, cfg.PG.SSLMode)

	dbConn, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		t.Errorf("failed to create connection pool: %v", err)
	}
	errPing := dbConn.Ping()
	if errPing != nil {
		t.Errorf("failed to ping the database: %v", errPing)
	}

	logger, err := logging.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	require.NoError(t, err)

	ctx := context.Background()

	repo := postgres.New(dbConn)

	useCase := usecase.New(repo)

	handler := api.New(useCase, logger)

	httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	t.Run("Баннер пользователя для админа", func(t *testing.T) {
		t.Parallel()
		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  true,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=false", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("token", os.Getenv("API_TOKEN"))
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusOK, resp.StatusCode)

		expectedBanner, err := useCase.GetUserBanner(ctx, requestForCreate.TagID[0], requestForCreate.FeatureID, true, false)
		require.NoError(t, err)
		require.Equal(t, expectedBanner, requestContent)
	})

	t.Run("Баннер пользователя для юзера с активным баннером", func(t *testing.T) {
		t.Parallel()
		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  true,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=false", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("token", os.Getenv("API_TOKEN"))
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusOK, resp.StatusCode)

		expectedBanner, err := useCase.GetUserBanner(ctx, requestForCreate.TagID[0], requestForCreate.FeatureID, true, true)
		require.NoError(t, err)
		require.Equal(t, expectedBanner, requestContent)
	})

	t.Run("Баннер пользователя для юзера с неактивным баннером", func(t *testing.T) {
		t.Parallel()
		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  false,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=false", http.NoBody)
		require.NoError(t, err)
		req.Header.Set("token", os.Getenv("API_TOKEN"))
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusOK, resp.StatusCode)

		expectedBanner, err := useCase.GetUserBanner(ctx, requestForCreate.TagID[0], requestForCreate.FeatureID, true, true)
		require.NoError(t, err)
		require.Equal(t, expectedBanner, requestContent)
	})

	t.Run("Некоректные данные", func(t *testing.T) {
		t.Parallel()

		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  true,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/user_banner", http.NoBody)
		require.NoError(t, err)

		req.Header.Set("token", os.Getenv("API_TOKEN"))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Пользователь не авторизован", func(t *testing.T) {
		t.Parallel()

		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  true,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=false", http.NoBody)
		require.NoError(t, err)

		req.Header.Set("token", "user_token_not_right")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Баннер не найден", func(t *testing.T) {
		t.Parallel()

		requestContent := entity.Content{
			Title: "some_title",
			Text:  "some_text",
			URL:   "some_url",
		}
		var requestForCreate = entity.Banner{
			TagID:     []uint64{1},
			FeatureID: 1,
			Content:   requestContent,
			IsActive:  true,
		}
		createBanner, err := repo.CreateBanner(ctx, &requestForCreate)
		require.NoError(t, err)
		require.NotZero(t, createBanner)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet,
			"http://localhost:8080/user_banner?tag_id=10000&feature_id=10000&use_last_revision=true", http.NoBody)
		require.NoError(t, err)

		req.Header.Set("token", os.Getenv("API_TOKEN"))

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { _ = resp.Body.Close() })

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
