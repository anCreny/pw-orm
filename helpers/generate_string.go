package helpers

import "math/rand/v2"

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		// Randomly pick a character from the charset
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
