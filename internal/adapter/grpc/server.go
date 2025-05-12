package grpc

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/gen/sessionpb"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/gen/userpb"
	interceptors "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/interceptors"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(deps *di.Dependency) *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor(deps)))
	reflection.Register(s)
	userpb.RegisterUserServiceServer(s, services.NewUserServiceServer(deps))
	sessionpb.RegisterSessionServiceServer(s, services.NewSessionServiceServer(deps))

	return s
}
