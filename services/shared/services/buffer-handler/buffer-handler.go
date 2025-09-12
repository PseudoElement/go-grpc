package bufferhandler

import (
	"io"
	"log"
	"strings"
	"sync"

	pb_bufferhandler "github.com/pseudoelement/go-grpc/protobuf/buffer-handler/generated"
	"github.com/pseudoelement/go-grpc/services/shared/services/buffer-handler/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BufferHandlerService struct {
	clients  map[string][]byte
	filesSrv *services.FilesSrv
	mu       *sync.Mutex
	pb_bufferhandler.UnimplementedBufferHandlerServer
}

func NewGrpcBufferHandlerService() *BufferHandlerService {
	return &BufferHandlerService{clients: make(map[string][]byte), filesSrv: services.NewFilesSrv()}
}

func (s *BufferHandlerService) UploadBufferChunks(
	stream grpc.ClientStreamingServer[pb_bufferhandler.UploadBufferChunksReq, pb_bufferhandler.UploadBufferChunksResp],
) error {
	first := true
	loop := true
	go func() {
		// called when ctx cancel() invoked on client
		<-stream.Context().Done()
		log.Println("<-stream.Context().Done() on SERVER")
		loop = false
	}()

	for loop {
		inMsg, err := stream.Recv()
		if err == io.EOF {
			println("[BufferHandlerService_UploadBufferChunks] io.EOF")
			resp := &pb_bufferhandler.UploadBufferChunksResp{Done: true, Error: nil}
			return stream.SendAndClose(resp)
		}

		if err != nil {
			println("[BufferHandlerService_UploadBufferChunks] stream.Recv() err -", err)
			errMsg := err.Error()
			resp := &pb_bufferhandler.UploadBufferChunksResp{Done: false, Error: &errMsg}
			return stream.SendAndClose(resp)
		}
		if inMsg.Name == nil {
			errMsg := "file without name not allowed"
			resp := &pb_bufferhandler.UploadBufferChunksResp{Done: false, Error: &errMsg}
			return stream.SendAndClose(resp)
		}

		if first {
			log.Println("FIRST CALLED")
			s.clients[*inMsg.Name] = make([]byte, 0)
			first = false
		}

		s.clients[*inMsg.Name] = append(s.clients[*inMsg.Name], inMsg.Chunk...)
		println("file size -", len(s.clients[*inMsg.Name]))

		if inMsg.Last {
			println("[BufferHandlerService_UploadBufferChunks] inMsg.Last")
			fileBytes := s.clients[*inMsg.Name]
			println("total file size -", len(fileBytes))
			splitted := strings.Split(*inMsg.Name, "/")
			fileName := splitted[len(splitted)-1]
			fileErr := s.filesSrv.CreateAndWriteFile(fileBytes, fileName)
			if fileErr != nil {
				fileErrPtr := fileErr.Error()
				resp := &pb_bufferhandler.UploadBufferChunksResp{Done: true, Error: &fileErrPtr}
				return stream.SendAndClose(resp)
			}
		}
	}

	return stream.SendAndClose(&pb_bufferhandler.UploadBufferChunksResp{Done: false, Error: nil})
}

func (s *BufferHandlerService) DownloadBufferChunks(
	req *pb_bufferhandler.DownloadBufferChunksReq,
	stream grpc.ServerStreamingServer[pb_bufferhandler.DownloadBufferChunksResp],
) error {
	return status.Errorf(codes.Unimplemented, "method DownloadBufferChunks not implemented")
}
