package di

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/mongo/repositories"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/domain"
)

type Dependency struct {
	Config       *config.Config
	Repositories database.Repository
	Services     services.IServices
	Actor        *domain.User
}

func NewDependency(cfg *config.Config) *Dependency {
	repositories, err := repositories.NewMongoRepository(cfg)
	if err != nil {
		panic(err)
	}
	services := services.NewServices(repositories, cfg)

	return &Dependency{
		Config:       cfg,
		Repositories: repositories,
		Services:     services,
	}
}
