package storage

import (
	"context"
	"time"
)

type Storage interface {
	GetUser(ctx context.Context, login string) (user User, err error)
	StoreUser(ctx context.Context, login string, passwordHash string) (err error)
	GetOrder(ctx context.Context, number string) (order Order, err error)
	GetOrders(ctx context.Context, login string) (orders []Order, err error)
	StoreOrder(ctx context.Context, number, login string) (err error)
	CloseConnection()
}

type User struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Order struct {
	Number     string    `json:"number"     binding:"required"`
	Status     string    `json:"status"     binding:"required"`
	Accrual    int64     `json:"accrual"    binding:"required"`
	Login      string    `json:"login"      binding:"required"`
	UploadedAt time.Time `json:"uploadedAt" binding:"required"`
}
