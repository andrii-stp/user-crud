// @title User API
// @version 0.1
// @description This is a User server.

// @host localhost:8080
// @BasePath /api/v1
// @schemes http
package main

import (
	"log/slog"
	"os"

	"github.com/andrii-stp/users-crud/config"
	"github.com/andrii-stp/users-crud/router"
	"github.com/andrii-stp/users-crud/storage"

	_ "github.com/swaggo/echo-swagger/example/docs"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	configFile := ".env"

	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Error("failed to load config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	db, err := storage.Connect(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if err = storage.InitDB(db); err != nil {
		logger.Error("failed to initialize database schema", slog.String("err", err.Error()))
		os.Exit(1)
	}

	repo := storage.NewPostgresRepository(logger, db)

	server := router.Router(logger, repo)
	port := ":" + cfg.Server.Port

	if err = server.Start(port); err != nil {
		server.Logger.Fatal("error when server is initializing")
	}
}
