package test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/opchaves/gin-web-app/app"
	"github.com/opchaves/gin-web-app/app/model"
	"github.com/stretchr/testify/assert"
)

func SetupTest(t *testing.T) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	_ = godotenv.Load("../../.env.test")

	srv, err := app.Setup()
	assert.NoError(t, err)

	cleanUpDatabase(t, srv)

	return srv.Router
}

func cleanUpDatabase(t *testing.T, config *app.Config) {
	queries := model.New(config.Db)

	err := queries.DeleteWorkspaces(config.Ctx)
	assert.NoError(t, err)
	err = queries.DeleteUsers(config.Ctx)
	assert.NoError(t, err)
}
