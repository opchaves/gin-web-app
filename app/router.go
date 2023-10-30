package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/handler"
	"github.com/opchaves/gin-web-app/app/handler/middleware"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/opchaves/gin-web-app/app/service"
)

func SetRoutes(c *Config) {
	queries := model.New(c.Db)
	userService := service.NewUserService(&service.ServiceConfig{
		Db:     c.Db,
		Q:      queries,
		Logger: c.Logger,
	})

	h := &handler.Handler{
		Db:           c.Db,
		Logger:       c.Logger,
		Cfg:          c.Cfg,
		MaxBodyBytes: c.MaxBodyBytes,
		UserService:  userService,
	}

	c.Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found. Go to docs for a list of all routes",
		})
	})

	c.Router.GET("/", h.GetHome)
	c.Router.POST("/add-car", h.PostAddCar)

	authGroup := c.Router.Group("/auth")
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.POST("/logout", h.Logout)

	authGroup.Use(middleware.AuthUser(c.Logger))
	authGroup.GET("/me", h.GetCurrent)
}
