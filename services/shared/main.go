package main

func main() {
	grpcServer := NewSharedServer(":9001")
	grpcServer.Run()
}
