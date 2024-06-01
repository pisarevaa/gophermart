package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

func NewRouter(cfg configs.Config, logger *zap.SugaredLogger, repo storage.Storage) *gin.Engine {
	server := handlers.Server{Config: cfg, Logger: logger, Repo: repo}
	if cfg.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	basePath := r.Group("/api/user")
	{
		basePath.POST("/register", server.RegisterUser)
		basePath.POST("/login", server.LoginUser)
	}
	return r
}
