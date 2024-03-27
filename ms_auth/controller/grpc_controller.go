package controller

import (
	"context"
	"fmt"
	"ms_auth/pb"
	"ms_auth/service"
)

type GRPC_MSAuth_Interface interface {
	Login(ctx context.Context, in *pb.LoginMessage) (*pb.LoginResponse, error)
	CreateUser(ctx context.Context, in *pb.CreateUserMessage) (*pb.CreateUserResponse, error)
}

type Server_GRPC_MS_Auth struct {
	userService service.UserService
	pb.UnimplementedAuthenServer
}

func (s *Server_GRPC_MS_Auth) Login(ctx context.Context, in *pb.LoginMessage) (*pb.LoginResponse, error) {
	// log.Printf("=> Request Login From: %v\n", in.GetEmail())
	result, err := s.userService.Login(in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Name:  result.Name,
			Email: result.Email,
		},
		SessionId:   fmt.Sprintf("%d", result.ID),
		AccessToken: fmt.Sprintf(">%d - Access: %s connect to system", result.ID, result.Email),
	}, nil
}

func (s *Server_GRPC_MS_Auth) CreateUser(ctx context.Context, in *pb.CreateUserMessage) (*pb.CreateUserResponse, error) {
	// log.Printf("=> Request register: %v\n", in.GetName())
	result, err := s.userService.CreateUser(in.GetName(), in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("=> Error: %v", err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Name:     result.Name,
			Email:    result.Email,
			Password: result.Password,
		},
	}, nil
}

func NewServer_GRPC_MS_Auth() *Server_GRPC_MS_Auth {
	return &Server_GRPC_MS_Auth{
		userService: service.NewUserService(),
	}
}
