package service

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/model"
)

type ServiceConfig struct {
	Q      *model.Queries
	Logger *slog.Logger
	Db     *pgxpool.Pool
}
