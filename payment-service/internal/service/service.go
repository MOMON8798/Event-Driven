package service

import (
	"context"
	"time"

	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/repository"
	"github.com/google/uuid"
)

const maxAutoApprovedAmount = 10000.0

type PaymentService interface {
	GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error)
	GetPaymentsByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error)
	ProcessPayment(ctx context.Context, orderID string, amount float64) (*domain.Payment, error)
}

type paymentService struct {
	repo repository.Repository
}

func NewPaymentService(repo repository.Repository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) GetPaymentByID(ctx context.Context, id string) (*domain.Payment, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *paymentService) GetPaymentsByOrderID(ctx context.Context, orderID string) ([]*domain.Payment, error) {
	return s.repo.GetPaymentsByOrderID(ctx, orderID)
}

func (s *paymentService) ProcessPayment(ctx context.Context, orderID string, amount float64) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	payment := &domain.Payment{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		Amount:    amount,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	payment.Status = decideOutcome(amount)

	if err := s.repo.UpdatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func decideOutcome(amount float64) domain.Status {
	if amount > maxAutoApprovedAmount {
		return domain.StatusFailed
	}
	return domain.StatusSucceeded
}
