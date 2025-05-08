package http

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/handlers"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/logger"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/savsgio/atreugo/v11"
)

func NewServer() *atreugo.Atreugo {
	cfg, _ := config.LoadConfig()
	config := atreugo.Config{
		Addr: cfg.App.Port,
	}

	server := atreugo.New(config)
	server.UseBefore(func(ctx *atreugo.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, X-Custom")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "OPTIONS,GET,HEAD,PUT,PATCH,POST,DELETE")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "Content-Length, Authorization")
		return ctx.Next()
	})
	server.UseBefore(logger.Handler)

	deps := di.NewDependency(cfg)
	handlers.InitRoutes(server, deps)

	return server
}
