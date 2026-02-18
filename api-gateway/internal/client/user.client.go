package client

import (
	"context"

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
		return nil, err
	}

	return &UserClient{
		client: proto.NewUserServiceClient(conn),
	}, nil
}

func (u *UserClient) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.UserResponse, error) {
	return u.client.RegisterUser(ctx, req)
}
