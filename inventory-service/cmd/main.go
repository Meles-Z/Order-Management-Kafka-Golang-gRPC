package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/order_management/iventory_service/configs"
	db "github.com/order_management/iventory_service/internal/db"
	"github.com/order_management/iventory_service/internal/message"
	"github.com/order_management/iventory_service/internal/repository"
	"github.com/order_management/iventory_service/internal/service"
	"github.com/order_management/iventory_service/pkg/logger"
	pkg "github.com/order_management/iventory_service/pkg/validate"
)

func main() {
	kafkaBootstrap := getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092")
	productTopic := getEnv("KAFKA_TOPIC", "product-topic")
	productGroupID := getEnv("KAFKA_GROUP_ID", "inventory-service-group")
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	dbConn, err := db.InitDB(cfg.DB)
	if err != nil {
		logger.Error("Failed to load db connection", "error", err)
	}

	inventoryRepo := repository.NewRepostiory(dbConn)
	inventorySvc := service.NewServices(inventoryRepo)

	if err := logger.Init(cfg.ENV.Env); err != nil {
		panic("logger init failed: " + err.Error())
	}
	defer logger.Sync()
	// 🧠 Only start the consumer AFTER the logger and db are ready
	if err := message.StartConsumeProduct(kafkaBootstrap, productGroupID, productTopic, inventorySvc); err != nil {
		log.Fatalf("❌ Failed to start product consumer: %v", err)
	}

	e := echo.New()
	e.Validator = &pkg.CustomValidator{Validator: validator.New()}

	// Health routes
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "✅ Inventory service running")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	port := "8383"
	// Start server in goroutine
	go func() {
		log.Printf("🚀 HTTP server running on port %s", port)
		if err := e.Start(":" + port); err != nil {
			log.Fatalf("❌ Echo server stopped: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT or SIGTERM
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("🛑 Shutdown signal received, stopping server...")
	if err := e.Shutdown(nil); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Println("✅ Server exited properly")

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
