package database

type Repository interface {
	User() IUserRepository
	Session() ISessionRepository
}
