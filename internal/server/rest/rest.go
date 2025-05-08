package rest

import (
	pb "SubPub_project/proto"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"context"
	"log"
	"net/http"
)

func RunServer(restPort int, grpcPort int) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterPubSubHandlerFromEndpoint(context.Background(), mux, "localhost:"+fmt.Sprintf("%d", grpcPort), opts)
	if err != nil {
		return err
	}

	log.Println("REST API запущен на порту:", restPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", restPort), mux)
}
