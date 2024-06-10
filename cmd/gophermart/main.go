package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/pisarevaa/gophermart/internal/server"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/tasks"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

// @title           Swagger Example API
// @version         2.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exit, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfg := configs.NewConfig()
	logger := server.NewLogger()
	repo := storage.NewDB(cfg.DatabaseUri, logger)
	r := server.NewRouter(cfg, logger, repo)
	logger.Info("Run Server")
	srv := &http.Server{
		Addr:    cfg.Host,
		Handler: r.Handler(),
	}

	// Поднимает http сервер
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// Запускаем фоновую задачу по обновлению статусов заказов
	client := utils.NewClient()
	task := tasks.NewTask(cfg, logger, repo, client)
	go task.RunUpdateOrderStatuses(exit)

	if err := srv.Shutdown(exit); err != nil {
		logger.Fatal("Server Shutdown:", err)
	}

	<-exit.Done()
	logger.Info("Server Shutdown!")
}
