package infrastructure

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
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

func loadKeyPemParam() error {
	// Load privateKey
	privateReader, err := ioutil.ReadFile("./infrastructure/private.pem")
	if err != nil {
		log.Println("No RSA private pem file: ", err)
		return err
	}

	privatePem, _ := pem.Decode(privateReader)
	privateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)

	// Load publicKey
	publicReader, err := ioutil.ReadFile("./infrastructure/public.pem")
	if err != nil {
		log.Println("No RSA public pem file: ", err)
		return err
	}

	publicPem, _ := pem.Decode(publicReader)
	publicKey, _ = x509.ParsePKIXPublicKey(publicPem.Bytes)
	return nil
}
