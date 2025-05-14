package sessionservices

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type ISessionService interface {
	CreateSession(userID string) (*domain.Tokens, error)
	GetSessionByID(id string) (*domain.Session, error)
	TerminateSession(id string) error
}

type sessionService struct {
	Dependencies *di.Dependency

	sessionRepo database.ISessionRepository
	userRepo    database.IUserRepository
	authservice authservices.IAuthSerivce
}

func NewSessionService(deps *di.Dependency) ISessionService {
	return sessionService{
		Dependencies: deps,
		sessionRepo:  di.Get[database.ISessionRepository](deps),
		userRepo:     di.Get[database.IUserRepository](deps),
		authservice:  di.Get[authservices.IAuthSerivce](deps),
	}
}

func (s sessionService) CreateSession(userID string) (*domain.Tokens, error) {
	session := &domain.Session{
		UserID: &userID,
	}
	_session, err := s.sessionRepo.InsertSession(session)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	tokenInfo := domain.TokenInfo{
		UserID:    userID,
		SessionID: *_session.ID,
		Scopes:    *user.Scopes,
	}
	accessToken, refreshToken, err := s.authservice.GenerateTokens(tokenInfo)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s sessionService) GetSessionByID(id string) (*domain.Session, error) {
	// Implementation for getting a session by ID
	session, err := s.sessionRepo.GetSessionByID(id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s sessionService) TerminateSession(id string) error {
	// Implementation for terminating a session
	if err := s.sessionRepo.TerminateSession(id); err != nil {
		return err
	}
	return nil
}
