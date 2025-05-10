package handlers

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/handlers/sessions"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/handlers/users"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/savsgio/atreugo/v11"
)

func InitRoutes(server *atreugo.Atreugo, deps *di.Dependency) {
	server.GET("/", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("Hello World!")
	})
	server.GET("/health", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("OK")
	})

	apiRouter := server.NewGroupPath("/api/v1")

	users.InitRoutes(apiRouter, deps)
	sessions.InitRoutes(apiRouter, deps)
}
