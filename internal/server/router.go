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

	api := r.Group("/api")
	{
		users := api.Group("/user")
		{
			users.POST("/register", s.RegisterUser)
			users.POST("/login", s.LoginUser)
			authorized := users.Group("/")
			authorized.Use(utils.JWTAuth(cfg.SecretKey))
			{
				authorized.POST("/orders", s.AddOrder)
				authorized.GET("/orders", s.GetOrders)
				authorized.GET("/balance", s.GetBalance)
				authorized.POST("/balance/withdraw", s.WithdrawBalance)
				authorized.GET("/withdrawals", s.Withdrawls)
			}
		}
		authorized := api.Group("/")
		authorized.Use(utils.JWTAuth(cfg.SecretKey))
		{
			authorized.POST("/orders/:number", s.GetOrder)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
