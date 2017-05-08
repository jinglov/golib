package utils

import (
	"math/rand"
	"time"
)

func MakeRandomBytes(size int) (b []byte) {
	if size < 1 {
		return
	}
	randomChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b = make([]byte, 0, size)
	textNum := len(randomChars)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		b = append(b, randomChars[r.Intn(textNum)])
	}
	return b
}
