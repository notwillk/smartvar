package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
)

// VerifySignature verifies a cryptographic signature over value.
// Supported algorithms: ed25519.
func VerifySignature(algorithm, keyStr, sigStr, value string) error {
	switch algorithm {
	case "ed25519":
		return verifyEd25519(keyStr, sigStr, value)
	default:
		return fmt.Errorf("unsupported signature algorithm: %s", algorithm)
	}
}

func verifyEd25519(keyStr, sigStr, value string) error {
	pubKey, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return fmt.Errorf("invalid public key encoding: %w", err)
	}
	if len(pubKey) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid ed25519 public key size: expected %d, got %d", ed25519.PublicKeySize, len(pubKey))
	}

	sig, err := base64.StdEncoding.DecodeString(sigStr)
	if err != nil {
		return fmt.Errorf("invalid signature encoding: %w", err)
	}

	if !ed25519.Verify(pubKey, []byte(value), sig) {
		return fmt.Errorf("signature does not match")
	}
	return nil
}

// SignEd25519 signs value with privateKey (base64-encoded) and returns a base64 signature.
// Used in tests and the test fixture setup script.
func SignEd25519(privateKeyStr, value string) (string, error) {
	privKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("invalid private key encoding: %w", err)
	}
	if len(privKey) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid ed25519 private key size")
	}
	sig := ed25519.Sign(privKey, []byte(value))
	return base64.StdEncoding.EncodeToString(sig), nil
}
