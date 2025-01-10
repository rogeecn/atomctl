package users

import (
	"context"

	userv1 "{{.ModuleName}}/pkg/proto/user/v1"
)

// @provider(grpc) userv1.RegisterUserServiceServer
type Users struct{}

func (u *Users) ListUsers(ctx context.Context, in *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	// userv1.UserServiceServer
	return &userv1.ListUsersResponse{}, nil
}

// GetUser implements userv1.UserServiceServer
func (u *Users) GetUser(ctx context.Context, in *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id: in.Id,
		},
	}, nil
}
