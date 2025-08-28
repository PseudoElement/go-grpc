package main

func main() {
	grpcServer := NewGRPCServer(":9001")
	grpcServer.Run()
}
