package userservices

import (
	"fmt"

	dbrepo "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type Port interface {
	GetUserByID(id string) (*domain.UserResponse, error)
	SearchUsers(user domain.User) ([]*domain.UserResponse, error)
	SearchUsersForAuth(user domain.User) ([]*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(id string, user *domain.User) (*domain.User, error)
	DeleteUser(id string) error

	Authenticate(user *domain.User) (*domain.User, error)
}

type userService struct {
	Repository dbrepo.Repository
	Config     *config.Config
}

func NewUserService(repo dbrepo.Repository, cfg *config.Config) Port {
	return &userService{
		Repository: repo,
		Config:     cfg,
	}
}

func (s *userService) GetUserByID(id string) (*domain.UserResponse, error) {
	// Implementation for getting a user by ID
	user, err := s.Repository.User().GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *userService) SearchUsers(user domain.User) ([]*domain.UserResponse, error) {
	// Implementation for searching users
	users, err := s.Repository.User().Search(user)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userService) SearchUsersForAuth(user domain.User) ([]*domain.User, error) {
	// Implementation for searching users
	users, err := s.Repository.User().SearchForAuth(user)
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

func (s *userService) Authenticate(user *domain.User) (*domain.User, error) {
	users, err := s.SearchUsersForAuth(domain.User{
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}
	var authenticatedUser *domain.User
	for _, u := range users {
		if err := s.comparePassword(*user.Password, *u.Password); err != nil {
			continue
		}

		authenticatedUser = u
		break
	}

	if authenticatedUser == nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	return users[0], nil
}

func (s userService) comparePassword(plainPassword string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}
