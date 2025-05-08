package services

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
)

type IServices interface {
	User() userservices.Port
}

type ServicesImpl struct {
	user userservices.Port
}

func (s *ServicesImpl) User() userservices.Port {
	return s.user
}

func NewServices(repo database.Repository) IServices {
	return &ServicesImpl{
		user: userservices.NewUserService(repo),
	}
}
