package repositories

import (
	"fmt"
	"log"

	"github.com/NunChatSpace/7-solutions-backend-challenge/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	User() IUserRepository
}

type RepositoryImpl struct {
	userRepo IUserRepository
}

func NewRepository(cfg *config.Config) Repository {
	var (
		pgConfig = cfg.Database.Postgres
		db       *gorm.DB
	)
	log.SetFlags(0)
	conn := postgres.Config{
		DSN: "host=" + pgConfig.Host + " user=" + pgConfig.User + " password=" + pgConfig.Password + " dbname=" + pgConfig.Name + " port=" + pgConfig.Port + " sslmode=disable TimeZone=" + cfg.App.Timezone, //nolint:lll
	}
	fmt.Println("Connection DSN:", conn.DSN)
	tmpdb, err := gorm.Open(postgres.New(conn))
	if err != nil {
		log.Println(err)
		log.Fatal("Cannot connect to database: Please check the connection")
	}
	db = tmpdb

	return &RepositoryImpl{
		userRepo: NewUserRepository(db),
	}
}
