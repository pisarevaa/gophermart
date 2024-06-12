package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetUser(ctx context.Context, login string) (user User, err error)
	StoreUser(ctx context.Context, login string, passwordHash string) (err error)
	GetOrder(ctx context.Context, number string) (order Order, err error)
	GetOrders(ctx context.Context, login string, onlyWithdrawn bool) (orders []Order, err error)
	GetOrdersCountToUpdate(ctx context.Context) (count int64, err error)
	StoreOrder(ctx context.Context, number, login string) (err error)
	BeginTransaction(ctx context.Context) (tx *DBTransaction, err error)
	CloseConnection()
}

type Transaction interface {
	GetOrderToUpdateStatus(ctx context.Context) (orderToUpdate OrderToUpdate, err error)
	UpdateOrderStatus(ctx context.Context, order OrderStatus) (err error)
	AccrualUserBalance(ctx context.Context, accraul int64, login string) (err error)
	GetUserWithLock(ctx context.Context, login string) (user User, err error)
	GetOrderWithLock(ctx context.Context, number string) (order Order, err error)
	WithdrawUserBalance(ctx context.Context, login string, withdraw int64) (err error)
	WithdrawOrderBalance(ctx context.Context, number string, withdraw int64) (err error)
	Commit(ctx context.Context) (err error)
}

type RegisterUser struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Login     string `json:"login"     binding:"required"`
	Password  string `json:"password"  binding:"required"`
	Balance   int64  `json:"balance"   binding:"required"`
	Withdrawn int64  `json:"withdrawn" binding:"required"`
}

type Order struct {
	Number      string    `json:"number"      binding:"required"`
	Status      string    `json:"status"      binding:"required"`
	Accrual     int64     `json:"accrual"     binding:"required"`
	Withdrawn   int64     `json:"withdrawn"   binding:"required"`
	Login       string    `json:"login"       binding:"required"`
	UploadedAt  time.Time `json:"uploadedAt"  binding:"required"`
	ProcessedAt time.Time `json:"processedAt" binding:"required"`
}

type OrderToUpdate struct {
	Number string `json:"number" binding:"required"`
	Login  string `json:"login"  binding:"required"`
}

type OrderStatus struct {
	Number  string `json:"number"  binding:"required"`
	Status  string `json:"status"  binding:"required"`
	Accrual int64  `json:"accrual" binding:"required"`
}
