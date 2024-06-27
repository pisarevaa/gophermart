package storage

import (
	"context"
	"errors"
	"time"
)

type MemoryStorage struct {
	Users  map[string]User
	Orders map[string]Order
}

type MemoryTransaction struct {
	Transaction *MemoryStorage
}

func NewMemory() *MemoryStorage {
	return &MemoryStorage{
		Users:  make(map[string]User),
		Orders: make(map[string]Order),
	}
}

func (m *MemoryStorage) GetUser(_ context.Context, login string) (User, error) {
	user, ok := m.Users[login]
	if !ok {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (m *MemoryStorage) StoreUser(_ context.Context, login string, passwordHash string) error {
	m.Users[login] = User{
		Login:     login,
		Password:  passwordHash,
		Balance:   0,
		Withdrawn: 0,
	}
	return nil
}

func (m *MemoryStorage) GetOrder(_ context.Context, number string) (Order, error) {
	order, ok := m.Orders[number]
	if !ok {
		return order, errors.New("order not found")
	}
	return order, nil
}

func (m *MemoryStorage) GetOrders(_ context.Context, login string, onlyWithdrawn bool) ([]Order, error) {
	var orders []Order
	for _, order := range m.Orders {
		if onlyWithdrawn {
			if order.Withdrawn > 0 && order.Login == login {
				orders = append(orders, order)
			}
		} else if order.Login == login {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (m *MemoryStorage) GetOrdersCountToUpdate(_ context.Context) (int64, error) {
	var count int64
	for _, order := range m.Orders {
		if order.Status == "NEW" || order.Status == "PROCESSING" || order.Status == "REGISTERED" {
			count++
		}
	}
	return count, nil
}

func (m *MemoryStorage) StoreOrder(_ context.Context, number, login string) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	m.Orders[number] = Order{
		Number:     number,
		Status:     "NEW",
		Accrual:    0,
		Withdrawn:  0,
		Login:      login,
		UploadedAt: time.Now().In(loc),
	}
	return nil
}

func (m *MemoryStorage) CloseConnection() {}

func (m *MemoryStorage) BeginTransaction(_ context.Context) (Transaction, error) {
	return m, nil
	// return m, nil
}

func (m *MemoryStorage) Commit(_ context.Context) error {
	return nil
}

func (m *MemoryStorage) Rollback(_ context.Context) error {
	return nil
}

func (m *MemoryStorage) GetOrderToUpdateStatus(_ context.Context) (OrderToUpdate, error) {
	for _, order := range m.Orders {
		if order.Status == "NEW" || order.Status == "PROCESSING" || order.Status == "REGISTERED" {
			return OrderToUpdate{
				Number: order.Number,
				Login:  order.Login,
			}, nil
		}
	}
	return OrderToUpdate{}, nil
}

func (m *MemoryStorage) UpdateOrderStatus(_ context.Context, order OrderStatus) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	if currentOrder, ok := m.Orders[order.Number]; ok {
		currentOrder.Status = order.Status
		currentOrder.Accrual = order.Accrual
		now := time.Now().In(loc)
		currentOrder.ProcessedAt = &now
		m.Orders[order.Number] = currentOrder
		return nil
	}
	return errors.New("order not found")
}

func (m *MemoryStorage) AccrualUserBalance(_ context.Context, accraul float32, login string) error {
	if currentUser, ok := m.Users[login]; ok {
		currentUser.Balance += accraul
		m.Users[login] = currentUser
		return nil
	}
	return errors.New("user not found")
}

func (m *MemoryStorage) GetUserWithLock(_ context.Context, login string) (User, error) {
	user, ok := m.Users[login]
	if !ok {
		return user, errors.New("user not found")
	}
	return user, nil
}

func (m *MemoryStorage) GetOrderWithLock(_ context.Context, number string, _ string) (Order, error) {
	order, ok := m.Orders[number]
	if !ok {
		return order, errors.New("order not found")
	}
	return order, nil
}

func (m *MemoryStorage) WithdrawUserBalance(_ context.Context, login string, withdraw float32) error {
	if currentUser, ok := m.Users[login]; ok {
		currentUser.Balance -= withdraw
		currentUser.Withdrawn += withdraw
		m.Users[login] = currentUser
		return nil
	}
	return errors.New("user not found")
}

func (m *MemoryStorage) WithdrawOrderBalance(_ context.Context, number string, withdraw float32) error {
	if currentOrder, ok := m.Orders[number]; ok {
		currentOrder.Withdrawn = withdraw
		m.Orders[number] = currentOrder
		return nil
	}
	return errors.New("order not found")
}
