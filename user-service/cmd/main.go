package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/server"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the server
	srv, err := server.NewServer(*config)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// If your server has a Stop or Shutdown method, call it here
	// e.g., srv.Stop() or srv.GracefulShutdown()

	log.Println("Server shut down gracefully.")
}
