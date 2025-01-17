package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/fvbock/endless"

	server "github.com/pisarevaa/gophermart/internal"
	"github.com/pisarevaa/gophermart/internal/configs"
	"github.com/pisarevaa/gophermart/internal/storage"
	"github.com/pisarevaa/gophermart/internal/tasks"
	"github.com/pisarevaa/gophermart/internal/utils"
)

// @title		Swagger Gophermart Service API
// @version	1.0
// @host		localhost:8080

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exit, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfg := configs.NewConfig()
	logger := server.NewLogger()
	repo := storage.NewDB(cfg.DatabaseURI, logger)
	r := server.NewRouter(cfg, logger, repo)

	// Запускаем фоновую задачу по обновлению статусов заказов
	client := utils.NewClient()
	task := tasks.NewTask(cfg, logger, repo, client)
	go task.RunUpdateOrderStatuses(exit)

	logger.Info("Run Server")
	logger.Fatal(endless.ListenAndServe(cfg.Host, r))

	<-exit.Done()
	logger.Info("Server Shutdown!")
}
