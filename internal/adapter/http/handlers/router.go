package handlers

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/handlers/users"
	"github.com/savsgio/atreugo/v11"
)

func InitRoutes(server *atreugo.Atreugo) {
	server.GET("/", func(ctx *atreugo.RequestCtx) error {
		return ctx.TextResponse("Hello World!")
	})

	apiRouter := server.NewGroupPath("/api")

	users.InitRoutes(apiRouter)
}
