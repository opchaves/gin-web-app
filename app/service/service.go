package service

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/redis/go-redis/v9"
)

type ServiceConfig struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
	Redis  *redis.Client
}
