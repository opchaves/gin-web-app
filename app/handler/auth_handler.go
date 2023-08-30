package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

/*
 * AuthHandler contains all routes related to account actions (/api/account)
 */

type registerReq struct {
	// Must be unique
	Email string `json:"email"`
	// Min 3, max 30 characters.
	Username string `json:"username"`
	// Min 6, max 150 characters.
	Password string `json:"password"`
} //@name RegisterRequest

func (r registerReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

func (h *Handler) Register(c *gin.Context) {
	var input registerReq

	if ok := bindData(c, &input); !ok {
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mesage": "user created"})
}

func (h *Handler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "current user"})
}
