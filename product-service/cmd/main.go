package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/product_service/configs"
	"github.com/order_management/product_service/internal/api"
	"github.com/order_management/product_service/internal/db"
	"github.com/order_management/product_service/internal/kafka"
	"github.com/order_management/product_service/internal/repository"
	"github.com/order_management/product_service/internal/service"
	pkg "github.com/order_management/product_service/pkg/validator"
)

func main() {
	// Load configuration
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Initialize database
	db, err := db.InitDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}

	// Initialize Kafka producer with context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	producer, err := kafka.NewProducer(1024, kafka.WithWorkerCount(6))
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	repo := repository.NewProductRepository(db)
	svc := service.NewServices(repo)
	handler := api.NewAPiService(svc, producer)

	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}
	// Register routes
	api.RegisterRoutes(e, handler)
	
	//  Start server using config values
	serverAddress := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", serverAddress)
	go func() {
		if err := e.Start(serverAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown with timeout
	ctx, shutdownCancel := context.WithTimeout(ctx, 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
