package interceptors

import (
	"context"
	"strings"

	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(deps *di.Dependency) grpc.UnaryServerInterceptor {
	skipAuth := map[string]bool{
		"/session.SessionService/Login": true, // Fully-qualified method name
		"/user.UserService/CreateUser":  true,
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if skipAuth[info.FullMethod] {
			// Skip auth check
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization token")
		}

		tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
		authService := di.Get[authservices.IAuthSerivce](deps)
		if _, err := authService.DecodeToken(tokenStr); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}
