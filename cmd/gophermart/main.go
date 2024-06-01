package main

import (
	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

func main() {
	cfg := configs.NewConfig()
	logger := server.NewLogger()
	repo := storage.NewDB(cfg.DatabaseUri, logger)
	r := server.NewRouter(cfg, logger, repo)
	logger.Fatal(r.Run(cfg.Host))
}
