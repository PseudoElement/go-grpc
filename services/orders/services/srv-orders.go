package services

import (
	"context"
	"math/rand"

	pb_orders "github.com/pseudoelement/go-grpc/protobuf/orders/generated"
)

type OrdersGrpcHandler struct {
	ordersDb []*pb_orders.Order
	pb_orders.UnimplementedOrderServiceServer
}

func NewGrpcOrdersService() *OrdersGrpcHandler {
	gRPCHandler := &OrdersGrpcHandler{ordersDb: make([]*pb_orders.Order, 0)}
	return gRPCHandler
}

func (h *OrdersGrpcHandler) GetOrders(ctx context.Context, req *pb_orders.GetOrdersRequest) (*pb_orders.GetOrdersResponse, error) {
	// o := h.ordersService.GetOrders(ctx)
	res := &pb_orders.GetOrdersResponse{
		Orders: h.ordersDb,
	}

	println("[grpc_GetOrders]", h.ordersDb)

	return res, nil
}

func (h *OrdersGrpcHandler) CreateOrder(ctx context.Context, req *pb_orders.CreateOrderRequest) (*pb_orders.CreateOrderResponse, error) {
	oderID := int32(rand.Intn(1000))
	order := &pb_orders.Order{
		OrderID:    &oderID,
		CustomerID: req.CustomerID,
		ProductID:  req.ProductID,
		Quantity:   req.Quantity,
	}
	h.ordersDb = append(h.ordersDb, order)

	println("[grpc_CreateOrder]", order)
	status := "success"
	res := &pb_orders.CreateOrderResponse{
		Status: &status,
	}

	return res, nil
}
