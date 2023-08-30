package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct {
	MaxBodyBytes int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func NewHandler(c *Config) {
	h := &Handler{
		MaxBodyBytes: c.MaxBodyBytes,
	}

	c.R.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No route found. Go to docs for a list of all routes",
		})
	})

	c.R.Static("/assets", "./assets")

	// Home page
	c.R.GET("/", h.Home)
	c.R.POST("/add-car", h.AddCar)
	// About page
	// c.R.GET("/about", h.About)

	// Create an account group
	ag := c.R.Group("api/account")
	ag.GET("/me", h.Me)
	ag.POST("/register", h.Register)
}
