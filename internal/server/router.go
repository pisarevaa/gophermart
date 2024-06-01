package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

func NewRouter(config Config, logger *zap.SugaredLogger, repo *storage.DBStorage) *gin.Engine {
	server := Server{Config: config, Logger: logger, Repo: repo}
	r := gin.Default()
	r.GET("/ping", server.Hello)
	return r
}
