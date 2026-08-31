package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
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

func (r *postgresRepository) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	const q = `SELECT id, client_id, name, total, status, created_at FROM orders WHERE id = $1`

	rows, err := r.pool.Query(ctx, q, id)
	if err != nil {
		return nil, err
	}

	order, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByPos[domain.Order])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *postgresRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	const q = `INSERT INTO orders (id, client_id, name, total, status, created_at)
               VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, q, order.ID, order.ClientID, order.Name, order.Total, order.Status, order.CreatedAt)
	return err
}

func (r *postgresRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	const q = `UPDATE orders SET name = $1, total = $2, status = $3 WHERE id = $4`

	tag, err := r.pool.Exec(ctx, q, order.Name, order.Total, order.Status, order.ID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *postgresRepository) DeleteOrder(ctx context.Context, id string) error {
	const q = `DELETE FROM orders WHERE id = $1`

	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrOrderNotFound
	}
	return nil
}

func (r *postgresRepository) GetAllOrders(ctx context.Context) ([]*domain.Order, error) {
	const q = `SELECT id, client_id, name, total, status, created_at FROM orders ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByPos[domain.Order])
	if err != nil {
		return nil, err
	}
	return orders, nil
}
