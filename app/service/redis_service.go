package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
	"github.com/redis/go-redis/v9"
)

// RedisService defines methods related to the redis db the service layer expects
// any service it interacts with to implement
type RedisService interface {
	SetResetToken(ctx context.Context, id string) (string, error)
}

type redisService struct {
	Logger *slog.Logger
	Db     *pgxpool.Pool
	Redis  *redis.Client
}

type RDConfig struct {
	Logger *slog.Logger
	Db     *pgxpool.Pool
	Redis  *redis.Client
}

func NewRedisService(c *RDConfig) RedisService {
	return &redisService{
		Logger: c.Logger,
		Db:     c.Db,
		Redis:  c.Redis,
	}
}

// Redis Prefixes
const (
	ForgotPasswordPrefix = "forgot-password"
)

// SetResetToken implements RedisService.
func (s *redisService) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := gonanoid.New()

	if err != nil {
		s.Logger.Error("failed to generate id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = s.Redis.Set(ctx, fmt.Sprintf("%s:%s", ForgotPasswordPrefix, uid), id, 24*time.Hour).Err(); err != nil {
		s.Logger.Error("failed to set link in redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return uid, err
}
