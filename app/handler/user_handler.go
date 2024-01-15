package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
	"github.com/opchaves/gin-web-app/app/service"
	"github.com/opchaves/gin-web-app/app/utils"
)

// TODO login view, form, handler, use htmx

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

	h.setUserSession(c, user.ID.String())

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *Handler) Login(c *gin.Context) {
	var req service.LoginInput

	if err := c.ShouldBind(&req); err != nil {
		errors := parseError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	user, err := h.UserService.Login(c, &req)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{"error": err})
		return
	}

	h.setUserSession(c, user.ID.String())

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) Logout(c *gin.Context) {
	c.Set("user", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()

	if err != nil {
		h.Logger.Warn("error clearing sessions", slog.AnyValue(err))
	}

	c.JSON(http.StatusOK, true)
}

func (h *Handler) GetCurrent(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	user, err := h.UserService.GetById(c, userId)

	if err != nil {
		h.Logger.Info("Unable to find user", slog.AnyValue(err))
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{"error": e})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req service.ForgotPasswordInput

	if err := c.ShouldBind(&req); err != nil {
		errors := parseError(err)
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	user, err := h.UserService.GetByEmail(c, req.Email)

	if err != nil {
		// No user with the email found
		if err.Error() == apperrors.NewNotFound("email", req.Email).Error() {
			c.JSON(http.StatusOK, true)
			return
		}

		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{"error": e})
		return
	}

	ctx := c.Request.Context()
	err = h.UserService.ForgotPassword(ctx, user)

	if err != nil {
		h.Logger.Warn("error sending reset password email", slog.AnyValue(err))
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
