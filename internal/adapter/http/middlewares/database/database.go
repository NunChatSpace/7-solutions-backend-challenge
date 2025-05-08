package database

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/savsgio/atreugo/v11"
)

type CtxKey struct{}

func FromContext(ctx *atreugo.RequestCtx) database.Repository {
	return ctx.Value(&CtxKey{}).(database.Repository)
}

func Handler(ctx *atreugo.RequestCtx, db database.Repository) error {
	ctx.RequestCtx.SetUserValue(&CtxKey{}, db)
	return ctx.Next()
}
