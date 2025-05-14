package userservices

import (
	"fmt"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserByID(id string) (*domain.UserResponse, error)
	SearchUsers(user domain.User) ([]*domain.UserResponse, error)
	SearchUsersForAuth(user domain.User) ([]*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(id string, user *domain.User) error
	DeleteUser(id string) error

	Authenticate(user *domain.User) error
}

type userService struct {
	Dependencies *di.Dependency

	userRepo database.IUserRepository
}

func NewUserService(deps *di.Dependency) IUserService {
	return userService{
		Dependencies: deps,
		userRepo:     di.Get[database.IUserRepository](deps),
	}
}

func (s userService) GetUserByID(id string) (*domain.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s userService) SearchUsers(user domain.User) ([]*domain.UserResponse, error) {
	users, err := s.userRepo.Search(user)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s userService) SearchUsersForAuth(user domain.User) ([]*domain.User, error) {
	users, err := s.userRepo.SearchForAuth(user)
	if err != nil {
		return nil, err
	}

	return users, nil
}
func (s userService) CreateUser(user *domain.User) error {
	if err := s.userRepo.InsertUser(user); err != nil {
		return err
	}
	return nil
}
func (s userService) UpdateUser(id string, user *domain.User) error {
	if err := s.userRepo.UpdateUser(id, user); err != nil {
		return err
	}

	return nil
}
func (s userService) DeleteUser(id string) error {
	if err := s.userRepo.DeleteUser(id); err != nil {
		return err
	}
	return nil
}

func (s userService) Authenticate(user *domain.User) error {
	users, err := s.SearchUsersForAuth(domain.User{
		Email: user.Email,
	})
	if err != nil {
		return err
	}

	for _, u := range users {
		if err := s.comparePassword(*user.Password, *u.Password); err != nil {
			continue
		}

		*user = *u
		break
	}

	if user.ID == nil {
		return fmt.Errorf("invalid email or password")
	}

	return nil
}

func (s userService) comparePassword(plainPassword string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return fmt.Errorf("invalid password")
	}
	return nil
}
