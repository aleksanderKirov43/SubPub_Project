package server

import (
	"SubPub_project/internal/app"
	"SubPub_project/internal/config"
	"SubPub_project/internal/server/rest"
	"SubPub_project/pkg/logger"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"context"
	"fmt"
	"net"
)

type Server struct {
	pb.UnimplementedPubSubServer
	pubsub subpub.SubPubInterface
	log    logger.Logger
}

func NewServer(pubsub subpub.SubPubInterface) *Server {
	return &Server{
		pubsub: pubsub,
	}
}

func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	ctx := stream.Context()

	sub, err := s.pubsub.Subscribe(ctx, req.Key, func(msg interface{}) {
		if str, ok := msg.(string); ok {
			s.log.Info("Отправка события подписчику:", str)
			_ = stream.Send(&pb.Event{Data: str})
		}
	})
	if err != nil {
		return status.Errorf(codes.Internal, "Ошибка подписки: %v", err)
	}
	<-ctx.Done()
	sub.Unsubscribe()
	return nil
}

func (s *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	s.log.Info("Публикация: key = %s, data = %s", req.Key, req.Data) // Для проверки Postman
	err := s.pubsub.Publish(ctx, req.Key, req.Data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка публикации: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func Run(ctx context.Context, cfg *config.Config) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.NewLogger().Error("Не удалось прослушать порт: %v", err)
	}
	logger.NewLogger().Error("gRPC сервер запущен на порту:", cfg.GRPCPort)

	appInstance := app.NewApp(ctx, logger.NewLogger())
	serverInstance := NewServer(appInstance.PubSub)

	grpcServer := grpc.NewServer()
	pb.RegisterPubSubServer(grpcServer, serverInstance)

	go func() {
		if err = rest.RunServer(cfg.RESTPort, cfg.GRPCPort); err != nil {
			logger.NewLogger().Error("Ошибка запуска REST сервера: %v", err)
		}

		if err = grpcServer.Serve(listener); err != nil {
			logger.NewLogger().Error("Ошибка gRPC сервера: %v", err)
		}
	}()
}
