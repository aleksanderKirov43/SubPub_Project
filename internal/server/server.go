package server

import (
	"SubPub_project/internal/config"
	"SubPub_project/internal/server/grpc"
	"SubPub_project/internal/server/rest"

	"fmt"
	"log"
	"net"
)

func Run(cfg *config.Config) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}
	log.Println("gRPC сервер запущен на порту:", cfg.GRPCPort)

	go func() {
		if err = rest.RunServer(cfg.RESTPort, cfg.GRPCPort); err != nil {
			log.Fatalf("Ошибка запуска REST сервера: %v", err)
		}
	}()

	grpc.RunServer(listener)
}
