package internal

import (
	"context"

	"github.com/wafi11/microservices/users-services/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	proto.UnimplementedUserServiceServer
	service *UserService
}

func NewGrpcServer(service *UserService) *GrpcServer {
	return &GrpcServer{service: service}
}

func (s *GrpcServer) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.UserResponse, error) {
	// convert proto → internal type
	userReq := protoToUser(req)

	// validasi
	if err := userReq.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// panggil service
	user, err := s.service.RegisterUser(userReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert internal type → proto
	return userToProto(user), nil
}
func protoToUser(req *proto.RegisterRequest) UserRegister {
	password := req.Password
	return UserRegister{
		FullName:    req.FullName,
		Email:       req.Email,
		Password:    &password,
		PhoneNumber: req.PhoneNumber,
	}
}

func userToProto(user User) *proto.UserResponse {
	return &proto.UserResponse{
		Id:          user.ID,
		FullName:    user.FullName,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		IsActive:    user.IsActive,
	}
}
