package main

import (
	"github.com/order_management/iventory_service/configs"
	"github.com/order_management/iventory_service/internal/database"
	"github.com/order_management/iventory_service/pkg/logger"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	_, err = database.InitDB(cfg.DB)
	if err != nil {
		logger.Error("Failed to load database connection", "error", err)
	}
	if err := logger.Init(cfg.ENV.Env); err != nil {
		panic("logger init failed: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Inventory service started", "env", cfg.ENV.Env, "service", "inventory")
}
