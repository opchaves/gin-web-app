package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"log/slog"

	"github.com/opchaves/gin-web-app/app/config"
)

// NOTE not used. current migrate command being used with Makefile

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)

	if len(os.Args) <= 3 {
		logger.Error(fmt.Sprintf("Usage:", os.Args[1], "command", "argument"))
		return errors.New("invalid command")
	}

	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Error("failed to load config", slog.Any("error", err))
	}

	switch os.Args[1] {
	case "migrate":
		err = Migrate(&cfg, logger, os.Args[3])
	case "seed":
		err = Seed(&cfg, logger, os.Args[3])
	default:
		err = errors.New("must specify a command")
	}

	return err
}

func Migrate(cfg *config.Config, logger *slog.Logger, command string) error {
	return nil
}

func Seed(cfg *config.Config, logger *slog.Logger, sqlFile string) error {
	return nil
}
