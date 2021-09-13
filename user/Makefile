
build-grpc:
	protoc -I=./internal/transport/grpc/ --go_out=./ --go-grpc_out=./ ./internal/transport/grpc/user.proto

run-server:
	go run internal/cmd/server/*.go