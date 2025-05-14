package testutils

import (
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database/mongo/repositories"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/core/services"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
)

func NewTestDependency(cfg *config.Config) *di.Dependency {
	deps := di.NewDependency(cfg)

	repositories.ProvideRepositories(deps)
	services.ProvideServices(deps)

	return deps
}
