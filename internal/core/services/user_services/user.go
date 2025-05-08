package userservices

import (
	dbrepo "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type Port interface {
	GetUserByID(id string) (*domain.User, error)
	SearchUsers(user domain.User) ([]*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(id string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
}

type userService struct {
	Repository dbrepo.Repository
}

func NewUserService(repo dbrepo.Repository) Port {
	return &userService{
		Repository: repo,
	}
}

func (s *userService) GetUserByID(id string) (*domain.User, error) {
	// Implementation for getting a user by ID
	user, err := s.Repository.User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *userService) SearchUsers(user domain.User) ([]*domain.User, error) {
	// Implementation for searching users
	users, err := s.Repository.User().Search(user)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s *userService) CreateUser(user *domain.User) error {
	if err := s.Repository.User().InsertUser(user); err != nil {
		return err
	}
	return nil
}
func (s *userService) UpdateUser(id string, user *domain.User) (*domain.User, error) {
	res, err := s.Repository.User().UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (s *userService) DeleteUser(id string) error {
	if err := s.Repository.User().DeleteUser(id); err != nil {
		return err
	}
	return nil
}
