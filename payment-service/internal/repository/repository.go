package repository

import (
	"context"
	"sync"

	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/domain"
)

type Repository interface {
	GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error)
	GetPaymentsByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error)
	CreatePayment(ctx context.Context, payment *domain.Payment) error
	UpdatePayment(ctx context.Context, payment *domain.Payment) error
}

type inMemoryRepository struct {
	payments map[string]*domain.Payment
	mu       sync.RWMutex
}

func NewInMemoryRepository() Repository {
	return &inMemoryRepository{
		payments: make(map[string]*domain.Payment),
	}
}

func (r *inMemoryRepository) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, exists := r.payments[id]
	if !exists {
		return nil, domain.ErrPaymentNotFound
	}
	return p, nil
}

func (r *inMemoryRepository) GetPaymentsByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Payment, 0)
	for _, p := range r.payments {
		if p.OrderID == orderID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *inMemoryRepository) CreatePayment(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[payment.ID] = payment
	return nil
}

func (r *inMemoryRepository) UpdatePayment(ctx context.Context, payment *domain.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.payments[payment.ID]; !exists {
		return domain.ErrPaymentNotFound
	}
	r.payments[payment.ID] = payment
	return nil
}
