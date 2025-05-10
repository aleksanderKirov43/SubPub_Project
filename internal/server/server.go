package server

import (
	"SubPub_project/internal/app"
	"SubPub_project/internal/config"
	"SubPub_project/pkg/logger"
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	fmt.Println(fmt.Sprintf("Публикация: key = %s, data = %s", req.Key, req.Data)) // Для проверки Postman
	err := s.pubsub.Publish(ctx, req.Key, req.Data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка публикации: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func Run(ctx context.Context, cfg *config.Config) {
	log := logger.NewLogger()

	// 1. Создание listener для gRPC
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Error("Не удалось прослушать gRPC порт: %v", err)
		return
	}

	// 2. Инициализация gRPC сервера
	appInstance := app.NewApp(ctx, log)
	grpcServer := grpc.NewServer()
	pb.RegisterPubSubServer(grpcServer, NewServer(appInstance.PubSub))

	// 3. Запуск HTTP Gateway сервера в отдельной горутине
	go func() {
		mux := runtime.NewServeMux()
		endpoint := fmt.Sprintf("localhost:%d", cfg.GRPCPort)
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

		if err := pb.RegisterPubSubHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			log.Error("Ошибка подключения HTTP Gateway к gRPC: %v", err)
			return
		}

		log.Info("HTTP сервер запущен на порту: %d", cfg.RESTPort)
		if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.RESTPort), mux); err != nil {
			log.Error("Ошибка HTTP сервера: %v", err)
		}
	}()

	// 4. Завершение по CTRL+C
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Info("Получен сигнал завершения, выключаем gRPC...")
		grpcServer.GracefulStop()
	}()

	log.Info("gRPC сервер запущен на порту: %d", cfg.GRPCPort)
	if err := grpcServer.Serve(listener); err != nil {
		log.Error("Ошибка gRPC сервера: %v", err)
	}
}
