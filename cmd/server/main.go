package main

import (
	"fmt"
	"log"
	"net"
	"os"

	productv1 "github.com/sekthor/grpc-streaming-example/api/product/v1"
	"github.com/sekthor/grpc-streaming-example/internal/service"
	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("PORT")
	service := service.ProductService{}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	productv1.RegisterProductServiceServer(grpcServer, service)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("could not start server")
	}
}
