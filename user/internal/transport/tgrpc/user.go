package tgrpc

import (
	"context"
	"fmt"

	"github.com/thalysonr/poc_go/common/log"
	mygrpc "github.com/thalysonr/poc_go/user/grpc"
	"github.com/thalysonr/poc_go/user/internal/app/model"
	"github.com/thalysonr/poc_go/user/internal/app/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	mygrpc.UnimplementedUserServiceServer
	userService service.UserService
}

func (u *userService) Create(ctx context.Context, createUser *mygrpc.CreateUser) (*mygrpc.UserID, error) {
	if createUser == nil {
		return nil, fmt.Errorf("user info is required")
	}
	user := model.User{
		BirthDate: createUser.BirthDate,
		Email:     createUser.Email,
		FirstName: createUser.FirstName,
		LastName:  createUser.LastName,
	}
	id, err := u.userService.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return &mygrpc.UserID{
		ID: uint32(id),
	}, nil
}

func (u *userService) Delete(ctx context.Context, userID *mygrpc.UserID) (*emptypb.Empty, error) {
	if userID == nil {
		return &emptypb.Empty{}, fmt.Errorf("user id is required")
	}

	err := u.userService.Delete(ctx, uint(userID.ID))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (u *userService) FindAll(ctx context.Context, empty *emptypb.Empty) (*mygrpc.UserResponse, error) {
	users, err := u.userService.FindAll(ctx)
	if err != nil {
		log.GetLogger().Error(ctx, "could not find users: %w", err)
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
	log.GetLogger().Debug(ctx, "Users found: ", users)
	return &mygrpc.UserResponse{
		Users: gUsers,
	}, nil
}
