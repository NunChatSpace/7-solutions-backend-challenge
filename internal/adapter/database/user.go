package database

import "github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"

type IUserRepository interface {
	InsertUser(user *domain.User) error
	GetUserByID(id string) (*domain.User, error)
	Search(user domain.User) ([]*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id string) error
}
