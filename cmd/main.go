package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/opchaves/gin-web-app/app/config"
	"github.com/opchaves/gin-web-app/cmd/server"
)

// Version set current code version
var Version = "0.0.1"

func main() {
	if err := run(); err != nil {
		slog.Error("error: ", slog.AnyValue(err))
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cfg, logger, err := initialize(ctx)
	if err != nil {
		logger.Error("failed to initialize", slog.AnyValue(err))
		return err
	}

	logger.Debug("Connecting to database...")

	db, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		return err
	}

	return server.Start(&server.Config{
		Logger:          logger,
		Db:              db,
		Cfg:             cfg,
		Ctx:             ctx,
		TimeoutDuration: time.Duration(cfg.HandlerTimeOut) * time.Second,
		MaxBodyBytes:    cfg.MaxBodyBytes,
	})
}

func initialize(ctx context.Context) (*config.Config, *slog.Logger, error) {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	buildInfo, _ := debug.ReadBuildInfo()
	logger := slog.New(handler).WithGroup("program_info")
	loggerChild := logger.With(
		slog.Int("pid", os.Getpid()),
		slog.String("version", Version),
		slog.String("go_version", buildInfo.GoVersion),
	)

	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			loggerChild.Error("Error loading .env file", slog.AnyValue(err))
			return nil, nil, err
		}
	}

	cfg, err := config.LoadConfig(ctx)

	return &cfg, loggerChild, err
}
