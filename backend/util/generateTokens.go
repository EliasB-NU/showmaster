package util

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateSessionToken Generate a random session token (e.g., 32-byte token)
func GenerateSessionToken() string {
	bytes := make([]byte, 32) // 32 bytes = 64 characters when hex-encoded
	if _, err := rand.Read(bytes); err != nil {
		// Handle error
		// For simplicity, we return an empty string in case of error
		return ""
	}
	return hex.EncodeToString(bytes)
}
