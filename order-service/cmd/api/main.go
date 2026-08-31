package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MOMON8798/Event-Driven.git/internal/config"
	"github.com/MOMON8798/Event-Driven.git/internal/handler"
	"github.com/MOMON8798/Event-Driven.git/internal/repository"
	"github.com/MOMON8798/Event-Driven.git/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	cfg := config.Load()

	repo, err := repository.NewPostgresRepository(cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	svc := service.NewOrderService(repo)
	h := handler.NewHandler(svc)

	router := gin.Default()
	apiGroup := router.Group("/api")
	h.RegisterRoutes(apiGroup)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Printf("listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown finished with an error: %v", err)
	}

	if closer, ok := repo.(interface{ Close() }); ok {
		closer.Close()
		log.Println("database connections closed")
	}

	log.Println("server stopped")
}
