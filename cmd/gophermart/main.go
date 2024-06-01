package main

import (
	// "net/http"
	// "time"
	// "github.com/gin-gonic/gin"

	"github.com/pisarevaa/gophermart/internal/server"
	// "github.com/pisarevaa/metrics/internal/server/storage"
)

const readTimeout = 5
const writeTimout = 10

func main() {
	config := server.NewConfig()
	logger := server.NewLogger()
	// repo = storage.NewDBStorage(config.DatabaseDSN, logger)
	// defer repo.CloseConnection()
	// logger.Info("Server is running on ", config.Host)
	// srv := &http.Server{
	// 	Addr:         config.Host,
	// 	Handler:      server.MetricsRouter(config, logger, repo),
	// 	ReadTimeout:  readTimeout * time.Second,
	// 	WriteTimeout: writeTimout * time.Second,
	// }
	// logger.Fatal(srv.ListenAndServe())
}
