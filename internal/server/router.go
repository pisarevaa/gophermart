package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/pisarevaa/gophermart/docs"
	"go.uber.org/zap"
)

func NewRouter(cfg configs.Config, logger *zap.SugaredLogger, repo storage.Storage) *gin.Engine {
	s := handlers.NewController(cfg, logger, repo)
	if cfg.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	basePath := r.Group("/api/user")
	{
		basePath.POST("/register", s.RegisterUser)
		basePath.POST("/login", s.LoginUser)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
