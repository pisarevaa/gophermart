package storage

import (
	"context"
	"time"

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
	err := dbpool.QueryRow(ctx, "SELECT login, password, balance FROM users WHERE login = $1", login).
		Scan(&user.Login, &user.Password, &user.Balance)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (dbpool *DBStorage) StoreUser(ctx context.Context, login string, passwordHash string) error {
	_, err := dbpool.Exec(ctx, `
			INSERT INTO users (login, password, balance) VALUES ($1, $2, $3)
		`, login, passwordHash, 0)
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBStorage) WithdrawUserBalance(ctx context.Context, login string, withdraw int64) error {
	_, err := dbpool.Exec(ctx, `
			UPDATE users SET balance = balance - $1 WHERE login = $2
		`, withdraw, login)
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBStorage) AccrualOrderBalance(ctx context.Context, number string, withdraw int64) error {
	_, err := dbpool.Exec(ctx, `
			UPDATE orders SET accrual = accrual + $1 WHERE number = $2
		`, withdraw, number)
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBStorage) GetOrder(ctx context.Context, number string) (Order, error) {
	var order Order
	err := dbpool.QueryRow(ctx, "SELECT number, status, accrual, login, uploaded_at FROM orders WHERE number = $1", number).
		Scan(&order.Number, &order.Status, &order.Accrual, &order.Login, &order.UploadedAt)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (dbpool *DBStorage) GetOrders(ctx context.Context, login string, onlyAccrual bool) ([]Order, error) {
	sql := "SELECT number, status, accrual, login, uploaded_at FROM orders WHERE login = $1 "
	if onlyAccrual {
		sql += " AND accrual  >  0" //nolint:ineffassign
	}
	var orders []Order
	rows, err := dbpool.Query(
		ctx,
		"SELECT number, status, accrual, login, uploaded_at FROM orders WHERE login = $1",
		login,
	)
	if err != nil {
		return []Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.Login, &o.UploadedAt)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (dbpool *DBStorage) GetUserWithdrawals(ctx context.Context, login string) (int64, error) {
	var sum int64
	err := dbpool.QueryRow(ctx, "SELECT SUM(accrual) FROM orders WHERE login = $1", login).
		Scan(&sum)
	if err != nil {
		return sum, err
	}
	return sum, nil
}

func (dbpool *DBStorage) StoreOrder(ctx context.Context, number, login string) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	_, err = dbpool.Exec(ctx, `
			INSERT INTO orders (number, status, accrual, login, uploaded_at) VALUES ($1, $2, $3, $4, $5)
		`, number, "NEW", 0, login, time.Now().In(loc))
	if err != nil {
		return err
	}
	return nil
}

func (dbpool *DBStorage) CloseConnection() {
	dbpool.Close()
}
