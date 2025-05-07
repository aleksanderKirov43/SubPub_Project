package app

import (
	"SubPub_project/pkg/subpub"
	pb "SubPub_project/proto"

	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedPubSubServer
	pubsub subpub.SubPub
}

func NewServer(pubsub subpub.SubPub) *Server {
	return &Server{
		pubsub: pubsub,
	}
}

func (s *Server) Subscribe(req *pb.SubscribeRequest, stream pb.PubSub_SubscribeServer) error {
	sub, err := s.pubsub.Subscribe(req.Key, func(msg interface{}) {
		if str, ok := msg.(string); ok {
			log.Println("Отправка события подписчик:", str)
			_ = stream.Send(&pb.Event{Data: str})
		}
	})
	if err != nil {
		return status.Errorf(codes.Internal, "Ошибка приложения: %v", err)
	}
	<-stream.Context().Done()
	sub.Unsubscribe()
	return nil
}

func (s *Server) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	log.Printf("Публикация: key = %s, data = %s", req.Key, req.Data) // Для проверки Postman
	err := s.pubsub.Publish(req.Key, req.Data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Ошибка публикации: %v", err)
	}
	return &emptypb.Empty{}, nil
}
