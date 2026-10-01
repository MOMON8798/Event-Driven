package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/handler"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePaymentClient struct{}

func (f *fakePaymentClient) CreatePayment(ctx context.Context, orderID string, amount float64) (*service.PaymentResult, error) {
	return &service.PaymentResult{ID: "fake-payment-id", Status: "succeeded"}, nil
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	repo := repository.NewInMemoryRepository()
	svc := service.NewOrderService(repo, &fakePaymentClient{})
	h := handler.NewHandler(svc)

	router := gin.New()
	api := router.Group("/api")
	h.RegisterRoutes(api)
	return router
}

func doRequest(t *testing.T, router *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCreateOrder_Success(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Laptop",
		"total":     999.99,
	})

	assert.Equal(t, http.StatusCreated, rec.Code)

	var got domain.Order
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.NotEmpty(t, got.ID)
	assert.Equal(t, domain.StatusCreated, got.Status)
}

func TestCreateOrder_MissingRequiredField(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"name":  "Laptop",
		"total": 999.99,
	})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_ZeroTotalRejected(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Laptop",
		"total":     0,
	})

	assert.Equal(t, http.StatusBadRequest, rec.Code, "total=0 should be rejected the same way as in updateOrderRequest (gt=0)")
}

func TestGetOrderByID_Found(t *testing.T) {
	router := newTestRouter()

	created := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Item",
		"total":     50,
	})
	var createdOrder domain.Order
	require.NoError(t, json.Unmarshal(created.Body.Bytes(), &createdOrder))

	rec := doRequest(t, router, http.MethodGet, "/api/orders/"+createdOrder.ID, nil)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got domain.Order
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, createdOrder.ID, got.ID)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodGet, "/api/orders/does-not-exist", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateOrder_Success(t *testing.T) {
	router := newTestRouter()

	created := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Item",
		"total":     50,
	})
	var createdOrder domain.Order
	require.NoError(t, json.Unmarshal(created.Body.Bytes(), &createdOrder))

	rec := doRequest(t, router, http.MethodPut, "/api/orders/"+createdOrder.ID, map[string]any{
		"name":   "Updated item",
		"total":  75,
		"status": "paid",
	})

	assert.Equal(t, http.StatusOK, rec.Code)

	var updated domain.Order
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updated))
	assert.Equal(t, "Updated item", updated.Name)
	assert.Equal(t, domain.StatusPaid, updated.Status)
}

func TestUpdateOrder_NotFound(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodPut, "/api/orders/does-not-exist", map[string]any{
		"name":   "Item",
		"total":  10,
		"status": "created",
	})

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateOrder_NegativeTotalRejected(t *testing.T) {
	router := newTestRouter()

	created := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Item",
		"total":     50,
	})
	var createdOrder domain.Order
	require.NoError(t, json.Unmarshal(created.Body.Bytes(), &createdOrder))

	rec := doRequest(t, router, http.MethodPut, "/api/orders/"+createdOrder.ID, map[string]any{
		"name":   "Item",
		"total":  -10,
		"status": "created",
	})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteOrder_Success(t *testing.T) {
	router := newTestRouter()

	created := doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1",
		"name":      "Item",
		"total":     50,
	})
	var createdOrder domain.Order
	require.NoError(t, json.Unmarshal(created.Body.Bytes(), &createdOrder))

	rec := doRequest(t, router, http.MethodDelete, "/api/orders/"+createdOrder.ID, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	getRec := doRequest(t, router, http.MethodGet, "/api/orders/"+createdOrder.ID, nil)
	assert.Equal(t, http.StatusNotFound, getRec.Code, "the order should not be found after deletion")
}

func TestDeleteOrder_NotFound(t *testing.T) {
	router := newTestRouter()

	rec := doRequest(t, router, http.MethodDelete, "/api/orders/does-not-exist", nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetAllOrders_ReturnsCreatedOrders(t *testing.T) {
	router := newTestRouter()

	doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-1", "name": "Item A", "total": 10,
	})
	doRequest(t, router, http.MethodPost, "/api/orders/", map[string]any{
		"client_id": "client-2", "name": "Item B", "total": 20,
	})

	rec := doRequest(t, router, http.MethodGet, "/api/orders/", nil)

	assert.Equal(t, http.StatusOK, rec.Code)

	var orders []domain.Order
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &orders))
	assert.Len(t, orders, 2)
}
