package services

import (
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	sessionservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/session_services"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
)

func ProvideServices(deps *di.Dependency) {
	di.Provide(deps, userservices.NewUserService(deps))
	di.Provide(deps, authservices.NewAuthService(deps))
	di.Provide(deps, sessionservices.NewSessionService(deps))
}
