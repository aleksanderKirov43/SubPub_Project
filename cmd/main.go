package main

import (
	"SubPub_project/internal/config"
	"SubPub_project/internal/server"

	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	server.Run(cfg)
}
