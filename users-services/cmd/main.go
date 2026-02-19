package main

import (
	"fmt"
	"log"
	"net"

	"github.com/wafi11/microservices/users-services/config"
	"github.com/wafi11/microservices/users-services/internal"
	"github.com/wafi11/microservices/users-services/proto"
	"google.golang.org/grpc"
)

func main() {
	db := config.NewDatabaseConfig()
	conn, err := db.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	// setup dependencies
	repo := internal.NewUserRepository(conn)
	service := internal.NewUserService(repo)

	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, internal.NewGrpcServer(service))

	log.Println("gRPC running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
