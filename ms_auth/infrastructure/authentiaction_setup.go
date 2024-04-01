package infrastructure

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"reflect"
)

const Algorithm = "HS256"

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
	publicKey_Any, _ := x509.ParsePKIXPublicKey(publicPem.Bytes)
	publicKey = reflect.ValueOf(publicKey_Any).String()
	return nil
}
