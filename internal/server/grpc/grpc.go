package grpc

import (
	"SubPub_project/internal/app"
	"SubPub_project/internal/server"
	"SubPub_project/pkg/logger"
	pb "SubPub_project/proto"

	"context"
	"net"

	"google.golang.org/grpc"
)

func RunServer(listener net.Listener) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logInstance := logger.NewLogger()
	appInstance := app.NewApp(ctx, logInstance)

	grpcServer := grpc.NewServer()
	serverInstance := server.NewServer(appInstance.PubSub)

	pb.RegisterPubSubServer(grpcServer, serverInstance)

	if err := grpcServer.Serve(listener); err != nil {
		logger.NewLogger().Error("Ошибка gRPC сервера: %v", err)
	}
}
