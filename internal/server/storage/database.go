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

func NewDB(databaseUri string, logger *zap.SugaredLogger) *DBStorage {
	dbpool, err := pgxpool.New(context.Background(), databaseUri)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	m, err := migrate.New("file://migrations", databaseUri)
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

func (dbpool *DBStorage) GetUser(ctx context.Context, login string) (User, error) {
	var user User
	err := dbpool.QueryRow(ctx, "SELECT login, password FROM users WHERE login = $1", login).
		Scan(&user.Login, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dbpool *DBStorage) StoreUser(ctx context.Context, login string, passwordHash string) error {
	_, err := dbpool.Exec(ctx, `
			INSERT INTO users (login, password) VALUES ($1, $2)
		`, login, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBStorage) CloseConnection() {
	dbpool.Close()
}
