package tasks

import (
	"context"
	"time"

	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

type Task struct {
	Config configs.Config
	Logger *zap.SugaredLogger
	Repo   storage.Storage
}

func NewTask(
	config configs.Config,
	logger *zap.SugaredLogger,
	repo storage.Storage,
) *Task {
	return &Task{
		Config: config,
		Logger: logger,
		Repo:   repo,
	}
}

func (s *Task) RunUpdateOrderStatuses(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.Config.TaskInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := s.UpdateOrderStatuses()
			if err != nil {
				s.Logger.Error("error to update order statuses:", err)
			}
		case <-ctx.Done():
			s.Logger.Info("ctx.Done -> exit RunUpdateOrderStatuses")
			return
		}
	}
}

func (s *Task) UpdateOrderStatuses() error {
	s.Logger.Info("UpdateOrderStatuses....")
	return nil
}
