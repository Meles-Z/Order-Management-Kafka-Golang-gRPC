package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/db"
	"github.com/order_management/user_service/internal/handlers"
	kafkamessage "github.com/order_management/user_service/internal/kafka"
	"github.com/order_management/user_service/internal/repository"
	"github.com/order_management/user_service/internal/routes"
	"github.com/order_management/user_service/internal/services"
	"gorm.io/gorm"
)

// Server represents the server instance with all necessary dependencies.
type Server struct {
	DB          *gorm.DB
	Config      configs.Config
	UserHandler handlers.UserHandler
}

// NewServer creates a new server instance with all the dependencies injected.
func NewServer(cfg configs.Config) (*Server, error) {
	// Initialize database connection
	dbConn, err := db.InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Set up repositories, services, and handlers
	userRepo := repository.NewUserRepository(dbConn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Server{
		DB:          dbConn,
		Config:      cfg,
		UserHandler: userHandler,
	}, nil
}

// Start initializes the Echo server and starts listening on the configured port.
// It also starts the Kafka consumer in a goroutine and handles graceful shutdown.
func (s *Server) Start() error {
	e := echo.New()

	// Set up all routes
	routes.SetupRoutes(e, s.UserHandler)

	// Create a context that cancels on SIGINT or SIGTERM for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start Kafka consumer in the background with context using generic kafkamessage consumer

	go kafkamessage.StartKafkaConsumer(
		ctx,
		s.Config.Broker,
		s.Config.Topic,
		s.Config.Group,
		kafkamessage.HandleUserEvents,
	)

	// Start the HTTP server in a goroutine so we can listen for shutdown signals
	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- e.Start(fmt.Sprintf(":%d", s.Config.ServerPort))
	}()

	// Wait for either the server to stop or a shutdown signal
	select {
	case err := <-serverErrChan:
		return err // server stopped with error or normally
	case <-ctx.Done():
		log.Println("[Server] Shutdown signal received, exiting...")
		// Add any additional cleanup here if needed
		return nil
	}
}
