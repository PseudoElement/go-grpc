package main

import (
	"log"
	"net"

	pb_encryptor "github.com/pseudoelement/go-grpc/protobuf/encryptor/generated"
	"github.com/pseudoelement/go-grpc/services/encryptor/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	addr string
}

func NewGRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// init services
	// YOU CAN INIT SERVICES FROM DIFFERENT .proto files
	encryptorSrv := services.NewGrpcEncryptorService()

	// register services
	pb_encryptor.RegisterEncryptorServer(grpcServer, encryptorSrv)
	// ... here register other services like routes in REST

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
