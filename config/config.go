package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Database Database
}

type App struct {
	Port int
}
type Database struct {
	PostgreDSN string
}

var (
	once sync.Once
	cfg  *Config
)

// LoadConfig from .env once
func LoadConfig() *Config {
	once.Do(func() {
		// Load .env file only in local development
		if os.Getenv("APP_ENV") == "local" {
			if err := godotenv.Load(); err != nil {
				fmt.Println("Warning: No .env file found, using system environment variables")
			}
		}

		cfg = &Config{}

		port, err := strconv.Atoi(os.Getenv("APP_PORT"))
		if err != nil {
			panic(fmt.Errorf("invalid APP_PORT in .env: %w", err))
		}
		cfg.App.Port = port

		cfg.Database.PostgreDSN = os.Getenv("DATABASE_URL")
	})

	return cfg
}
