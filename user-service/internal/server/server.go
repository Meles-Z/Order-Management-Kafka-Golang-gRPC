package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	cfg "github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/db"
	"github.com/order_management/user_service/internal/handlers"
	"github.com/order_management/user_service/internal/repository"
	"github.com/order_management/user_service/internal/routes"
	"github.com/order_management/user_service/internal/services"
	"gorm.io/gorm"
)

// Server represents the server instance with all necessary dependencies.
type Server struct {
	DB          *gorm.DB
	Config      cfg.Config
	UserHandler handlers.UserHandler
}

// NewServer creates a new server instance with all the dependencies injected.
func NewServer(cfg cfg.Config) (*Server, error) {
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
func (s *Server) Start() error {
	e := echo.New()

	// Set up all routes
	routes.SetupRoutes(e, s.UserHandler)

	// Start the server
	return e.Start(fmt.Sprintf(":%d", s.Config.ServerPort))
}
