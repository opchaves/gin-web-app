package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Gin Web App",
	})
}

func (h *Handler) AddCar(c *gin.Context) {
	carName := c.PostForm("car")

	c.HTML(http.StatusOK, "add_car.html", gin.H{"name": carName})
}
