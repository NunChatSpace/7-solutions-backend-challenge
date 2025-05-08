package config

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/savsgio/atreugo/v11"
)

type ctxKey struct{}

func FromContext(ctx *atreugo.RequestCtx) *config.Config {
	if ctx.Value(&ctxKey{}) == nil {
		cfg, _ := config.LoadConfig()
		return cfg
	}

	return ctx.Value(&ctxKey{}).(*config.Config)
}

func Handler(ctx *atreugo.RequestCtx, cfg *config.Config) error {
	ctx.RequestCtx.SetUserValue(&ctxKey{}, cfg)
	return ctx.Next()
}
