package config

import (
	"context"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	DatabaseUrl    string `env:"DATABASE_URL,required"`
	RedisUrl       string `env:"REDIS_URL,required"`
	Port           string `env:"PORT,default=8080"`
	Domain         string `env:"DOMAIN,required"`
	CorsOrigin     string `env:"CORS_ORIGIN,default=*"`
	HandlerTimeOut int64  `env:"HANDLER_TIMEOUT,default=5"`
	MaxBodyBytes   int64  `env:"MAX_BODY_BYTES,default=4194304"`
	RootPath       string `env:"ROOT_PATH,default=src/github.com/opchaves/gin-web-app"`
	TemplatesGlob  string `env:"TEMPLATES_GLOB,default=app/templates/**/*"`
	AssetsDir      string `env:"ASSETS_DIR,default=assets"`
	SessionSecret  string `env:"SESSION_SECRET,default=sup3rs3cr37"`
	RateLimit      int64  `env:"RATE_LIMIT,default=1000"`

	MailMailer     string `env:"MAIL_MAILER,default=smtp"`
	MailHost       string `env:"MAIL_HOST,default=localhost"`
	MailPort       string `env:"MAIL_PORT,default=1025"`
	MailUsername   string `env:"MAIL_USERNAME"`
	MailPassword   string `env:"MAIL_PASSWORD"`
	MailEncryption string `env:"MAIL_ENCRYPTION"`
}

func LoadConfig(ctx context.Context) (config Config, err error) {
	err = envconfig.Process(ctx, &config)

	config.TemplatesGlob = filepath.Join(os.Getenv("GOPATH"), config.RootPath, config.TemplatesGlob)
	config.AssetsDir = filepath.Join(os.Getenv("GOPATH"), config.RootPath, config.AssetsDir)

	return
}
