package http

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/handlers"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/authen"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/logger"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/savsgio/atreugo/v11"
)

func NewServer(deps *di.Dependency, cfg atreugo.Config) *atreugo.Atreugo {
	server := atreugo.New(cfg)
	server.Static("/docs", "./static/swagger-ui")

	server.UseBefore(func(ctx *atreugo.RequestCtx) error {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type, X-Custom")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "OPTIONS,GET,HEAD,PUT,PATCH,POST,DELETE")
		ctx.Response.Header.Set("Access-Control-Expose-Headers", "Content-Length, Authorization")
		return ctx.Next()
	})
	server.UseBefore(func(rc *atreugo.RequestCtx) error {
		return logger.Handler(rc, deps)
	})
	server.UseBefore(func(rc *atreugo.RequestCtx) error {
		err := authen.Handler(rc, deps)
		if err != nil {
			return rc.ErrorResponse(err)
		}
		return rc.Next()
	})
	handlers.InitRoutes(server, deps)

	return server
}
