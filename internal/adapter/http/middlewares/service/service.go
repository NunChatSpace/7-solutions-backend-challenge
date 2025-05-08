package service

import (
	coreservice "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
	"github.com/savsgio/atreugo/v11"
)

type CtxKey struct{}

func FromContext(ctx *atreugo.RequestCtx) coreservice.IServices {
	return ctx.Value(&CtxKey{}).(coreservice.IServices)
}

func Handler(ctx *atreugo.RequestCtx, services coreservice.IServices) error {
	ctx.RequestCtx.SetUserValue(&CtxKey{}, services)
	return ctx.Next()
}
