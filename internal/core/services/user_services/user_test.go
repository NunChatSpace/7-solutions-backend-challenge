package userservices_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/common"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	userservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/user_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	testutils "github.com/NunChatSpace/7-solutions-backend-challenge/internal/test_utils"
	"github.com/NunChatSpace/7-solutions-backend-challenge/mocks"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run()) // This is crucial to actually run the tests
}

func TestGetUserByID(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := "1234"
	Name := "John Doe"
	Email := "john@doe.com"
	expectedUser := &domain.UserResponse{
		ID:    &id,
		Name:  &Name,
		Email: &Email,
	}

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockUserRepo.EXPECT().GetUserByID("1").Return(expectedUser, nil).AnyTimes()
	mockUserRepo.EXPECT().GetUserByID("2").Return(nil, errors.New("user not found")).AnyTimes()

	di.Provide[database.IUserRepository](deps, mockUserRepo)

	userService := userservices.NewUserService(deps)

	// Define the test cases
	tests := []struct {
		name     string
		userID   string
		expected *domain.UserResponse
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
	deps := testutils.NewTestDependency(cfg)
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
	expectedUsers := []*domain.UserResponse{
		&domain.UserResponse{
			ID:    &id,
			Name:  &Name,
			Email: &Email,
		},
	}

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockUserRepo.EXPECT().Search(*expectedUser).Return(expectedUsers, nil).AnyTimes()
	mockUserRepo.EXPECT().Search(domain.User{Name: &Name}).Return(nil, errors.New("user not found")).AnyTimes()
	mockUserRepo.EXPECT().Search(domain.User{Email: &Email}).Return(nil, errors.New("user not found")).AnyTimes()

	di.Provide[database.IUserRepository](deps, mockUserRepo)
	userService := userservices.NewUserService(deps)
	// Define the test cases
	tests := []struct {
		name     string
		user     domain.User
		expected []*domain.UserResponse
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
	deps := testutils.NewTestDependency(cfg)
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

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockUserRepo.EXPECT().InsertUser(expectedUser).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().InsertUser(&domain.User{Name: &Name}).Return(errors.New("user not found")).AnyTimes()
	di.Provide[database.IUserRepository](deps, mockUserRepo)

	userService := userservices.NewUserService(deps)
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
	deps := testutils.NewTestDependency(cfg)
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
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockUserRepo.EXPECT().UpdateUser(gomock.Any(), expectedUser).Return(nil).AnyTimes()
	mockUserRepo.EXPECT().UpdateUser(gomock.Any(), &domain.User{Name: &Name}).Return(errors.New("user not found")).AnyTimes()

	di.Provide[database.IUserRepository](deps, mockUserRepo)

	userService := userservices.NewUserService(deps)
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
			err := userService.UpdateUser(tt.userID, tt.user)
			if tt.err != nil && err == nil {
				t.Errorf("expected error %v, got nil", tt.err)
			}
			if tt.err == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	// Mock the database repository
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockUserRepo.EXPECT().DeleteUser("1").Return(nil).AnyTimes()
	mockUserRepo.EXPECT().DeleteUser("2").Return(errors.New("user not found")).AnyTimes()
	di.Provide[database.IUserRepository](deps, mockUserRepo)

	userService := userservices.NewUserService(deps)
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

func TestAuthenticate(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup
	mockUserRepo := mocks.NewMockIUserRepository(ctrl)

	// Password setup
	rawPassword := "secure123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)

	email := "test@example.com"
	name := "Test User"
	id := "1"

	mockedUser := &domain.User{
		ID:       &id,
		Name:     &name,
		Email:    &email,
		Password: common.Ptr(string(hashedPassword)),
	}

	// Test cases
	tests := []struct {
		name         string
		inputUser    *domain.User
		mockedResult []*domain.User
		mockedError  error
		expectedErr  string
	}{
		{
			name: "Successful authentication",
			inputUser: &domain.User{
				Email:    &email,
				Password: &rawPassword,
			},
			mockedResult: []*domain.User{mockedUser},
			mockedError:  nil,
			expectedErr:  "",
		},
		{
			name: "Wrong password",
			inputUser: &domain.User{
				Email:    &email,
				Password: common.Ptr("wrongpassword"),
			},
			mockedResult: []*domain.User{mockedUser},
			mockedError:  nil,
			expectedErr:  "invalid email or password",
		},
		{
			name: "User not found by email",
			inputUser: &domain.User{
				Email:    common.Ptr("unknown@example.com"),
				Password: &rawPassword,
			},
			mockedResult: nil,
			mockedError:  fmt.Errorf("user not found"),
			expectedErr:  "user not found",
		},
	}

	mockUserRepo.EXPECT().SearchForAuth(gomock.Any()).DoAndReturn(func(user domain.User) ([]*domain.User, error) {
		for _, tc := range tests {
			if user.Email != nil && *user.Email == *tc.inputUser.Email {
				// convert to []*domain.UserResponse
				var result []*domain.User
				for _, u := range tc.mockedResult {
					result = append(result, &domain.User{
						ID:       u.ID,
						Name:     u.Name,
						Email:    u.Email,
						Password: u.Password,
					})
				}
				return result, tc.mockedError
			}
		}
		return nil, fmt.Errorf("user not found")
	}).AnyTimes()

	di.Provide[database.IUserRepository](deps, mockUserRepo)

	service := userservices.NewUserService(deps)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Authenticate(tt.inputUser)
			if tt.expectedErr != "" {
				if err == nil || err.Error() != tt.expectedErr {
					t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			}
		})
	}
}
