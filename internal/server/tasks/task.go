package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/pisarevaa/gophermart/internal/server/configs"
	"github.com/pisarevaa/gophermart/internal/server/storage"
	"go.uber.org/zap"
)

type Task struct {
	Config configs.Config
	Logger *zap.SugaredLogger
	Repo   storage.Storage
	Client *resty.Client
}

func NewTask(
	config configs.Config,
	logger *zap.SugaredLogger,
	repo storage.Storage,
	client *resty.Client,
) *Task {
	return &Task{
		Config: config,
		Logger: logger,
		Repo:   repo,
		Client: client,
	}
}

func (s *Task) RunUpdateOrderStatuses(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.Config.TaskInterval) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := s.UpdateOrderStatuses(ctx)
			if err != nil {
				s.Logger.Error("error to update order statuses:", err)
			}
		case <-ctx.Done():
			s.Logger.Info("ctx.Done -> exit RunUpdateOrderStatuses")
			return
		}
	}
}

func (s *Task) UpdateOrderStatuses(ctx context.Context) error {
	s.Logger.Info("UpdateOrderStatuses....")

	count, err := s.Repo.GetOrdersCountToUpdate(ctx)
	if err != nil {
		return err
	}
	if count == 0 {
		s.Logger.Info("Not active orders to proccess")
		return nil
	}
	tx, err := s.Repo.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // ignore check
	orderToUpdate, err := tx.GetOrderToUpdateStatus(ctx)
	if err != nil {
		return err
	}
	status, err := s.GetOrderStatus(orderToUpdate.Number)
	if err != nil {
		return err
	}
	if status.Status == "NEW" || status.Status == "PROCESSING" || status.Status == "REGISTERED" {
		s.Logger.Info("order is not ready")
		return nil
	}
	err = tx.UpdateOrderStatus(ctx, status)
	if err != nil {
		return err
	}
	err = tx.AccrualUserBalance(ctx, status.Accrual, orderToUpdate.Login)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	s.Logger.Info("order is updated successfully ", orderToUpdate.Number)
	return nil
}

func (s *Task) GetOrderStatus(number string) (storage.OrderStatus, error) {
	var orderStatus storage.OrderStatus
	// requestURL := fmt.Sprintf("http://%v/api/orders/%v", s.Config.AccrualSystemAddress, number)
	requestURL := fmt.Sprintf("%v/api/orders/%v", s.Config.AccrualSystemAddress, number)
	resp, err := s.Client.R().SetResult(&orderStatus).SetHeader("Content-Type", "application/json").Get(requestURL)
	if err != nil {
		s.Logger.Info("Request to ", requestURL, " with Error: ", err)
		return orderStatus, err
	}
	s.Logger.Info("Request to ", requestURL, " with status code:  ", resp.StatusCode())
	s.Logger.Info("Request to ", requestURL, " with resp.RawResponse:  ", resp.RawResponse)
	return orderStatus, nil
}
