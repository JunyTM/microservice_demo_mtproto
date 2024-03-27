package infrastructure

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
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

	loadHandshake()
}

func loadKeyPemParam() error {
	// Load privateKey
	privateReader, err := ioutil.ReadFile("./private.pem")
	if err != nil {
		log.Println("No RSA private pem file: ", err)
		return err
	}

	privatePem, _ := pem.Decode(privateReader)
	privateKey, err = x509.ParsePKCS1PrivateKey(privatePem.Bytes)

	// Load publicKey
	publicReader, err := ioutil.ReadFile("./public.pem")
	if err != nil {
		log.Println("No RSA public pem file: ", err)
		return err
	}

	publicPem, _ := pem.Decode(publicReader)
	publicKey, _ = x509.ParsePKIXPublicKey(publicPem.Bytes)
	return nil
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
