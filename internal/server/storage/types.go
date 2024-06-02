package storage

import (
	"context"
)

type Storage interface {
	GetUser(ctx context.Context, login string) (user User, err error)
	StoreUser(ctx context.Context, login string, passwordHash string) (err error)
	CloseConnection()
}

type User struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}
