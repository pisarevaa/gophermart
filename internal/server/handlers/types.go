package handlers

import (
	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

type Server struct {
	Config configs.Config
	Logger *zap.SugaredLogger
	Repo   *storage.DBStorage
}
