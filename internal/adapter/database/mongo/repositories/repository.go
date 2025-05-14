package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/di"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryImpl struct {
	userRepo    database.IUserRepository
	sessionRepo database.ISessionRepository
}

func NewMongoRepository(cfg *config.Config) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=admin",
		cfg.Database.MongoDB.Username,
		cfg.Database.MongoDB.Password,
		cfg.Database.MongoDB.Host,
		cfg.Database.MongoDB.Port,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(fmt.Errorf("failed to connect to MongoDB: %w", err))
	}

	if err := client.Ping(ctx, nil); err != nil {
		panic(fmt.Errorf("failed to ping MongoDB: %w", err))
	}

	return client.Database(cfg.Database.MongoDB.DatabaseName)
}

func (r *RepositoryImpl) SetUserRepo(userRepo database.IUserRepository) {
	r.userRepo = userRepo
}

func (r *RepositoryImpl) SetSessionRepo(sessionRepo database.ISessionRepository) {
	r.sessionRepo = sessionRepo
}

func ProvideRepositories(deps *di.Dependency) {
	cfg := di.Get[*config.Config](deps)
	repositories := NewMongoRepository(cfg)

	di.Provide(deps, repositories)
	di.Provide(deps, NewUserRepository(repositories))
	di.Provide(deps, NewSessionRepository(repositories))
}
