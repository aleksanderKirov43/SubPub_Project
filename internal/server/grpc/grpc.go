package grpc

import (
	"SubPub_project/internal/app"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"google.golang.org/grpc"

	"log"
	"net"
)

func RunServer(listener net.Listener) {
	grpcServer := grpc.NewServer()
	pubsub := subpub.NewSubPub()
	pb.RegisterPubSubServer(grpcServer, app.NewServer(pubsub))

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка gRPC сервера: %v", err)
	}
}
