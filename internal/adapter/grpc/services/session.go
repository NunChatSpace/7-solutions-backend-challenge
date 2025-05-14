package services

import (
	"context"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/gen/sessionpb"
	sessionservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/session_services"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type sessionServiceServer struct {
	Dependencies *di.Dependency // Your business logic service
	sessionpb.UnimplementedSessionServiceServer
}

func NewSessionServiceServer(deps *di.Dependency) sessionpb.SessionServiceServer {
	return &sessionServiceServer{
		Dependencies: deps,
	}
}

func (s *sessionServiceServer) Login(ctx context.Context, req *sessionpb.LoginRequest) (*sessionpb.LoginResponse, error) {
	userService := di.Get[userservices.IUserService](s.Dependencies)
	sessionService := di.Get[sessionservices.ISessionService](s.Dependencies)
	user := domain.User{
		Email:    &req.Email,
		Password: &req.Password,
	}
	if err := userService.Authenticate(&user); err != nil {
		return nil, err
	}
	session, err := sessionService.CreateSession(*user.ID)
	if err != nil {
		return nil, err
	}

	return &sessionpb.LoginResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}
