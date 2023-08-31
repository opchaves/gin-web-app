package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opchaves/gin-web-app/app/config"
)

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	Db              *pgxpool.Pool
	Cfg             *config.Config
	Ctx             context.Context
	Logger          *slog.Logger
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

func Start(c *Config) error {
	corsConfig := cors.DefaultConfig()
	// corsConfig.AllowAllOrigins = false
	// corsConfig.AllowedOrigins = []string{cfg.CorsOrigin}

	router := gin.Default()
	router.Use(cors.New(corsConfig))
	router.LoadHTMLGlob("app/templates/**/*")
	router.Static("/assets", "./assets")

	SetRoutes(c, router)

	srv := &http.Server{
		Addr:    ":" + c.Cfg.Port,
		Handler: router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.Error("failed to initialize server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	c.Logger.Debug(fmt.Sprintf("Listening on port %v", srv.Addr))

	// Wait for kill signal of channel
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	c.Logger.Debug("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		c.Logger.Debug("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	return nil
}
