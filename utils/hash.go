package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(input string) string {
	hash := sha256.New()                   // Create a new SHA-256 hash
	hash.Write([]byte(input))              // Add a string to the hash
	hashedBytes := hash.Sum(nil)           // We get the final hash as a byte slice
	return hex.EncodeToString(hashedBytes) // Convert bytes to a string in hex format
}
