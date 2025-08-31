package main

import (
	"log"
	"net"

	pb_bufferhandler "github.com/pseudoelement/go-grpc/protobuf/buffer-handler/generated"
	pb_encryptor "github.com/pseudoelement/go-grpc/protobuf/encryptor/generated"

	bufferhandler "github.com/pseudoelement/go-grpc/services/shared/services/buffer-handler"
	"github.com/pseudoelement/go-grpc/services/shared/services/encryptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type SharedServer struct {
	addr string
}

func NewSharedServer(addr string) *SharedServer {
	return &SharedServer{addr: addr}
}

func (s *SharedServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// YOU CAN INIT SERVICES FROM DIFFERENT .proto files
	encryptorSrv := encryptor.NewGrpcEncryptorService()
	bufferHandlerSrv := bufferhandler.NewGrpcBufferHandlerService()

	// register services
	pb_encryptor.RegisterEncryptorServer(grpcServer, encryptorSrv)
	pb_bufferhandler.RegisterBufferHandlerServer(grpcServer, bufferHandlerSrv)

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on", s.addr)

	return grpcServer.Serve(lis)
}
