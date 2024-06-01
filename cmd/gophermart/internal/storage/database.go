package storage

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBStorage struct {
	*pgxpool.Pool
}

func NewDBStorage(dsn string, logger *zap.SugaredLogger) *DBStorage {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
	}
	err = m.Up()
	if err != nil {
		logger.Error("Unable to migrate tables: ", err)
	}
	db := &DBStorage{dbpool}
	return db
}

func (dbpool *DBStorage) CloseConnection() {
	dbpool.Close()
}
