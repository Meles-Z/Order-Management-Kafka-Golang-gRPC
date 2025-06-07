package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/order_management/product_service/configs"
	"github.com/order_management/product_service/internal/api"
	"github.com/order_management/product_service/internal/db"
	"github.com/order_management/product_service/internal/repository"
	"github.com/order_management/product_service/internal/service"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error to load configuration:%s", err)
	}
	db, err := db.InitDB(cfg.DBConfig)
	repo := repository.NewProductRepository(db)
	svc := service.NewServices(repo)
	api := api.NewAPiService(svc)

	e := echo.New()
	product := e.Group("/product")
	product.POST("/create", api.CreateProduct())
	e.Start(":8282")
}
