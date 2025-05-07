package main

import (
	"SubPub_project/internal/app"
	"SubPub_project/internal/config"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runRESTServer(restPort string, grpcPort string) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterPubSubHandlerFromEndpoint(context.Background(), mux, "localhost:"+grpcPort, opts)
	if err != nil {
		return err
	}

	log.Println("REST API запущен на порту:", restPort)
	return http.ListenAndServe(":"+restPort, mux)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	pubsub := subpub.NewSubPub()
	pb.RegisterPubSubServer(grpcServer, app.NewServer(pubsub))

	go func() {
		if err := runRESTServer(fmt.Sprintf("%d", cfg.RESTPort), fmt.Sprintf("%d", cfg.GRPCPort)); err != nil {
			log.Fatalf("Ошибка запуска REST сервера: %v", err)
		}
	}()

	log.Println("gRPC сервер запущен на порту:", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}
