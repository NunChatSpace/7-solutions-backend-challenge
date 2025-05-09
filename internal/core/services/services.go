package services

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	sessionservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/session_services"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
)

type IServices interface {
	User() userservices.Port
	Auth() authservices.Port
	Session() sessionservices.Port
}

type ServicesImpl struct {
	user    userservices.Port
	auth    authservices.Port
	session sessionservices.Port
}

func (s *ServicesImpl) User() userservices.Port {
	return s.user
}
func (s *ServicesImpl) Auth() authservices.Port {
	return s.auth
}
func (s *ServicesImpl) Session() sessionservices.Port {
	return s.session
}

func NewServices(repo database.Repository) IServices {
	return &ServicesImpl{
		user: userservices.NewUserService(repo),
		auth: authservices.NewAuthService(repo),
	}
}
