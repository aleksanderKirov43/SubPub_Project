package main

import (
	"SubPub_project/internal/config"
	"SubPub_project/internal/server"
	"SubPub_project/pkg/logger"

	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.NewLogger().Error("Ошибка загрузки конфигурации: %v", err)
	}
	server.Run(ctx, cfg)
}
