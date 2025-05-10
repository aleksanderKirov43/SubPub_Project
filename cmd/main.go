package main

import (
	"SubPub_project/internal/config"
	"SubPub_project/internal/server"
	"log"

	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Ошибка загрузки конфигурации: %v", err)
	}
	server.Run(ctx, cfg)
}
