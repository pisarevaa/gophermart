package main

import (
	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/storage"
)

func main() {
	config := server.NewConfig()
	logger := server.NewLogger()
	repo := storage.NewDB(config.DatabaseUri, logger)
	r := server.NewRouter(config, logger, repo)
	logger.Fatal(r.Run(config.Host))
}
