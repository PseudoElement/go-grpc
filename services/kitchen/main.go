package main

func main() {
	httpServer := NewHttpServer(":1000")
	httpServer.Run()
}
