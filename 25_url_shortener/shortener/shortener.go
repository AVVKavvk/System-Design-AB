package shortener

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"

// ShortenURL generates a unique short key from a long URL.
// length controls how many characters the key has (min 4, max 32).
func ShortenURL(longURL string, length int) (string, error) {
	if length < 4 {
		return "", fmt.Errorf("length must be at least 4")
	}
	if length > 32 {
		return "", fmt.Errorf("length must be at most 32")
	}

	key := generateKey(longURL, length)
	shortURL := key
	return shortURL, nil
}

// generateKey hashes the URL + a timestamp salt using SHA-256,
// then maps the raw bytes onto the allowed charset.
// The timestamp salt ensures uniqueness across repeated calls for the same URL.
func generateKey(longURL string, length int) string {
	// Salt with nanosecond timestamp for uniqueness
	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	input := longURL + "|" + salt

	hash := sha256.Sum256([]byte(input)) // 32 bytes

	var builder strings.Builder
	base := uint64(len(charset)) // 63

	// Consume the hash in 8-byte (uint64) windows, folding when we run out
	for i := 0; i < length; i++ {
		// Cycle through the 32-byte hash using modular indexing
		offset := (i * 3) % (len(hash) - 7) // stride of 3 bytes to spread entropy
		window := binary.BigEndian.Uint64(hash[offset : offset+8])
		idx := window % base
		builder.WriteByte(charset[idx])
	}

	return builder.String()
}
