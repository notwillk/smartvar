package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// Decrypt decrypts an ENC:-prefixed value using the given algorithm and key.
// The value must start with "ENC:" followed by base64(nonce || ciphertext).
// Supported algorithms: aes256-gcm.
func Decrypt(algorithm, keyStr, value string) (string, error) {
	if !strings.HasPrefix(value, "ENC:") {
		return "", fmt.Errorf("value does not start with ENC:")
	}
	encoded := strings.TrimPrefix(value, "ENC:")

	switch algorithm {
	case "aes256-gcm":
		return decryptAES256GCM(keyStr, encoded)
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func decryptAES256GCM(keyStr, encoded string) (string, error) {
	key, err := decodeKey(keyStr, 32)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}

	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("invalid base64 ciphertext: %w", err)
	}

	// nonce is the first 12 bytes; remainder is ciphertext+tag
	if len(data) < 12+16 {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce := data[:12]
	ciphertext := data[12:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}
	return string(plaintext), nil
}

// EncryptAES256GCM encrypts plaintext and returns an ENC:-prefixed base64 string.
// Used in tests and the test fixture setup script.
func EncryptAES256GCM(keyStr, plaintext string) (string, error) {
	key, err := decodeKey(keyStr, 32)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Use a fixed test nonce for deterministic output in tests.
	// In production use, a random nonce should be generated.
	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	data := append(nonce, ciphertext...)
	return "ENC:" + base64.StdEncoding.EncodeToString(data), nil
}

// decodeKey decodes a key string into exactly `size` bytes.
// Tries hex (2*size chars), then base64, then raw bytes (pad/truncate).
func decodeKey(keyStr string, size int) ([]byte, error) {
	if len(keyStr) == size*2 {
		if key, err := hex.DecodeString(keyStr); err == nil {
			return key, nil
		}
	}
	if key, err := base64.StdEncoding.DecodeString(keyStr); err == nil && len(key) == size {
		return key, nil
	}
	raw := []byte(keyStr)
	if len(raw) >= size {
		return raw[:size], nil
	}
	padded := make([]byte, size)
	copy(padded, raw)
	return padded, nil
}
