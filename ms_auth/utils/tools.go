package utils

import (
	"crypto/rand"
	"io"
	"time"
	
	uuid "github.com/satori/go.uuid"
	math "math/rand"
)

func init() {
	math.Seed(time.Now().UnixNano())
}

func GenSalt(n int) []byte {
	nonce := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	return (nonce)
}

func RandomString(length int) string {
	myuuid := uuid.NewV4().String()
	return myuuid[:length]
}

func RandomInt(min, max int64) int64 {
	return math.Int63n(max-min+1) + min
}
