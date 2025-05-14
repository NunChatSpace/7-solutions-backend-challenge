package sessionservices_test

import (
	"errors"
	"os"
	"testing"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	sessionservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/session_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
	testutils "github.com/NunChatSpace/7-solutions-backend-challenge/internal/test_utils"
	"github.com/NunChatSpace/7-solutions-backend-challenge/mocks"
	"github.com/golang/mock/gomock"
)

var cfg *config.Config

func TestMain(m *testing.M) {
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestCreateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)

		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)
		mockUserRepo := mocks.NewMockIUserRepository(ctrl)
		mockAuthService := mocks.NewMockIAuthSerivce(ctrl)

		userID := "u-123"
		sessionID := "s-456"
		scopes := map[string]interface{}{"users": 15}
		createdSession := &domain.Session{ID: &sessionID, UserID: &userID}
		user := &domain.UserResponse{ID: &userID, Scopes: &scopes}

		mockSessionRepo.EXPECT().
			InsertSession(gomock.Any()).
			Return(createdSession, nil)

		mockUserRepo.EXPECT().
			GetUserByID(userID).
			Return(user, nil)

		mockAuthService.EXPECT().
			GenerateTokens(gomock.Any()).
			Return("access-token", "refresh-token", nil)

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)
		di.Provide[database.IUserRepository](deps, mockUserRepo)
		di.Provide[authservices.IAuthSerivce](deps, mockAuthService)

		service := sessionservices.NewSessionService(deps)

		tokens, err := service.CreateSession(userID)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if tokens.AccessToken == "" || tokens.RefreshToken == "" {
			t.Error("expected non-empty tokens")
		}
	})

	t.Run("insert session fails", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)

		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

		mockSessionRepo.EXPECT().
			InsertSession(gomock.Any()).
			Return(nil, errors.New("insert failed"))

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)

		service := sessionservices.NewSessionService(deps)

		_, err := service.CreateSession("u-123")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("get user fails", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)

		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)
		mockUserRepo := mocks.NewMockIUserRepository(ctrl)

		userID := "u-123"
		sessionID := "s-456"
		createdSession := &domain.Session{ID: &sessionID, UserID: &userID}

		mockSessionRepo.EXPECT().
			InsertSession(gomock.Any()).
			Return(createdSession, nil)

		mockUserRepo.EXPECT().
			GetUserByID(userID).
			Return(nil, errors.New("user not found"))

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)
		di.Provide[database.IUserRepository](deps, mockUserRepo)

		service := sessionservices.NewSessionService(deps)

		_, err := service.CreateSession(userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("generate tokens fails", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)

		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)
		mockUserRepo := mocks.NewMockIUserRepository(ctrl)
		mockAuthService := mocks.NewMockIAuthSerivce(ctrl)

		userID := "u-123"
		sessionID := "s-456"
		scopes := map[string]interface{}{"users": 15}
		createdSession := &domain.Session{ID: &sessionID, UserID: &userID}
		user := &domain.UserResponse{ID: &userID, Scopes: &scopes}

		mockSessionRepo.EXPECT().
			InsertSession(gomock.Any()).
			Return(createdSession, nil)

		mockUserRepo.EXPECT().
			GetUserByID(userID).
			Return(user, nil)

		mockAuthService.EXPECT().
			GenerateTokens(gomock.Any()).
			Return("", "", errors.New("token generation failed"))

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)
		di.Provide[database.IUserRepository](deps, mockUserRepo)
		di.Provide[authservices.IAuthSerivce](deps, mockAuthService)

		service := sessionservices.NewSessionService(deps)

		_, err := service.CreateSession(userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCreateSession_InsertFails(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

	userID := "u-123"
	mockSessionRepo.EXPECT().
		InsertSession(gomock.Any()).
		Return(nil, errors.New("insert failed"))

	di.Provide[database.ISessionRepository](deps, mockSessionRepo)
	service := sessionservices.NewSessionService(deps)

	_, err := service.CreateSession(userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSessionByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("success", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)
		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

		sessionID := "s-123"
		session := &domain.Session{ID: &sessionID}

		mockSessionRepo.EXPECT().
			GetSessionByID(sessionID).
			Return(session, nil)

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)

		service := sessionservices.NewSessionService(deps)

		result, err := service.GetSessionByID(sessionID)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if result.ID == nil || *result.ID != sessionID {
			t.Errorf("expected session ID %s, got %v", sessionID, result.ID)
		}
	})

	t.Run("repo returns error", func(t *testing.T) {
		deps := testutils.NewTestDependency(cfg)
		mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

		sessionID := "s-123"

		mockSessionRepo.EXPECT().
			GetSessionByID(sessionID).
			Return(nil, errors.New("not found"))

		di.Provide[database.ISessionRepository](deps, mockSessionRepo)

		service := sessionservices.NewSessionService(deps)

		_, err := service.GetSessionByID(sessionID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestTerminateSession(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

	sessionID := "s-789"
	mockSessionRepo.EXPECT().TerminateSession(sessionID).Return(nil)
	di.Provide[database.ISessionRepository](deps, mockSessionRepo)

	service := sessionservices.NewSessionService(deps)

	err := service.TerminateSession(sessionID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestTerminateSession_Error(t *testing.T) {
	deps := testutils.NewTestDependency(cfg)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionRepo := mocks.NewMockISessionRepository(ctrl)

	sessionID := "s-789"
	mockSessionRepo.EXPECT().TerminateSession(sessionID).Return(errors.New("not found"))
	di.Provide[database.ISessionRepository](deps, mockSessionRepo)
	service := sessionservices.NewSessionService(deps)

	err := service.TerminateSession(sessionID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
