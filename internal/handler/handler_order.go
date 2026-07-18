package handler

import (
	"net/http"

	"github.com/MOMON8798/Event-Driven.git/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.OrderService
}

type createOrderRequest struct {
	ClientID string  `json:"client_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Total    float64 `json:"total" binding:"required"`
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

func (handler *Handler) GetOrderByID(ctx *gin.Context) {
	UserID := ctx.Param("id")
	order, err := handler.service.GetOrderByID(UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
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
	order, err := handler.service.CreateOrder(request.ClientID, request.Name, request.Total)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	ctx.JSON(http.StatusCreated, order)
}

func (handler *Handler) UpdateOrder(ctx *gin.Context) {
	ID := ctx.Param("id")
	var request updateOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingOrder, err := handler.service.GetOrderByID(ID)
	if err != nil {
		if err == domain.ErrOrderNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}
	existingOrder.Name = request.Name
	existingOrder.Total = request.Total
	existingOrder.Status = domain.Status(request.Status)

	if err := handler.service.UpdateOrder(existingOrder); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}
	ctx.JSON(http.StatusOK, existingOrder)
}

func (handler *Handler) DeleteOrder(ctx *gin.Context) {
	ID := ctx.Param("id")
	if err := handler.service.DeleteOrder(ID); err != nil {
		if err == domain.ErrOrderNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (handler *Handler) GetAllOrders(ctx *gin.Context) {
	orders, err := handler.service.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}
