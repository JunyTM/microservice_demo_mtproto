package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"ms_gmail/pb"
	"time"
)

func main() {
	// Your code here
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial("0.0.0.0:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAuthenClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Login(ctx, &pb.LoginMessage{Email: "manhtokim@gmail.com", Password: "123456"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetAccessToken())

	
}
