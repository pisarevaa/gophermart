package server

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	docs "github.com/pisarevaa/gophermart/docs"
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/handlers"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func NewRouter(cfg configs.Config, logger *zap.SugaredLogger, repo storage.Storage) *gin.Engine {
	s := handlers.NewController(cfg, logger, repo)
	if cfg.GinMode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	docs.SwaggerInfo.BasePath = "/api/v1"
	basePath := r.Group("/api/user")
	{
		basePath.POST("/register", s.RegisterUser)
		basePath.POST("/login", s.LoginUser)
		authorized := basePath.Group("/")
		authorized.Use(utils.JWTAuth(cfg.SecretKey))
		{
			authorized.POST("/orders", s.AddOrder)
			authorized.GET("/orders", s.GetOrders)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
