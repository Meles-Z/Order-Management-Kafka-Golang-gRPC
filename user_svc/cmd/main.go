package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/user_svc/internal/db"
	"github.com/order_management/user_svc/internal/handler"
	"github.com/order_management/user_svc/internal/kafka"
	"github.com/order_management/user_svc/internal/repository"
	"github.com/order_management/user_svc/internal/services"
	pkg "github.com/order_management/user_svc/pkg/validator"
)

func main() {
	fmt.Println("Starting User Service...")
	file := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	fmt.Println(file)
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("DB error: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	producer, err := kafka.NewProducer(1024, kafka.WithWorkerCount(6))
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	repo := repository.NewUserRepository(dbConn)
	svc := services.NewUserService(repo)
	h := handler.NewHandler(svc, producer)
	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}
	handler.RegisterRoutes(e, h)

	//  Start server using config values
	serverAddress := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
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
