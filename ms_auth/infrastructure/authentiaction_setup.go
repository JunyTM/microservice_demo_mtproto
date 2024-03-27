package infrastructure

import (
	"crypto/sha256"
	"crypto/x509"
	"fmt"
)

const Algorithm = "HS256"

// Authen key algorithm - Used for regist Client TPC
func GenerateAuthKeyFromPublicKey(publicKey []byte) (string, error) {
	//  Generate the auth key from the public key
	pubKey, err := x509.ParsePKCS1PublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("could not parse public key: %v", err)
	}

	// Marshal the public key to a byte slice
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)

	// Calculate SHA256 hash
	authKey := sha256.Sum256(pubKeyBytes)
	return string(authKey[:]), nil
}
