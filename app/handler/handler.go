package handler

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/config"
	"github.com/opchaves/gin-web-app/app/service"
)

// Handler struct holds required services for handler to function
type Handler struct {
	Db           *pgxpool.Pool
	Logger       *slog.Logger
	Cfg          *config.Config
	MaxBodyBytes int64
	UserService  service.UserService
}
