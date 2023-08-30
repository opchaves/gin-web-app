package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home_index.html", gin.H{
		"title": "Gin Web App",
	})
}

func (h *Handler) AddCar(c *gin.Context) {
	c.HTML(http.StatusOK, "add_car.html", gin.H{"name": "Prisma"})
}
