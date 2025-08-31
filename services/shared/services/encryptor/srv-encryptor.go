package encryptor

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strconv"
	"strings"

	pb_encryptor "github.com/pseudoelement/go-grpc/protobuf/encryptor/generated"
	encryptorutils "github.com/pseudoelement/go-grpc/services/shared/services/encryptor/utils"
	"google.golang.org/grpc"
)

type EncryptorService struct {
	operationsCount int
	pb_encryptor.UnimplementedEncryptorServer
}

func NewGrpcEncryptorService() *EncryptorService {
	return &EncryptorService{operationsCount: 0}
}

func (s *EncryptorService) Encrypt(ctx context.Context, req *pb_encryptor.EncryptReq) (*pb_encryptor.EncryptResp, error) {
	println("[grpc_Encrypt]", s.operationsCount)
	resp := &pb_encryptor.EncryptResp{Error: "", EncryptionType: req.EncryptionType}

	switch req.EncryptionType {
	case pb_encryptor.EncryptionType_SHA_256:
		hasher := sha256.New()
		hasher.Write([]byte(req.Value))
		resp.Encrypted = hex.EncodeToString(hasher.Sum(nil))
	case pb_encryptor.EncryptionType_MD_5:
		hasher := md5.New()             // Create a new MD5 hash object
		hasher.Write([]byte(req.Value)) // Write the string's bytes to the hasher
		resp.Encrypted = hex.EncodeToString(hasher.Sum(nil))
	case pb_encryptor.EncryptionType_BASE_64:
		resp.Encrypted = base64.StdEncoding.EncodeToString([]byte(req.Value))
	default:
		resp.Error = strconv.Itoa(int(req.EncryptionType)) + "is invalid encrypt type"
	}

	s.operationsCount++

	return resp, nil
}

func (s *EncryptorService) DecimalToHex(ctx context.Context, req *pb_encryptor.DecimalToHexReq) (*pb_encryptor.DecimalToHexResp, error) {
	println("[grpc_DecimalToHex]", s.operationsCount)

	hex, err := encryptorutils.DecimalToHex(req.DecimalStr)
	s.operationsCount++

	return &pb_encryptor.DecimalToHexResp{HexStr: hex}, err
}

func (s *EncryptorService) DecimalToHexStream(stream grpc.BidiStreamingServer[pb_encryptor.DecimalToHexStreamReq, pb_encryptor.DecimalToHexStreamResp]) error {
	for {
		inMsg, err := stream.Recv()
		log.Printf("[grpc_DecimalToHexStream] inMsg - %+v\n", inMsg)

		if err != nil && strings.Contains(err.Error(), "DeadlineExceeded") {
			println("[grpc_DecimalToHexStream] context done!")
			return nil
		}
		if inMsg == nil && err != nil {
			println("[grpc_DecimalToHexStream] unexpected nil message ==>", err.Error())
			return nil
		}

		hex, err := encryptorutils.DecimalToHex(inMsg.DecimalStr)
		if err != nil {
			log.Printf("[grpc_DecimalToHexStream] err != nil - %v\n", err)
			errPtr := err.Error()
			stream.Send(&pb_encryptor.DecimalToHexStreamResp{
				HexStr: hex,
				Error:  &errPtr,
				Stop:   true,
			})
			return nil
		}

		resp := &pb_encryptor.DecimalToHexStreamResp{
			HexStr: hex,
			Error:  nil,
			Stop:   false,
		}

		log.Printf("[grpc_DecimalToHexStream] Send - %v\n", resp)
		stream.Send(resp)

		if inMsg.Stop {
			return nil
		}
	}
}
