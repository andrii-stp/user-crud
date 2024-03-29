package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   *Server
	Database *Database
}

type Server struct {
	Port string
}

type Database struct {
	Driver   string
	Host     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load(envFile string) (*Config, error) {
	if _, err := os.Stat(envFile); err == nil {
		_ = godotenv.Load(envFile)
	}

	if err := getMissingEnvs(); err != nil {
		return nil, err
	}

	return &Config{
		Server: &Server{
			Port: os.Getenv("SERVER_PORT"),
		},
		Database: &Database{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}, nil
}

func getMissingEnvs() error {
	var (
		missingEnvs []string
		envs        = []string{"SERVER_PORT", "DB_DRIVER", "DB_HOST", "DB_USERNAME", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"}
	)

	for _, env := range envs {
		if val := os.Getenv(env); len(val) == 0 {
			missingEnvs = append(missingEnvs, env)
		}
	}

	if len(missingEnvs) > 0 {
		return fmt.Errorf("missing env variables: %v", missingEnvs)
	}

	return nil
}
