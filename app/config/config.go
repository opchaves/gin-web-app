package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl    string `env:"DATABASE_URL,required`
	Port           string `env:"PORT",default=8080`
	CorsOrigin     string `env:"CORS_ORIGIN,default=*"`
	HandlerTimeOut int64  `env:"HANDLER_TIMEOUT,default=5"`
	MaxBodyBytes   int64  `env:"MAX_BODY_BYTES,default=4194304"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	return
}
