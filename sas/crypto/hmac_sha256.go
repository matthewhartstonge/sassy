package crypto

import (
	// Standard Library Imports
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// HMACSHA256 generates a base64 encoded HMAC-SHA256 authentication code.
func HMACSHA256(key []byte, message []byte) string {
	hash := hmac.New(sha256.New, key)
	hash.Write(message)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
