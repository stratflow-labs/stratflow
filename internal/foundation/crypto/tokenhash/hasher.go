package tokenhash

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
)

var (
	ErrEmptyToken  = errors.New("token is empty")
	ErrEmptySecret = errors.New("secret is empty")
)

// Hash creates a deterministic HMAC-SHA256 hash of the token using the provided secret.
// The same token + secret combination will always produce the same hash.
// This allows fast lookups in the database while keeping tokens secure.
func Hash(token string, secret []byte) (string, error) {
	if token == "" {
		return "", ErrEmptyToken
	}
	if len(secret) == 0 {
		return "", ErrEmptySecret
	}

	// HMAC-SHA256(secret, token) - deterministic hash
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(token))
	hash := h.Sum(nil)

	return hex.EncodeToString(hash), nil
}

// Verify checks if the provided token matches the stored hash using constant-time comparison.
// Protects against timing attacks.
func Verify(token string, storedHash string, secret []byte) bool {
	if token == "" || storedHash == "" || len(secret) == 0 {
		return false
	}

	// Recompute hash with same secret
	computedHash, err := Hash(token, secret)
	if err != nil {
		return false
	}

	// Constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(storedHash)) == 1
}
