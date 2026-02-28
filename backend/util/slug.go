package util

import (
	"crypto/sha256"
	"encoding/binary"
	mathrand "math/rand"
)

const (
	slugAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	slugBase     = int64(len(slugAlphabet)) // 62
	slugLength   = 5
)

// SlugGenerator encodes a monotonic counter into a fixed-length, non-sequential
// slug by mapping each base-62 digit through a secret-key–derived alphabet permutation.
type SlugGenerator struct {
	shuffled []byte
	length   int
}

// NewSlugGenerator returns a SlugGenerator whose alphabet permutation is derived
// deterministically from secretKey.  length is the number of characters in each
// generated slug (defaults to 5 when ≤ 0).
func NewSlugGenerator(secretKey string, length int) *SlugGenerator {
	if length <= 0 {
		length = slugLength
	}
	return &SlugGenerator{
		shuffled: shuffleAlphabet(secretKey),
		length:   length,
	}
}

// FromCounter encodes counter into a slug of g.length characters.
// The mapping is bijective: every counter value produces a unique slug.
func (g *SlugGenerator) FromCounter(counter int64) string {
	encoded := toBase62Fixed(counter, g.length)
	result := make([]byte, g.length)
	for i, ch := range []byte(encoded) {
		for j := range slugBase {
			if slugAlphabet[j] == ch {
				result[i] = g.shuffled[j]
				break
			}
		}
	}
	return string(result)
}

// toBase62Fixed encodes n as a fixed-length base-62 string, left-padding with
// slugAlphabet[0] ('0') when n is small.
func toBase62Fixed(n int64, length int) string {
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		result[i] = slugAlphabet[n%slugBase]
		n /= slugBase
	}
	return string(result)
}

// shuffleAlphabet performs a deterministic Fisher-Yates shuffle of slugAlphabet
// seeded by the SHA-256 hash of secretKey.
func shuffleAlphabet(secretKey string) []byte {
	hash := sha256.Sum256([]byte(secretKey))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))

	//nolint:gosec // deterministic shuffle, not used for cryptographic randomness
	r := mathrand.New(mathrand.NewSource(seed))

	shuffled := []byte(slugAlphabet)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}
