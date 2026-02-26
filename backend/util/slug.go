package util

import (
	"crypto/rand"
	"math/big"
)

const (
	slugCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	slugLength  = 5
)

// GenerateSlug creates a cryptographically random base62 slug of the given length.
func GenerateSlug(length int) (string, error) {
	if length <= 0 {
		length = slugLength
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(slugCharset)))

	for i := range result {
		n, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = slugCharset[n.Int64()]
	}

	return string(result), nil
}
