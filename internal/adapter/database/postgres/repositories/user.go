package repositories

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"gorm.io/gorm"
)

func (r *RepositoryImpl) User() IUserRepository {
	return r.userRepo
}

type IUserRepository interface {
	InsertUser(user *domain.User) error
	GetUserByID(id int) (*domain.User, error)
	Search(user domain.User) ([]*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id int) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db}
}
func (u *UserRepository) InsertUser(user *domain.User) error {
	// Implementation for inserting a user into the database
	return nil
}
func (u *UserRepository) GetUserByID(id int) (*domain.User, error) {
	// Implementation for finding a user by ID in the database
	return nil, nil
}
func (u *UserRepository) Search(user domain.User) ([]*domain.User, error) {
	// Implementation for searching users in the database
	return nil, nil
}
func (u *UserRepository) UpdateUser(user *domain.User) error {
	// Implementation for updating a user in the database
	return nil
}
func (u *UserRepository) DeleteUser(id int) error {
	// Implementation for deleting a user from the database
	return nil
}
