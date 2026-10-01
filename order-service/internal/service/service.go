package service

import (
	"context"
	"time"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, id string) (*domain.Order, error)
	CreateOrder(ctx context.Context, clientID string, name string, total float64) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id string) error
	GetAllOrders(ctx context.Context) ([]*domain.Order, error)
	PayOrder(ctx context.Context, id string) (*domain.Order, error)
}

type PaymentClient interface {
	CreatePayment(ctx context.Context, orderID string, amount float64) (*PaymentResult, error)
}

type orderService struct {
	repo          repository.Repository
	paymentClient PaymentClient
}

type PaymentResult struct {
	ID     string
	Status string
}

func NewOrderService(repo repository.Repository, paymentClient PaymentClient) OrderService {
	return &orderService{repo: repo, paymentClient: paymentClient}
}

func (s *orderService) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *orderService) CreateOrder(ctx context.Context, clientID string, name string, total float64) (*domain.Order, error) {
	order := &domain.Order{
		ID:        uuid.New().String(),
		ClientID:  clientID,
		Name:      name,
		Total:     total,
		Status:    domain.StatusCreated,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) UpdateOrder(ctx context.Context, order *domain.Order) error {
	return s.repo.UpdateOrder(ctx, order)
}

func (s *orderService) DeleteOrder(ctx context.Context, id string) error {
	return s.repo.DeleteOrder(ctx, id)
}

func (s *orderService) GetAllOrders(ctx context.Context) ([]*domain.Order, error) {
	return s.repo.GetAllOrders(ctx)
}

func (s *orderService) PayOrder(ctx context.Context, id string) (*domain.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if order.Status != domain.StatusCreated {
		return nil, domain.ErrOrderNotPayable
	}

	result, err := s.paymentClient.CreatePayment(ctx, order.ID, order.Total)
	if err != nil {
		return nil, err
	}

	if result.Status == "succeeded" {
		order.Status = domain.StatusPaid
		if err := s.repo.UpdateOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	return order, nil
}
