package config

import (
	"log"

	"github.com/jinzhu/configor"
	"gorm.io/gorm/logger"
)

type Config struct {
	App struct {
		URL       string `default:"http://localhost:8080" env:"APP__URL"`
		Port      string `default:":8080" env:"APP__PORT"`
		DebugMode bool   `default:"false" env:"APP__DEBUG_MODE"`
		// Timezone refer in https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
		Timezone          string `default:"Asia/Manila" env:"APP__TIME_ZONE"`
		Version           string `default:"" env:"APP__VERSION"`
		Currency          string `default:"PHP" env:"APP__CURRENCY"`
		Env               string `default:"dev" env:"APP__ENV"`
		EnabledMonitoring bool   `default:"true" env:"APP__ENABLE_MONITORING"`
	}

	Cors struct {
		// Domains need single quotes if you're using asterisk: https://github.com/jinzhu/configor/issues/55
		Domains string `default:"'*'" env:"CORS__DOMAINS"`
	}

	Database struct {
		Postgres struct {
			Host     string `default:"localhost" env:"DATABASE__POSTGRES__HOST"`
			User     string `default:"postgres" env:"DATABASE__POSTGRES__USER"`
			Password string `default:"postgres" env:"DATABASE__POSTGRES__PASSWORD"`
			Port     string `default:"54321" env:"DATABASE__POSTGRES__PORT"`
			Name     string `default:"notify-chat-local" env:"DATABASE__POSTGRES__NAME"`
			SSLMode  string `default:"disable" env:"DATABASE__POSTGRES__SSL_MODE"`
		}
		LogLevel logger.LogLevel `default:"4" env:"DATABASE__LOG_LEVEL"`
	}

	Log struct {
		Level  string `env:"LOG__LEVEL" default:"trace"`
		Format string `env:"LOG__FORMAT" default:"console"`
	}

	Smtp struct {
		Email    string `env:"SMTP__EMAIL" default:""`
		Password string `env:"SMTP__PASSWORD" default:""`
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

	return &config, err
}
