package services

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/postgres/repositories"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/http/middlewares/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/savsgio/atreugo/v11"
)

type UserPort interface {
	GetUserByID(id int) (*domain.User, error)
	SearchUsers(user domain.User) ([]*domain.User, error)
	CreateUser(user *domain.User) (*domain.User, error)
	UpdateUser(id int, user *domain.User) (*domain.User, error)
	DeleteUser(id int) error
}

type UserService struct {
	Repository repositories.Repository
}

func NewUserService(ctx *atreugo.RequestCtx) UserPort {
	repo := database.FromContext(ctx)
	return &UserService{
		Repository: repo,
	}
}

func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	// Implementation for getting a user by ID
	user, err := s.Repository.User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *UserService) SearchUsers(user domain.User) ([]*domain.User, error) {
	// Implementation for searching users
	users, err := s.Repository.User().Search(user)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s *UserService) CreateUser(user *domain.User) (*domain.User, error) {
	// Implementation for creating a user
	return nil, nil
}
func (s *UserService) UpdateUser(id int, user *domain.User) (*domain.User, error) {
	// Implementation for updating a user
	return nil, nil
}
func (s *UserService) DeleteUser(id int) error {
	// Implementation for deleting a user
	return nil
}
