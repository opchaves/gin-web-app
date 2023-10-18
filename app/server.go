package app

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/opchaves/gin-web-app/app/config"
)

// Version set current code version
var Version = "0.0.1"

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Db              *pgxpool.Pool
	Cfg             *config.Config
	Ctx             context.Context
	Logger          *slog.Logger
	Router          *gin.Engine
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func Setup() (*Config, error) {
	ctx := context.Background()
	cfg, logger, err := initialize(ctx)

	if err != nil {
		logger.Error("failed to initialize", slog.AnyValue(err))
		return nil, err
	}

	logger.Debug("Connecting to database...")

	db, err := pgxpool.New(ctx, cfg.DatabaseUrl)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("error", err.Error()))
		return nil, err
	}

	corsConfig := cors.DefaultConfig()
	router := gin.Default()
	router.Use(cors.New(corsConfig))
	router.LoadHTMLGlob(cfg.TemplatesGlob)
	router.Static("/assets", cfg.AssetsDir)

	config := &Config{
		Db:              db,
		Cfg:             cfg,
		Logger:          logger,
		Router:          router,
		Ctx:             ctx,
		TimeoutDuration: time.Duration(cfg.HandlerTimeOut) * time.Second,
		MaxBodyBytes:    cfg.MaxBodyBytes,
	}

	SetRoutes(config)

	return config, nil
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
