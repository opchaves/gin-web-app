package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opchaves/gin-web-app/app"
)

func main() {
	config, err := app.Setup()

	if err != nil {
		slog.Error("error: ", slog.AnyValue(err))
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:    ":" + config.Cfg.Port,
		Handler: config.Router,
	}

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			config.Logger.Error("failed to initialize server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	config.Logger.Debug(fmt.Sprintf("Listening on port %v", srv.Addr))

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
	config.Logger.Debug("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		config.Logger.Debug("Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	return
}
