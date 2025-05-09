package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/adapter/database"
	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepositoryImpl struct {
	userRepo    database.IUserRepository
	sessionRepo database.ISessionRepository
}

func NewMongoRepository(cfg *config.Config) (database.Repository, error) {
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
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}
	fmt.Println("Connected to MongoDB with database:", cfg.Database.MongoDB.DatabaseName)
	db := client.Database(cfg.Database.MongoDB.DatabaseName)

	return &RepositoryImpl{
		userRepo: NewUserRepository(db),
	}, nil
}
