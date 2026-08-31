package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/gin-gonic/gin"
)

const requestTimeout = 3 * time.Second

type Handler struct {
	service service.OrderService
}

type createOrderRequest struct {
	ClientID string  `json:"client_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Total    float64 `json:"total" binding:"required,gt=0"`
}

type updateOrderRequest struct {
	Name   string  `json:"name" binding:"required"`
	Total  float64 `json:"total" binding:"required,gt=0"`
	Status string  `json:"status" binding:"required"`
}

func NewHandler(service service.OrderService) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(group *gin.RouterGroup) {
	orderGroup := group.Group("/orders")
	orderGroup.GET("/:id", handler.GetOrderByID)
	orderGroup.POST("/", handler.CreateOrder)
	orderGroup.PUT("/:id", handler.UpdateOrder)
	orderGroup.DELETE("/:id", handler.DeleteOrder)
	orderGroup.GET("/", handler.GetAllOrders)
}

func requestContext(ctx *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Request.Context(), requestTimeout)
}

func (handler *Handler) GetOrderByID(ctx *gin.Context) {
	orderID := ctx.Param("id")

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	order, err := handler.service.GetOrderByID(reqCtx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

func (handler *Handler) CreateOrder(ctx *gin.Context) {
	var request createOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	order, err := handler.service.CreateOrder(reqCtx, request.ClientID, request.Name, request.Total)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

func (handler *Handler) UpdateOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")
	var request updateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := domain.Status(request.Status)
	if !status.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	existingOrder, err := handler.service.GetOrderByID(reqCtx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	existingOrder.Name = request.Name
	existingOrder.Total = request.Total
	existingOrder.Status = domain.Status(request.Status)

	if err := handler.service.UpdateOrder(reqCtx, existingOrder); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}
	ctx.JSON(http.StatusOK, existingOrder)
}

func (handler *Handler) DeleteOrder(ctx *gin.Context) {
	orderID := ctx.Param("id")

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	if err := handler.service.DeleteOrder(reqCtx, orderID); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (handler *Handler) GetAllOrders(ctx *gin.Context) {
	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	orders, err := handler.service.GetAllOrders(reqCtx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
