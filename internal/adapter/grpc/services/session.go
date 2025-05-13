package services

import (
	"context"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/gen/sessionpb"
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
	userService := s.Dependencies.Services.User()
	sessionService := s.Dependencies.Services.Session()
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
