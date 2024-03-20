package utils

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxy"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func RandomInt(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}
