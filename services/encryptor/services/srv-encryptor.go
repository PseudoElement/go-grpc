package services

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strconv"

	pb_encryptor "github.com/pseudoelement/go-grpc/protobuf/encryptor/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	decimalNum, err := strconv.Atoi(req.DecimalStr) // Must be int64 for FormatInt
	hex := strconv.FormatInt(int64(decimalNum), 16) // "ff"
	s.operationsCount++

	return &pb_encryptor.DecimalToHexResp{HexStr: hex}, err
}
func (s *EncryptorService) DecimalToHexStream(data grpc.BidiStreamingServer[pb_encryptor.DecimalToHexReq, pb_encryptor.DecimalToHexResp]) error {
	s.operationsCount++

	return status.Errorf(codes.Unimplemented, "method DecimalToHexStream not implemented")
}
