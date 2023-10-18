package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
	"github.com/opchaves/gin-web-app/app/service"
	"github.com/opchaves/gin-web-app/app/utils"
)

func (h *Handler) Register(c *gin.Context) {
	var input service.RegisterInput

	if err := c.ShouldBind(&input); err != nil {
		h.Logger.Info("failed to validate", slog.String("error", err.Error()))
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]model.FieldError, len(ve))
			for i, fe := range ve {
				fmt.Println(fe.Type(), fe.Param(), fe.Value())
				out[i] = model.FieldError{
					Field:   fe.Field(),
					Message: msgForTag(fe.Tag(), fe.Param(), fe.Type()),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": out})
		}
		return
	}

	user, err := h.UserService.Register(c, &input)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			utils.ToFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}

		c.JSON(apperrors.Status(err), gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": h.UserService.GetRegisterResponse(user),
		"code": http.StatusCreated,
	})
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
