package middleware

import (
	"log/slog"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
)

// AuthUser checks if the request contains a valid session
// and saves the session's userId in the context
func AuthUser(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("userId")

		if id == nil {
			e := apperrors.NewAuthorization(apperrors.InvalidSession)
			c.JSON(e.Status(), gin.H{"error": e})
			c.Abort()
			return
		}

		userId := id.(string)

		c.Set("userId", userId)

		// Recreate session to extend its lifetime
		session.Set("userId", id)
		if err := session.Save(); err != nil {
			logger.Error("Failed to create session", slog.AnyValue(err.Error()))
		}

		c.Next()
	}
}
