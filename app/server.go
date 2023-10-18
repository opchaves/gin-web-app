package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-contrib/sessions"
	sessionRedis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/opchaves/gin-web-app/app/config"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// Version set current code version
var Version = "0.0.1"

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Db              *pgxpool.Pool
	Cfg             *config.Config
	Ctx             context.Context
	RedisClient     *redis.Client
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

	// Initialize redis connection
	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		logger.Error("error parsing the redis url", slog.String("error", err.Error()))
		return nil, err
	}

	logger.Debug("Connecting to redis...")
	rdb := redis.NewClient(opt)

	// verify redis connection
	_, err = rdb.Ping(ctx).Result()

	if err != nil {
		logger.Error("error connecting to redis", slog.String("error", err.Error()))
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
		RedisClient:     rdb,
		Ctx:             ctx,
		TimeoutDuration: time.Duration(cfg.HandlerTimeOut) * time.Second,
		MaxBodyBytes:    cfg.MaxBodyBytes,
	}

	redisURL := rdb.Options().Addr
	password := rdb.Options().Password

	// initialize session store
	store, err := sessionRedis.NewStore(10, "tcp", redisURL, password, []byte(cfg.SessionSecret))

	if err != nil {
		logger.Error("could not initialize redis session store", slog.String("error", err.Error()))
		return nil, err
	}

	store.Options(sessions.Options{
		Domain:   cfg.Domain,
		MaxAge:   60 * 60 * 24 * 7, // 7 days
		Secure:   gin.Mode() == gin.ReleaseMode,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions(model.CookieName, store))

	// add rate limit
	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1500,
	}

	limitStore, _ := sredis.NewStore(rdb)

	rateLimiter := mgin.NewMiddleware(limiter.New(limitStore, rate))
	router.Use(rateLimiter)

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
