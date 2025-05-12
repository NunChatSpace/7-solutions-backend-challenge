package interceptors

import (
	"context"
	"fmt"
	"strings"

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
		fmt.Println("AuthInterceptor: ", info.FullMethod)
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
		if _, err := deps.Services.Auth().DecodeToken(tokenStr); err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		return handler(ctx, req)
	}
}
