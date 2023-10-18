package handler

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
	"github.com/opchaves/gin-web-app/app/service"
	"github.com/opchaves/gin-web-app/app/utils"
)

func (h *Handler) Register(c *gin.Context) {
	var req service.RegisterInput

	if err := c.ShouldBind(&req); err != nil {
		errors := parseError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	user, err := h.UserService.Register(c, &req)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			utils.ToFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}

		c.JSON(apperrors.Status(err), gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
	return
}

func msgForTag(tag string, param string, t reflect.Type) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		if t.Name() == "string" {
			return fmt.Sprintf("cannot have less than %s characters", param)
		}
		return fmt.Sprintf("min value is %s, %s", param, t)
	case "max":
		if t.Name() == "string" {
			return fmt.Sprintf("cannot be longer than %s characters", param)
		}
		return fmt.Sprintf("max value is %s", param)
	}

	return ""
}
