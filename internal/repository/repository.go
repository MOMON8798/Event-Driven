package repository

import (
	"sync"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
)

type Repository interface {
	GetOrderByID(id string) (*domain.Order, error)
	CreateOrder(order *domain.Order) error
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id string) error
	GetAllOrders() ([]*domain.Order, error)
}

type inMemoryRepository struct {
	orders map[string]*domain.Order
	mu     sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		orders: make(map[string]*domain.Order),
		mu:     sync.RWMutex{},
	}
}

func (r *inMemoryRepository) GetOrderByID(id string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (r *inMemoryRepository) CreateOrder(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepository) UpdateOrder(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.ID]; !exists {
		return domain.ErrOrderNotFound
	}
	r.orders[order.ID] = order
	return nil
}

func (r *inMemoryRepository) DeleteOrder(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[id]; !exists {
		return domain.ErrOrderNotFound
	}
	delete(r.orders, id)
	return nil
}

func (r *inMemoryRepository) GetAllOrders() ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	orders := make([]*domain.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}
