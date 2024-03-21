package utils

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxy"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	myuuid := uuid.NewV4().String()
	return myuuid[:length]
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}
