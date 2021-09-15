package grpc

import (
	"fmt"
	"net"
	"reflect"

	mygrpc "github.com/thalysonr/poc_go/user/grpc"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"github.com/thalysonr/poc_go/user/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GrpcServer struct {
	cfg         config.Config
	listener    net.Listener
	server      *grpc.Server
	userService service.UserService
}

func NewGrpcServer(userService service.UserService) *GrpcServer {
	return &GrpcServer{
		userService: userService,
	}
}

func (g *GrpcServer) ConfigChanged(cfg config.Config) bool {
	return !reflect.DeepEqual(g.cfg.Server.Grpc, cfg.Server.Grpc)
}

func (g *GrpcServer) Start(cfg config.Config) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Grpc.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	g.listener = lis

	server := grpc.NewServer()
	g.server = server
	g.cfg = cfg
	mygrpc.RegisterUserServiceServer(server, &userService{
		userService: g.userService,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (g *GrpcServer) Stop() error {
	if g.server != nil {
		g.server.Stop()
	}
	if g.listener != nil {
		g.listener.Close()
	}
	return nil
}
