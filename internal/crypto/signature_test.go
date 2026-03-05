package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
)

func TestVerifyEd25519RoundTrip(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	pubB64 := base64.StdEncoding.EncodeToString(pub)
	privB64 := base64.StdEncoding.EncodeToString(priv)

	value := "hello-world"

	sig, err := SignEd25519(privB64, value)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	if err := VerifySignature("ed25519", pubB64, sig, value); err != nil {
		t.Errorf("verify: %v", err)
	}
}

func TestVerifyEd25519WrongKey(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	pub2, _, _ := ed25519.GenerateKey(rand.Reader)

	privB64 := base64.StdEncoding.EncodeToString(priv)
	pub2B64 := base64.StdEncoding.EncodeToString(pub2)

	value := "payload"
	sig, _ := SignEd25519(privB64, value)

	err := VerifySignature("ed25519", pub2B64, sig, value)
	if err == nil {
		t.Error("expected verification to fail with wrong public key")
	}
}

func TestVerifyEd25519WrongValue(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	pubB64 := base64.StdEncoding.EncodeToString(pub)
	privB64 := base64.StdEncoding.EncodeToString(priv)

	sig, _ := SignEd25519(privB64, "original")

	err := VerifySignature("ed25519", pubB64, sig, "tampered")
	if err == nil {
		t.Error("expected verification to fail for tampered value")
	}
}

func TestVerifyUnsupportedAlgorithm(t *testing.T) {
	err := VerifySignature("rsa", "key", "sig", "value")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}
