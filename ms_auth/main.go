package main

import (
	"context"
	"fmt"
	"log"
	"ms_auth/pb"
	"ms_auth/service"
	"net"

	"google.golang.org/grpc"
)

// Implements the Server gRPC
type Server_GRPC_MS_Auth struct {
	userService service.UserService
	pb.UnimplementedAuthenServer
}

func (s *Server_GRPC_MS_Auth) Login(ctx context.Context, in *pb.LoginMessage) (*pb.LoginResponse, error) {
	log.Printf("=> Request Login From: %v\n", in.GetEmail())
	s.userService.Login(in.GetEmail(), in.GetPassword())



	var user *pb.User
	return &pb.LoginResponse{
		User:        user,
		SessionId:   "",
		AccessToken: "connection",
	}, nil
}

func (s *Server_GRPC_MS_Auth) Register(ctx context.Context, in *pb.CreateUserMessage) (*pb.CreateUserResponse, error) {
	log.Printf("=> Request register: %v\n", in.GetName())
	result, err := s.userService.CreateUser(in.GetName(), in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("=> Error: %v", err)
	}

	return &pb.CreateUserResponse{
		User:        &pb.User{
			Name:     result.Name,
            Email:    result.Email,
            Password: result.Password,
		},
		Code: "Success",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthenServer(s, &Server_GRPC_MS_Auth{})
	log.Printf("Starting micro_auth: port - %s\n", "9090")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
