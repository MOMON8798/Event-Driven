package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(dsn string) (Repository, error) {
	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConns = 25
	poolCfg.MinConns = 5
	poolCfg.MaxConnLifetime = 5 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &postgresRepository{pool: pool}, nil
}

func (r *postgresRepository) Close() {
	r.pool.Close()
}

func (r *postgresRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	const q = `SELECT id, order_id, amount, status, created_at FROM payments WHERE id = $1`

	rows, err := r.pool.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}

	payment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[domain.Payment])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *postgresRepository) GetPaymentsByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	const q = `SELECT id, order_id, amount, status, created_at FROM payments WHERE order_id = $1 ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q, orderID)
	if err != nil {
		return nil, err
	}

	payments, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[domain.Payment])
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *postgresRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	const q = `INSERT INTO payments (id, order_id, amount, status, created_at)
               VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, q, payment.ID, payment.OrderID, payment.Amount, payment.Status, payment.CreatedAt)
	return err
}

func (r *postgresRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	const q = `UPDATE payments SET status = $1 WHERE id = $2`

	tag, err := r.pool.Exec(ctx, q, payment.Status, payment.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}
