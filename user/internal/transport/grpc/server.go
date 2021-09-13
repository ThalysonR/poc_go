package grpc

import (
	"fmt"
	"net"

	mygrpc "github.com/thalysonr/poc_go/user/grpc"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	server      *grpc.Server
	userService service.UserService
}

func NewGrpcServer(userService service.UserService) *GrpcServer {
	return &GrpcServer{
		userService: userService,
	}
}

func (g *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := userService{
		userService: g.userService,
	}

	server := grpc.NewServer()
	g.server = server
	mygrpc.RegisterUserServiceServer(server, &s)
	reflection.Register(server)
	return server.Serve(lis)
}

func (g *GrpcServer) Stop() error {
	if g.server != nil {
		g.server.Stop()
	}
	return nil
}
