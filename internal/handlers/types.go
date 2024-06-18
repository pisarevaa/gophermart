package handlers

import (
	"github.com/pisarevaa/gophermart/internal/configs"
	"github.com/pisarevaa/gophermart/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	Config configs.Config
	Logger *zap.SugaredLogger
	Repo   storage.Storage
}

func NewController(
	config configs.Config,
	logger *zap.SugaredLogger,
	repo storage.Storage,
) *Service {
	return &Service{
		Config: config,
		Logger: logger,
		Repo:   repo,
	}
}
