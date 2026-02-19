package client

import (
	"context"
	"fmt"

	"github.com/wafi11/microservices/users-services/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	client proto.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpcs server: %v", err)
	}

	return &UserClient{
		client: proto.NewUserServiceClient(conn),
	}, nil
}

func (u *UserClient) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.UserResponse, error) {
	return u.client.RegisterUser(ctx, req)
}

func (u *UserClient) LoginUser(c context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	return u.client.LoginUser(c, req)
}
func (u *UserClient) FindMe(c context.Context, req *proto.FindMeRequest) (*proto.FindMeResponse, error) {
	return u.client.FindMe(c, req)
}
