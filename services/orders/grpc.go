package main

import (
	"log"
	"net"

	pb_orders "github.com/pseudoelement/go-grpc/protobuf/orders/generated"
	"github.com/pseudoelement/go-grpc/services/orders/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type OrdersServer struct {
	addr        string
	extCallsSrv services.OrdersExtCallsSrv
}

func NewOrdersServer(addr string) *OrdersServer {
	return &OrdersServer{addr: addr, extCallsSrv: services.NewOrdersExtCallsSrv()}
}

func (s *OrdersServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// init services
	ordersSrv := services.NewGrpcOrdersService()

	// register services
	pb_orders.RegisterOrderServiceServer(grpcServer, ordersSrv)

	reflection.Register(grpcServer)

	log.Println("Starting gRPC server on", s.addr)

	go s.extCallsSrv.UploadFileToBufferHandler()

	return grpcServer.Serve(lis)
}
