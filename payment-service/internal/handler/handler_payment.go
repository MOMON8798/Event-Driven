package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/domain"
	"github.com/MOMON8798/Event-Driven.git/payment-service/internal/service"
	"github.com/gin-gonic/gin"
)

const requestTimeout = 3 * time.Second

type Handler struct {
	service service.PaymentService
}

type createPaymentRequest struct {
	OrderID string  `json:"order_id" binding:"required"`
	Amount  float64 `json:"amount" binding:"required,gt=0"`
}

func NewHandler(service service.PaymentService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	payments := group.Group("/payments")
	payments.POST("/", h.CreatePayment)
	payments.GET("/:id", h.GetPaymentByID)
	payments.GET("/", h.GetPaymentsByOrderID)
}

func requestContext(ctx *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Request.Context(), requestTimeout)
}

func (h *Handler) CreatePayment(ctx *gin.Context) {
	var request createPaymentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	payment, err := h.service.ProcessPayment(reqCtx, request.OrderID, request.Amount)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidAmount) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}
	ctx.JSON(http.StatusCreated, payment)
}

func (h *Handler) GetPaymentByID(ctx *gin.Context) {
	id := ctx.Param("id")

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	payment, err := h.service.GetPaymentByID(reqCtx, id)
	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}
	ctx.JSON(http.StatusOK, payment)
}

func (h *Handler) GetPaymentsByOrderID(ctx *gin.Context) {
	orderID := ctx.Query("order_id")
	if orderID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "order_id query parameter is required"})
		return
	}

	reqCtx, cancel := requestContext(ctx)
	defer cancel()

	payments, err := h.service.GetPaymentsByOrderID(reqCtx, orderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve payments"})
		return
	}
	ctx.JSON(http.StatusOK, payments)
}
