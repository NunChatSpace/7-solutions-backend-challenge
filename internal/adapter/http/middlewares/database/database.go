package database

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/postgres/repositories"
	"github.com/savsgio/atreugo/v11"
)

type CtxKey struct{}

func FromContext(ctx *atreugo.RequestCtx) repositories.Repository {
	return ctx.Value(&CtxKey{}).(repositories.Repository)
}

func Handler(ctx *atreugo.RequestCtx, db repositories.Repository) error {
	ctx.RequestCtx.SetUserValue(&CtxKey{}, db)
	return ctx.Next()
}
