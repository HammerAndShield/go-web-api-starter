package apiutils

import (
	"go-web-api-starter/internal/common"
	"go-web-api-starter/internal/vcs"
	"log/slog"
	"os"
	"sync"
)

type ApiConfig struct {
	Port        int
	Env         string
	Version     string
	Wg          sync.WaitGroup
	Logger      *slog.Logger
	CorsOptions *Cors
}

type Cors struct {
	TrustedOrigins []string
}

type Option func(*ApiConfig)

func WithLoggerOptions(logger *slog.Logger) Option {
	return func(config *ApiConfig) {
		config.Logger = logger
	}
}

func WithCorsOptions(trustedOrigins []string) Option {
	return func(config *ApiConfig) {
		config.CorsOptions = &Cors{
			TrustedOrigins: trustedOrigins,
		}
	}
}

func NewApiConfig(getEnv func(string) string, portKey string, opts ...Option) *ApiConfig {
	env := common.StringEnv(getEnv, "ENV", "dev")

	cfg := &ApiConfig{
		Port:        common.IntEnv(getEnv, portKey, 8080),
		Env:         env,
		Logger:      slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		Version:     vcs.Version(),
		Wg:          sync.WaitGroup{},
		CorsOptions: &Cors{TrustedOrigins: []string{"https://*", "http://*"}},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
