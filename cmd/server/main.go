package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"tf-grpc-svc/pkg/pb"
)

const (
	DefaultPort = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	grpcEndpoint := fmt.Sprintf(":%s", port)
	log.Printf("gRPC endpoint [%s]", grpcEndpoint)

	grpcServer := grpc.NewServer()
	service := NewServer()
	pb.RegisterTerraformServer(grpcServer, service)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting: gRPC Listener [%s]\n", grpcEndpoint)
	log.Fatal(grpcServer.Serve(listen))
}
