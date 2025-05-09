package sessionservices

import (
	dbrepo "github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	authservices "github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services/auth_services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type Port interface {
	CreateSession(userID string) (*domain.Tokens, error)
	GetSessionByID(id string) (*domain.Session, error)
	TerminateSession(id string) error
}

type sessionService struct {
	Repository dbrepo.Repository
}

func NewSessionService(repo dbrepo.Repository) Port {
	return &sessionService{
		Repository: repo,
	}
}

func (s *sessionService) CreateSession(userID string) (*domain.Tokens, error) {
	// Implementation for creating a session
	session := &domain.Session{
		UserID: &userID,
	}
	_session, err := s.Repository.Session().InsertSession(session)
	if err != nil {
		return nil, err
	}

	jwtService := authservices.NewAuthService(s.Repository)
	accessToken, refreshToken, err := jwtService.GenerateTokens(userID, *_session.ID)
	if err != nil {
		return nil, err
	}

	return &domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *sessionService) GetSessionByID(id string) (*domain.Session, error) {
	// Implementation for getting a session by ID
	session, err := s.Repository.Session().GetSessionByID(id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionService) TerminateSession(id string) error {
	// Implementation for terminating a session
	if err := s.Repository.Session().TerminateSession(id); err != nil {
		return err
	}
	return nil
}
