
build-grpc:
	protoc -I=./internal/transport/grpc/ --go_out=./ --go-grpc_out=./ ./internal/transport/grpc/user.proto

run-server:
	RUN_MODE=LOCAL go run internal/cmd/server/*.go

mockgen:
	mockery --dir internal/app/repository --all --output internal/mocks