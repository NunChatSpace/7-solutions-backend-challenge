package di

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/mongo/repositories"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
)

type Dependency struct {
	Repositories database.Repository
	Services     services.IServices
}

func NewDependency(cfg *config.Config) *Dependency {
	repositories, err := repositories.NewMongoRepository(cfg)
	if err != nil {
		panic(err)
	}
	services := services.NewServices(repositories)

	return &Dependency{
		Repositories: repositories,
		Services:     services,
	}
}
