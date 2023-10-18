package handler

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

// setUserSession saves the users ID in the session
func (h *Handler) setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		h.Logger.Error("error setting the session", slog.AnyValue(err.Error()))
	}
}
