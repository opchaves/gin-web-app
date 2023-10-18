package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO maybe move routes to app/handlers and only leave router.go and server.go under cmd/server

func (h *Handler) GetHome(c *gin.Context) {
	fmt.Println(h.Db.Config().ConnString())

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Gin Web App",
	})
}

func (h *Handler) PostAddCar(c *gin.Context) {
	carName := c.PostForm("car")

	c.HTML(http.StatusOK, "add_car.html", gin.H{"name": carName})
}
