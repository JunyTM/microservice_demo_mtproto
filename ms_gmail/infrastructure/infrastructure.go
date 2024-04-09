package infrastructure

import (
	"crypto/rsa"
	"errors"
	"log"
)

var (
	serverHost = "localhost:9090"
	publicKey  interface{}
	privateKey *rsa.PrivateKey

	authKey         string
	serverPublicKey string
)

func init() {
	if err := loadKeyPemParam(); err != nil {
		log.Fatal(err)
	}

	if err := loadHandshake(); err != nil {
		log.Fatal(errors.New("=> Load handshake error: No authen key available"))
	}
	// log.Println("Handshake success")
	log.Println("Handshake success - Auth_key: ", authKey)
}

func GetServerHost() string {
	return serverHost
}

func GetAuthKey() string {
	return authKey
}

func GetServerPublicKey() string {
	return serverPublicKey
}
