package storage

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DBStorage struct {
	*pgxpool.Pool
}

type DBTransaction struct {
	pgx.Tx
}

func NewDB(databaseURI string, logger *zap.SugaredLogger) *DBStorage {
	dbpool, err := pgxpool.New(context.Background(), databaseURI)
	if err != nil {
		logger.Error("Unable to create connection pool: %v", err)
		return nil
	}
	m, err := migrate.New("file://migrations", databaseURI)
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
	err := dbpool.QueryRow(ctx, "SELECT login, password, balance, withdrawn FROM users WHERE login = $1", login).
		Scan(&user.Login, &user.Password, &user.Balance, &user.Withdrawn)
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

func (dbpool *DBStorage) GetOrder(ctx context.Context, number string) (Order, error) {
	var order Order
	err := dbpool.QueryRow(ctx, "SELECT number, status, accrual, withdrawn, login, uploaded_at, processed_at FROM orders WHERE number = $1", number).
		Scan(&order.Number, &order.Status, &order.Accrual, &order.Withdrawn, &order.Login, &order.UploadedAt, &order.ProcessedAt)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (dbpool *DBStorage) GetOrders(ctx context.Context, login string, onlyWithdrawn bool) ([]Order, error) {
	sql := "SELECT number, status, accrual, withdrawn, login, uploaded_at, processed_at FROM orders WHERE login = $1 "
	if onlyWithdrawn {
		sql += " AND withdrawn  >  0"
	}
	sql += " ORDER BY uploaded_at ASC"
	var orders []Order
	rows, err := dbpool.Query(
		ctx, sql, login,
	)
	if err != nil {
		return []Order{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var o Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.Withdrawn, &o.Login, &o.UploadedAt, &o.ProcessedAt)
		if err != nil {
			return []Order{}, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (dbpool *DBStorage) GetOrdersCountToUpdate(ctx context.Context) (int64, error) {
	var count int64
	err := dbpool.QueryRow(ctx, "SELECT COUNT(*) AS count FROM orders WHERE status = 'NEW' OR status = 'PROCESSING' OR status = 'REGISTERED'").
		Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
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

func (dbpool *DBStorage) BeginTransaction(ctx context.Context) (*DBTransaction, error) {
	tx, err := dbpool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	dbTx := &DBTransaction{tx}
	return dbTx, nil
}

func (tx *DBTransaction) GetOrderToUpdateStatus(ctx context.Context) (OrderToUpdate, error) {
	var order OrderToUpdate
	err := tx.QueryRow(ctx, "SELECT number, login FROM orders WHERE status = 'NEW' OR status = 'PROCESSING' OR status = 'REGISTERED' LIMIT 1 FOR UPDATE SKIP LOCKED").
		Scan(&order.Number, &order.Login)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (tx *DBTransaction) UpdateOrderStatus(ctx context.Context, order OrderStatus) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `
			UPDATE orders SET status = $1, accrual = $2, processed_at = $3 WHERE number = $4
		`, order.Status, order.Accrual, time.Now().In(loc), order.Number)
	if err != nil {
		return err
	}
	return nil
}

func (tx *DBTransaction) AccrualUserBalance(ctx context.Context, accraul float32, login string) error {
	_, err := tx.Exec(ctx, `
			UPDATE users SET balance = balance + $1 WHERE login = $2
		`, accraul, login)
	if err != nil {
		return err
	}
	return nil
}

func (tx *DBTransaction) GetUserWithLock(ctx context.Context, login string) (User, error) {
	var user User
	err := tx.QueryRow(ctx, "SELECT login, password, balance FROM users WHERE login = $1 FOR UPDATE", login).
		Scan(&user.Login, &user.Password, &user.Balance)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (tx *DBTransaction) GetOrderWithLock(ctx context.Context, number string) (Order, error) {
	var order Order
	err := tx.QueryRow(ctx, "SELECT number, status, accrual, withdrawn, login, uploaded_at, processed_at FROM orders WHERE number = $1", number).
		Scan(&order.Number, &order.Status, &order.Accrual, &order.Withdrawn, &order.Login, &order.UploadedAt, &order.ProcessedAt)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (tx *DBTransaction) WithdrawUserBalance(ctx context.Context, login string, withdraw float32) error {
	_, err := tx.Exec(ctx, `
			UPDATE users SET balance = balance - $1, withdrawn = withdrawn + $2 WHERE login = $3
		`, withdraw, withdraw, login)
	if err != nil {
		return err
	}
	return nil
}

func (tx *DBTransaction) WithdrawOrderBalance(ctx context.Context, number string, withdraw float32) error {
	_, err := tx.Exec(ctx, `
			UPDATE orders SET withdrawn = withdrawn + $1 WHERE number = $2
		`, withdraw, number)
	if err != nil {
		return err
	}
	return nil
}
