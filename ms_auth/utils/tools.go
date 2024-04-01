package utils

import (
	"crypto/rand"
	"io"
)

func GenSalt(n int) []byte {
	nonce := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return (nonce)
}

