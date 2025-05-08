package grpc

import (
	"SubPub_project/internal/app"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

func RunServer(listener net.Listener) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcServer := grpc.NewServer()
	pubsub := subpub.NewSubPub()
	server := app.NewServer(ctx, pubsub)

	pb.RegisterPubSubServer(grpcServer, server)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка gRPC сервера: %v", err)
	}
}
