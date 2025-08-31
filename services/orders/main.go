package main

func main() {
	grpcServer := NewOrdersServer(":9000")
	grpcServer.Run()
}
