package main

import (
	"SubPub_project/internal/app"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func runRESTServer() error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterPubSubHandlerFromEndpoint(context.Background(), mux, "localhost:8082", opts)
	if err != nil {
		return err
	}

	log.Println("REST API запущен на :8081")
	return http.ListenAndServe(":8081", mux)
}

func main() {

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Не удалось прослушать порт: %v", err)
	}

	grpcServer := grpc.NewServer()
	pubsub := subpub.NewSubPub()
	pb.RegisterPubSubServer(grpcServer, app.NewServer(pubsub))

	go func() {
		if err := runRESTServer(); err != nil {
			log.Fatalf("Ошибка запуска REST сервера: %v", err)
		}
	}()

	log.Println("gRPC сервер читает порт:8082")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

}
