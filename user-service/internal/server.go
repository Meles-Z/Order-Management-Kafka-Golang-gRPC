package internal

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	cfg "github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/db"
	"github.com/order_management/user_service/internal/handlers"
	"github.com/order_management/user_service/internal/repository"
	"github.com/order_management/user_service/internal/services"
	"gorm.io/gorm"
)

// IStart defines the interface for starting the server.
type IStart interface {
	Start() error
}

// Server represents the server instance.
type Server struct {
	DB          *gorm.DB
	c           cfg.Config
	userHandler handlers.UserHandler
}

// NewServer creates a new server instance.
func NewServer(cfg cfg.Config) IStart {
	dbconn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error to initialize db: %v", err)
	}

	userRepo := repository.NewUserRepository(dbconn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	return &Server{
		DB:          dbconn,
		c:           cfg,
		userHandler: userHandler,
	}
}

// Start starts the server.
func (s *Server) Start() error {
	e := echo.New()
	Routes(e, *s)
	return e.Start(fmt.Sprintf(":%d", s.c.ServerPort))
}
