package main

import (
	"log"
	"ms_auth/controller"
	"ms_auth/pb"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthenServer(s, controller.NewServer_GRPC_MS_Auth())
	log.Printf("Starting micro_auth: port - %s\n", "9090")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
