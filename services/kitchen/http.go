package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"text/template"
	"time"

	pb_encryptor "github.com/pseudoelement/go-grpc/protobuf/encryptor/generated"
	pb_orders "github.com/pseudoelement/go-grpc/protobuf/orders/generated"
)

type httpServer struct {
	addr string
}

func NewHttpServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) Run() error {
	router := http.NewServeMux()

	ordersConn := NewGRPCClient(":9000")
	defer ordersConn.Close()

	encryptorConn := NewGRPCClient(":9001")
	defer encryptorConn.Close()

	ordersGRPCClient := pb_orders.NewOrderServiceClient(ordersConn)
	encryptorGRPCClient := pb_encryptor.NewEncryptorClient(encryptorConn)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		res, err := ordersGRPCClient.GetOrders(ctx, &pb_orders.GetOrdersRequest{})
		if err != nil {
			log.Fatalf("client error: %v", err)
		}

		t := template.Must(template.New("orders").Parse(ordersTemplate))

		if err := t.Execute(w, res.GetOrders()); err != nil {
			log.Fatalf("template error: %v", err)
		}
	})

	router.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		id := int32(rand.Intn(1000))
		productID := int32(rand.Intn(1000))
		log.Printf("id - %v, prod_id - %v\n", id, productID)

		resp, err := ordersGRPCClient.CreateOrder(ctx, &pb_orders.CreateOrderRequest{
			CustomerID: id,
			ProductID:  productID,
			Quantity:   1,
		})

		if err != nil {
			log.Fatalf("create error: %v", err)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(resp)
	})

	router.HandleFunc("/get-orders", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		orders, _ := ordersGRPCClient.GetOrders(ctx, nil)

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(orders)
	})

	router.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		decimalValue := r.URL.Query().Get("decimal")
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		hexResp, _ := encryptorGRPCClient.DecimalToHex(ctx, &pb_encryptor.DecimalToHexReq{
			DecimalStr: decimalValue,
		})
		encryptResp, _ := encryptorGRPCClient.Encrypt(ctx, &pb_encryptor.EncryptReq{
			EncryptionType: 0,
			Value:          "EncodeToString",
		})

		resp := []any{hexResp, encryptResp}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}

var ordersTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kitchen Orders</title>
</head>
<body>
    <h1>Orders List</h1>
	<button onclick="createOrder()">Create order</button>
	<button onclick="getOrders()">Get orders</button>
	<br>
	<input type="string" />
	<button onclick="encrypt()">Encrypt</button>

    <table border="1">
        <tr>
            <th>Order ID</th>
            <th>Customer ID</th>
            <th>Quantity</th>
        </tr>
        {{range .}}
        <tr>	
            <td>{{.OrderID}}</td>
            <td>{{.CustomerID}}</td>
            <td>	{{.Quantity}}</td>
        </tr>
        {{end}}
    </table>

	<script>
		const input = document.querySelector('input');
		const createOrder = () => fetch('http://localhost:1000/create').then(r => r.json()).then(console.log)
		const getOrders = () => fetch('http://localhost:1000/get-orders').then(r => r.json()).then(console.log)
		const encrypt = () => fetch('http://localhost:1000/encrypt?decimal=' + input.value).then(r => r.json()).then(console.log)
	</script>
</body>
</html>`
