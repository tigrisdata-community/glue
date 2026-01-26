package internal

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256sum computes a cryptographic hash. Still used for proof-of-work challenges
// where we need the security properties of a cryptographic hash function.
func SHA256sum(text string) string {
	hash := sha256.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}
