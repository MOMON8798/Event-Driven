package service_test

import (
	"context"
	"testing"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestService() service.OrderService {
	repo := repository.NewInMemoryRepository()
	return service.NewOrderService(repo)
}

func TestCreateOrder_SetsDefaults(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, "client-1", "Laptop", 999.99)
	require.NoError(t, err)

	assert.NotEmpty(t, order.ID, "CreateOrder should generate a UUID itself")
	assert.Equal(t, "client-1", order.ClientID)
	assert.Equal(t, "Laptop", order.Name)
	assert.Equal(t, 999.99, order.Total)
	assert.Equal(t, domain.StatusCreated, order.Status, "a new order must get the created status, not any other")
	assert.False(t, order.CreatedAt.IsZero(), "CreatedAt should be set by the server")
}

func TestCreateOrder_GeneratesUniqueIDs(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	first, err := svc.CreateOrder(ctx, "client-1", "Item A", 10)
	require.NoError(t, err)

	second, err := svc.CreateOrder(ctx, "client-1", "Item B", 20)
	require.NoError(t, err)

	assert.NotEqual(t, first.ID, second.ID, "two different orders must not get the same ID")
}

func TestGetOrderByID_Found(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	created, err := svc.CreateOrder(ctx, "client-1", "Item", 100)
	require.NoError(t, err)

	fetched, err := svc.GetOrderByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)
	assert.Equal(t, created.Name, fetched.Name)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.GetOrderByID(ctx, "non-existent-id")

	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestUpdateOrder_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, "client-1", "Item", 100)
	require.NoError(t, err)

	order.Name = "Updated name"
	order.Total = 200
	order.Status = domain.StatusPaid

	err = svc.UpdateOrder(ctx, order)
	require.NoError(t, err)

	fetched, err := svc.GetOrderByID(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated name", fetched.Name)
	assert.Equal(t, 200.0, fetched.Total)
	assert.Equal(t, domain.StatusPaid, fetched.Status)
}

func TestUpdateOrder_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	fakeOrder := &domain.Order{ID: "non-existent-id", Name: "x", Total: 1, Status: domain.StatusCreated}

	err := svc.UpdateOrder(ctx, fakeOrder)

	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestDeleteOrder_Success(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	order, err := svc.CreateOrder(ctx, "client-1", "Item", 100)
	require.NoError(t, err)

	err = svc.DeleteOrder(ctx, order.ID)
	require.NoError(t, err)

	_, err = svc.GetOrderByID(ctx, order.ID)
	assert.ErrorIs(t, err, domain.ErrOrderNotFound, "the order should not be found after deletion")
}

func TestDeleteOrder_NotFound(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	err := svc.DeleteOrder(ctx, "non-existent-id")

	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestGetAllOrders_ReturnsAllCreated(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	_, err := svc.CreateOrder(ctx, "client-1", "Item A", 10)
	require.NoError(t, err)
	_, err = svc.CreateOrder(ctx, "client-2", "Item B", 20)
	require.NoError(t, err)

	all, err := svc.GetAllOrders(ctx)
	require.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestGetAllOrders_EmptyByDefault(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	all, err := svc.GetAllOrders(ctx)
	require.NoError(t, err)
	assert.Empty(t, all, "a fresh repository should return an empty list, not a nil error")
}
