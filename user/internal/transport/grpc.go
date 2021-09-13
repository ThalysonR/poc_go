package transport

import (
	"context"
	"fmt"
	"net"

	"github.com/thalysonr/poc_go/common/log"
	mygrpc "github.com/thalysonr/poc_go/user/grpc"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type userService struct {
	mygrpc.UnimplementedUserServiceServer
	userService service.UserService
}

func (u *userService) FindAll(ctx context.Context, empty *emptypb.Empty) (*mygrpc.UserResponse, error) {
	users, err := u.userService.FindAll(ctx)
	if err != nil {
		log.GetLogger().Error("could not find users: %w", err)
		return nil, err
	}

	var gUsers []*mygrpc.User
	for _, user := range users {
		gUsers = append(gUsers, &mygrpc.User{
			ID:        uint32(user.ID),
			BirthDate: user.BirthDate,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}
	log.GetLogger().Debug("Users found: ", users)
	return &mygrpc.UserResponse{
		Users: gUsers,
	}, nil
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
