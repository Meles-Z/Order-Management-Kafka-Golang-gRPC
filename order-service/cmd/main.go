package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/order_management/order_service/internal/database"
)

func main() {
	fmt.Println("Starting Order Service...")

	_, err := database.InitDb()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ech := echo.New()
	ech.GET("/", func(c echo.Context) error {
		fmt.Println("New change happend!")
		return c.String(200, "Order Service is running.")
	})

	ech.Logger.Fatal(ech.Start(":" + port))
}
