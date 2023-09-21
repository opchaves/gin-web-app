package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl    string `env:"DATABASE_URL,required"`
	Port           string `env:"PORT,default=8080"`
	CorsOrigin     string `env:"CORS_ORIGIN,default=*"`
	HandlerTimeOut int64  `env:"HANDLER_TIMEOUT,default=5"`
	MaxBodyBytes   int64  `env:"MAX_BODY_BYTES,default=4194304"`
	RootPath       string `env:"ROOT_PATH,default=src/github.com/opchaves/gin-web-app"`
	TemplatesGlob  string `env:"TEMPLATES_GLOB,default=app/templates/**/*"`
	AssetsDir      string `env:"ASSETS_DIR,default=assets"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	config.TemplatesGlob = filepath.Join(os.Getenv("GOPATH"), config.RootPath, config.TemplatesGlob)
	config.AssetsDir = filepath.Join(os.Getenv("GOPATH"), config.RootPath, config.AssetsDir)

	return
}
