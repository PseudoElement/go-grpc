package bufferhandler

import pb_bufferhandler "github.com/pseudoelement/go-grpc/protobuf/buffer-handler/generated"

type BufferHandlerService struct {
	pb_bufferhandler.UnimplementedBufferHandlerServer
}

func NewGrpcBufferHandlerService() *BufferHandlerService {
	return &BufferHandlerService{}
}
