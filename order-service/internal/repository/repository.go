package repository

import (
	"context"
	"sync"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
)

type Repository interface {
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	CreateOrder(ctx context.Context, order *domain.Order) error
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetAllOrders(ctx context.Context) ([]*domain.Order, error)
}

type inMemoryRepository struct {
	orders map[string]*domain.Order
	mu     sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *inMemoryRepository) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (r *inMemoryRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.ID]; !exists {
		return domain.ErrOrderNotFound
	}
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepository) DeleteOrder(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[id]; !exists {
		return domain.ErrOrderNotFound
	}
	delete(r.orders, id)
	return nil
}

func (r *inMemoryRepository) GetAllOrders(ctx context.Context) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]*domain.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}
