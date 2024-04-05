package main

import (
	"banner_service/internal/config"
	"banner_service/internal/controller/api"
	"banner_service/internal/repository/postgres"
	"banner_service/internal/usecase"
	"banner_service/pkg/httpserver"
	"banner_service/pkg/logging"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.New(cfg.Logger.LogFilePath, cfg.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = logger.Sync() }()

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.PG)
	if err != nil {
		logger.Fatalf("failed to connect to postgres db: %s", err)
	}
	defer db.Close()

	repo := postgres.New(db.Pool)

	useCase := usecase.New(repo)

	logger.Info("Starting api server...")
	handler := api.New(useCase, logger)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app: main: signal: " + s.String())
	case err = <-httpServer.Notify():
		logger.Errorf("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		logger.Errorf("app: main: httpServer.Shutdown: %v", err)
	}
}
