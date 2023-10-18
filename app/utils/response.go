package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/model"
)

func ToFieldErrorResponse(c *gin.Context, field, message string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"errors": []model.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	})
}
