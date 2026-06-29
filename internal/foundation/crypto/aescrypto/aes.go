package aescrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

const (
	// DefaultKey is a placeholder; replace with a secure value provided via configuration.
	DefaultKey = "8e71bbce7451ba2835de5aea73e4f3f96821455240823d2fd8174975b8321bfc!"
)

// Encrypt performs AES-GCM encryption using the provided hex-encoded key.
func Encrypt(keyHex, plaintext string) (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("decode key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("read nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(ciphertext), nil
}

// EncryptWithDefaultKey encrypts using DefaultKey.
func EncryptWithDefaultKey(plaintext string) (string, error) {
	return Encrypt(DefaultKey, plaintext)
}

// Decrypt reverses AES-GCM encryption using the provided hex-encoded key.
func Decrypt(keyHex, encrypted string) (string, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return "", fmt.Errorf("decode key: %w", err)
	}

	enc, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("new gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}

	return string(plaintext), nil
}

// DecryptWithDefaultKey decrypts using DefaultKey.
func DecryptWithDefaultKey(encrypted string) (string, error) {
	return Decrypt(DefaultKey, encrypted)
}
