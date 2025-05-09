package database

import "github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"

type ISessionRepository interface {
	InsertSession(session *domain.Session) (*domain.Session, error)
	GetSessionByID(id string) (*domain.Session, error)
	TerminateSession(id string) error
}
