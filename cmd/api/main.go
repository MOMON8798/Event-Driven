package main

import (
	"github.com/MOMON8798/Event-Driven.git/internal/handler"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	repository := repository.NewInMemoryRepository()
	service := service.NewOrderService(repository)
	handler := handler.NewHandler(service)

	server := gin.Default()
	apiGroup := server.Group("/api")
	handler.RegisterRoutes(apiGroup)

	server.Run(":8080")
}
