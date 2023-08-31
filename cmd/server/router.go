package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/config"
)

// Handler struct holds required services for handler to function
type Handler struct {
	Db           *pgxpool.Pool
	Logger       *slog.Logger
	Cfg          *config.Config
	MaxBodyBytes int64
}

func SetRoutes(c *Config, router *gin.Engine) {
	h := &Handler{
		Db:           c.Db,
		Logger:       c.Logger,
		Cfg:          c.Cfg,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found. Go to docs for a list of all routes",
		})
	})

	router.GET("/", h.GetHome)
	router.POST("/add-car", h.PostAddCar)
}
