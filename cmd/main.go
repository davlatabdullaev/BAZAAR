package main

import (
	"bazaar/api"
	"bazaar/config"
	"bazaar/pkg/logger"
	"bazaar/storage/postgres"
	"context"
)

func main() {

	cfg := config.Load()

	log := logger.New(cfg.ServiceName)

	pgStore, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("error while connecting to db", logger.Error(err))
		return
	}
	defer pgStore.CloseDB()

	server := api.New(pgStore, log)

	log.Info("Server is running on", logger.Int("port", 8080))
	if err = server.Run("localhost:8080"); err != nil {
		log.Error("error while server run")
		return
	}

}
