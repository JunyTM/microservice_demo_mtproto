package infrastructure

import (
	"context"
	"ms_gmail/pb"
	"reflect"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func loadHandshake() error {
	// Contact to server auth
	conn, err := grpc.Dial(serverHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	clientAuthRPC := pb.NewAuthenClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	publicKeyStr := reflect.ValueOf(publicKey).String()
	serverResponse, err := clientAuthRPC.Handshake(ctx, &pb.HandshakeRequest{PublicKey: publicKeyStr})
	if err != nil {
		return err
	}
	authKey = serverResponse.AuthKey
	serverPublicKey = serverResponse.PublicKey
	return nil
}
