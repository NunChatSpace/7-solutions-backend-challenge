package userservices_test

import (
	"errors"
	"testing"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	"github.com/NunChatSpace/7-solutions-backend-challenge/mocks"
	"github.com/golang/mock/gomock"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	if _cfg, err := config.LoadConfig(); err != nil {
		panic(err)
	} else {
		cfg = _cfg
	}
}

func TestGetUserByID(t *testing.T) {
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1234"
	Name := "John Doe"
	Email := "john@doe.com"
	expectedUser := &domain.User{
		ID:    &id,
		Name:  &Name,
		Email: &Email,
	}
	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	mockRepo.EXPECT().User().Return(mockUserRepo).AnyTimes()
	mockUserRepo.EXPECT().GetUserByID(1).Return(expectedUser, nil).AnyTimes()
	mockUserRepo.EXPECT().GetUserByID(2).Return(nil, errors.New("user not found")).AnyTimes()

	userService := userservices.NewUserService(mockRepo, cfg)

	// Define the test cases
	tests := []struct {
		name     string
		userID   string
		expected *domain.User
		err      error
	}{
		{
			name:     "Valid User ID",
			userID:   "1",
			expected: expectedUser,
			err:      nil,
		},
		{
			name:     "User Not Found",
			userID:   "2",
			expected: nil,
			err:      errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.GetUserByID(tt.userID)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result != nil && *result.ID != *tt.expected.ID {
				t.Errorf("expected user ID %s, got %s", *tt.expected.ID, *result.ID)
			}
			if result != nil && *result.Name != *tt.expected.Name {
				t.Errorf("expected user Name %s, got %s", *tt.expected.Name, *result.Name)
			}
			if result != nil && *result.Email != *tt.expected.Email {
				t.Errorf("expected user Email %s, got %s", *tt.expected.Email, *result.Email)
			}
		})
	}
}

func TestSearchUsers(t *testing.T) {
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1234"
	Name := "John Doe"
	Email := ""
	expectedUser := &domain.User{
		ID:    &id,
		Name:  &Name,
		Email: &Email,
	}
	expectedUsers := []*domain.User{expectedUser}
	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRepo.EXPECT().User().Return(mockUserRepo).AnyTimes()
	mockUserRepo.EXPECT().Search(*expectedUser).Return(expectedUsers, nil).AnyTimes()
	mockUserRepo.EXPECT().Search(domain.User{Name: &Name}).Return(nil, errors.New("user not found")).AnyTimes()
	mockUserRepo.EXPECT().Search(domain.User{Email: &Email}).Return(nil, errors.New("user not found")).AnyTimes()

	userService := userservices.NewUserService(mockRepo, cfg)
	// Define the test cases
	tests := []struct {
		name     string
		user     domain.User
		expected []*domain.User
		err      error
	}{
		{
			name:     "Valid User",
			user:     *expectedUser,
			expected: expectedUsers,
			err:      nil,
		},
		{
			name:     "User Not Found",
			user:     domain.User{Name: &Name},
			expected: nil,
			err:      errors.New("user not found"),
		},
		{
			name:     "User Not Found",
			user:     domain.User{Email: &Email},
			expected: nil,
			err:      errors.New("user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.SearchUsers(tt.user)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result != nil && len(result) != len(tt.expected) {
				t.Errorf("expected %d users, got %d", len(tt.expected), len(result))
			}
			if result != nil && *result[0].ID != *tt.expected[0].ID {
				t.Errorf("expected user ID %s, got %s", *tt.expected[0].ID, *result[0].ID)
			}
			if result != nil && *result[0].Name != *tt.expected[0].Name {
				t.Errorf("expected user Name %s, got %s", *tt.expected[0].Name, *result[0].Name)
			}
			if result != nil && *result[0].Email != *tt.expected[0].Email {
				t.Errorf("expected user Email %s, got %s", *tt.expected[0].Email, *result[0].Email)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := "1234"
	Name := "John Doe"
	Email := ""
	expectedUser := &domain.User{
		ID:    &id,
		Name:  &Name,
		Email: &Email,
	}
	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRepo.EXPECT().User().Return(mockUserRepo).AnyTimes()
	mockUserRepo.EXPECT().InsertUser(expectedUser).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().InsertUser(&domain.User{Name: &Name}).Return(errors.New("user not found")).AnyTimes()

	userService := userservices.NewUserService(mockRepo, cfg)
	// Define the test cases
	tests := []struct {
		name string
		user *domain.User
		err  error
	}{
		{
			name: "Valid User",
			user: expectedUser,
			err:  nil,
		},
		{
			name: "User Not Found",
			user: &domain.User{Name: &Name},
			err:  errors.New("user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userService.CreateUser(tt.user)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1234"
	Name := "John Doe"
	Email := ""
	expectedUser := &domain.User{
		ID:    &id,
		Name:  &Name,
		Email: &Email,
	}
	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRepo.EXPECT().User().Return(mockUserRepo).AnyTimes()
	mockUserRepo.EXPECT().UpdateUser(expectedUser).Return(expectedUser, nil).AnyTimes()
	mockUserRepo.EXPECT().UpdateUser(&domain.User{Name: &Name}).Return(nil, errors.New("user not found")).AnyTimes()

	userService := userservices.NewUserService(mockRepo, cfg)
	// Define the test cases
	tests := []struct {
		name     string
		userID   string
		user     *domain.User
		expected *domain.User
		err      error
	}{
		{
			name:     "Valid User",
			userID:   "1",
			user:     expectedUser,
			expected: expectedUser,
			err:      nil,
		},
		{
			name:     "User Not Found",
			userID:   "2",
			user:     &domain.User{Name: &Name},
			expected: nil,
			err:      errors.New("user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := userService.UpdateUser(tt.userID, tt.user)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if result != nil && *result.ID != *tt.expected.ID {
				t.Errorf("expected user ID %d, got %s", *tt.expected.ID, *result.ID)
			}
			if result != nil && *result.Name != *tt.expected.Name {
				t.Errorf("expected user Name %s, got %s", *tt.expected.Name, *result.Name)
			}
			if result != nil && *result.Email != *tt.expected.Email {
				t.Errorf("expected user Email %s, got %s", *tt.expected.Email, *result.Email)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockRepo.EXPECT().User().Return(mockUserRepo).AnyTimes()
	mockUserRepo.EXPECT().DeleteUser(1).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().DeleteUser(2).Return(errors.New("user not found")).AnyTimes()

	userService := userservices.NewUserService(mockRepo, cfg)
	// Define the test cases
	tests := []struct {
		name   string
		userID string
		err    error
	}{
		{
			name:   "Valid User",
			userID: "1",
			err:    nil,
		},
		{
			name:   "User Not Found",
			userID: "2",
			err:    errors.New("user not found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userService.DeleteUser(tt.userID)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
