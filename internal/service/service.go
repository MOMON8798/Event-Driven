package service

import (
	"time"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/google/uuid"
)

type OrderService interface {
	GetOrderByID(id string) (*domain.Order, error)
	CreateOrder(clientID string, name string, total float64) (*domain.Order, error)
	UpdateOrder(order *domain.Order) error
	DeleteOrder(id string) error
	GetAllOrders() ([]*domain.Order, error)
}

type orderService struct {
	repo repository.Repository
}

func NewOrderService(repo repository.Repository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetOrderByID(id string) (*domain.Order, error) {
	return s.repo.GetOrderByID(id)
}

func (s *orderService) CreateOrder(clientID string, name string, total float64) (*domain.Order, error) {
	order := &domain.Order{
		ID:        uuid.New().String(),
		ClientID:  clientID,
		Name:      name,
		Total:     total,
		Status:    domain.StatusCreated,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateOrder(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) UpdateOrder(order *domain.Order) error {
	return s.repo.UpdateOrder(order)
}

func (s *orderService) DeleteOrder(id string) error {
	return s.repo.DeleteOrder(id)
}

func (s *orderService) GetAllOrders() ([]*domain.Order, error) {
	return s.repo.GetAllOrders()
}
