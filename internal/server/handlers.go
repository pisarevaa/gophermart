package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

type Server struct {
	Config Config
	Logger *zap.SugaredLogger
	Repo   *storage.DBStorage
}

func (s *Server) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
