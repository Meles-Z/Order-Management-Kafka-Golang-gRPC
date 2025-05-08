package main

import (
	"log"

	"github.com/order_management/user_service/internal/configs"
	"github.com/order_management/user_service/internal/server"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize and start the server
	server, err := server.NewServer(*config)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
