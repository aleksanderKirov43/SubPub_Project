package main

import (
	"SubPub_project/internal/config"
	"SubPub_project/internal/server"
	"SubPub_project/pkg/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.NewLogger().Error("Ошибка загрузки конфигурации: %v", err)
	}
	server.Run(cfg)
}
