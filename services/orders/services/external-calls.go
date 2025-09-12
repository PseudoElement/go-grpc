package services

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb_bufhandler "github.com/pseudoelement/go-grpc/protobuf/buffer-handler/generated"
	grpcutils "github.com/pseudoelement/go-grpc/services/common/grpc"
)

type OrdersExtCallsSrv struct {
	bufHandlerGRPCClient pb_bufhandler.BufferHandlerClient
}

func NewOrdersExtCallsSrv() OrdersExtCallsSrv {
	extCallsSrv := OrdersExtCallsSrv{}

	sharedConn := grpcutils.NewGRPCClient(":9001")
	extCallsSrv.bufHandlerGRPCClient = pb_bufhandler.NewBufferHandlerClient(sharedConn)

	return extCallsSrv
}

func (s *OrdersExtCallsSrv) UploadFileToBufferHandler() error {
	<-time.After(3 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := s.bufHandlerGRPCClient.UploadBufferChunks(ctx)
	if err != nil {
		panic("[OrdersExtCallsSrv_UploadFileToBufferHandler] Open file error: " + err.Error())
	}

	pwd, _ := os.Getwd()
	path := pwd + "/services/orders/services/assets/borrow-image.jpg"
	println("[OrdersExtCallsSrv_UploadFileToBufferHandler] file path -", path)

	f, err := os.Open(path)
	if err != nil {
		panic("[OrdersExtCallsSrv_UploadFileToBufferHandler] Open file error: " + err.Error())
	}

	fileName := f.Name()
	bufSize := 1024
	buf := make([]byte, bufSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				println("file ended", n)
				err = stream.Send(&pb_bufhandler.UploadBufferChunksReq{
					Name:  &fileName,
					Last:  true,
					Chunk: buf,
				})
				break
			}
			panic("[OrdersExtCallsSrv_UploadFileToBufferHandler] Read file error: " + err.Error())
		}

		err = stream.Send(&pb_bufhandler.UploadBufferChunksReq{
			Name:  &fileName,
			Last:  false,
			Chunk: buf,
		})
		if err != nil {
			panic("[OrdersExtCallsSrv_UploadFileToBufferHandler] stream.Send() file error: " + err.Error())
		}

		println("read from file chunk ", n)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("[OrdersExtCallsSrv_UploadFileToBufferHandler] CloseAndRecv() error -", err)
	}
	log.Printf("[OrdersExtCallsSrv_UploadFileToBufferHandler] success reply on finish - %+v\n", reply)

	return nil
}
