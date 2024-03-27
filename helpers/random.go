package helpers

import (
	"math/rand"
	"time"
)

func GenerateRandomCode() string {
	rand.NewSource(time.Now().UnixNano())

	// Define the character set from which the code will be generated.
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}
