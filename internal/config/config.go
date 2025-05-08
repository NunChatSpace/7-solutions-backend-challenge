package config

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
	"gorm.io/gorm/logger"
)

type Config struct {
	App struct {
		URL       string `default:"http://localhost:8888" env:"APP__URL"`
		Port      string `default:":8888" env:"APP__PORT"`
		DebugMode bool   `default:"false" env:"APP__DEBUG_MODE"`
		// Timezone refer in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
		Timezone string `default:"Asia/Manila" env:"APP__TIME_ZONE"`
		Version  string `default:"" env:"APP__VERSION"`
		Env      string `default:"dev" env:"APP__ENV"`
	}

	Cors struct {
		// Domains need single quotes if you're using asterisk: https://github.com/jinzhu/configor/issues/55
		Domains string `default:"'*'" env:"CORS__DOMAINS"`
	}

	Database struct {
		MongoDB struct {
			Username     string `default:"mongou" env:"DATABASE__MONGODB__USERNAME"`
			Password     string `default:"mongop" env:"DATABASE__MONGODB__PASSWORD"`
			Host         string `default:"localhost" env:"DATABASE__MONGODB__HOST"`
			Port         string `default:"27017" env:"DATABASE__MONGODB__PORT"`
			DatabaseName string `default:"challengeApp" env:"DATABASE__MONGODB__NAME"`
		}
		LogLevel logger.LogLevel `default:"4" env:"DATABASE__LOG_LEVEL"`
	}

	Log struct {
		Level  string `env:"LOG__LEVEL" default:"trace"`
		Format string `env:"LOG__FORMAT" default:"console"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config
	var err error

	err = loadEnv()
	if err == nil {
		log.Println("Env file loaded")
	}

	err = configor.
		New(&configor.Config{AutoReload: false}).
		Load(&config)

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading config")
	}

	fmt.Println("config.App.Env", config.App.Env)
	if config.App.Env == "test" {
		config.Database.MongoDB.DatabaseName = "_test"
	}

	return &config, err
}
