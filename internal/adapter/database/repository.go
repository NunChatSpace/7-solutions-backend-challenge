package database

type Repository interface {
	User() IUserRepository
	Session() ISessionRepository

	SetUserRepo(repo IUserRepository)
	SetSessionRepo(repo ISessionRepository)
}
