run-orders:
	go run services/orders/*.go

run-kitchen:
	go run services/kitchen/*.go

gen-orders:
	protoc \
		--proto_path=protobuf "protobuf/orders/orders.proto" \
		--go_out=protobuf/orders/generated --go_opt=paths=source_relative \
  		--go-grpc_out=protobuf/orders/generated --go-grpc_opt=paths=source_relative

gen-encryptor:
	protoc \
		--proto_path=protobuf "protobuf/encryptor/encryptor.proto" \
		--go_out=protobuf/encryptor/generated --go_opt=paths=source_relative \
  		--go-grpc_out=protobuf/encryptor/generated --go-grpc_opt=paths=source_relative