package main

import (
	"log"

	"github.com/order_management/user_service/internal"
	"github.com/order_management/user_service/internal/configs"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error to load config file: %v", err)
	}

	inrnal := internal.NewServer(*config)
	inrnal.Start()
}
