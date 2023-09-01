package server

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/config"
	"github.com/opchaves/gin-web-app/app/model"
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

func SetRoutes(c *Config, router *gin.Engine) {
	queries := model.New(c.Db)
	userService := service.NewUserService(&service.UserServiceConfig{Db: c.Db, Q: queries, Logger: c.Logger})

	h := &Handler{
		Db:           c.Db,
		Logger:       c.Logger,
		Cfg:          c.Cfg,
		MaxBodyBytes: c.MaxBodyBytes,
		UserService:  userService,
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found. Go to docs for a list of all routes",
		})
	})

	router.GET("/", h.GetHome)
	router.POST("/add-car", h.PostAddCar)

	authGroup := router.Group("/auth")
	authGroup.POST("/register", h.Register)
}

func toFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []model.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	})
}
