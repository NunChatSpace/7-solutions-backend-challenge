package services

import (
	"context"
	"fmt"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/grpc/gen/userpb"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/common"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type userServiceServer struct {
	Dependencies *di.Dependency // Your business logic service
	userpb.UnimplementedUserServiceServer
}

func NewUserServiceServer(deps *di.Dependency) userpb.UserServiceServer {
	return &userServiceServer{Dependencies: deps}
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	user := domain.User{
		Email:    &req.Email,
		Password: &req.Password,
	}
	if err := s.Dependencies.Services.User().CreateUser(&user); err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.UserResponse{
			Id:        common.SafeString(user.ID),
			Name:      common.SafeString(user.Name),
			Email:     common.SafeString(user.Email),
			CreatedAt: common.SafeTime(user.CreatedAt),
			UpdatedAt: common.SafeTime(user.UpdatedAt),
			Scopes:    s.convertScopes(common.SafeMap(user.Scopes)),
		},
	}, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := s.Dependencies.Services.User().GetUserByID(req.Id)
	if err != nil {
		return nil, err
	}

	fmt.Printf("grpc user: %v\n", user)
	return &userpb.GetUserResponse{
		User: &userpb.UserResponse{
			Id:        common.SafeString(user.ID),
			Name:      common.SafeString(user.Name),
			Email:     common.SafeString(user.Email),
			CreatedAt: common.SafeTime(user.CreatedAt),
			UpdatedAt: common.SafeTime(user.UpdatedAt),
			Scopes:    s.convertScopes(common.SafeMap(user.Scopes)),
		},
	}, nil
}

func (s userServiceServer) convertScopes(scopes map[string]interface{}) map[string]int32 {
	result := make(map[string]int32)
	for key, val := range scopes {
		switch v := val.(type) {
		case int:
			result[key] = int32(v)
		case int32:
			result[key] = v
		case int64:
			result[key] = int32(v)
		case float64:
			result[key] = int32(v) // common when unmarshaling from JSON
		default:
			// skip or log unexpected type
			continue
		}
	}
	return result
}
