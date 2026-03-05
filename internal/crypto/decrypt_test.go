package crypto

import (
	"strings"
	"testing"
)

func TestEncryptDecryptRoundtrip(t *testing.T) {
	key := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	plaintext := "my-secret-value"

	encrypted, err := EncryptAES256GCM(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	if !strings.HasPrefix(encrypted, "ENC:") {
		t.Errorf("expected ENC: prefix, got: %s", encrypted)
	}

	decrypted, err := Decrypt("aes256-gcm", key, encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("roundtrip: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20"
	wrongKey := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	plaintext := "secret"

	encrypted, err := EncryptAES256GCM(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	_, err = Decrypt("aes256-gcm", wrongKey, encrypted)
	if err == nil {
		t.Error("expected decryption to fail with wrong key")
	}
}

func TestDecryptNoPrefix(t *testing.T) {
	_, err := Decrypt("aes256-gcm", "key", "notencrypted")
	if err == nil {
		t.Error("expected error for value without ENC: prefix")
	}
}

func TestDecryptUnsupportedAlgorithm(t *testing.T) {
	_, err := Decrypt("rot13", "key", "ENC:abc")
	if err == nil {
		t.Error("expected error for unsupported algorithm")
	}
}

func TestDecodeKeyBase64(t *testing.T) {
	// base64("AQIDBAUG") = 6 bytes: 01 02 03 04 05 06
	b64key := "AQIDBAUG"
	key, err := decodeKey(b64key, 6)
	if err != nil {
		t.Fatalf("decodeKey: %v", err)
	}
	if len(key) != 6 {
		t.Errorf("expected 6 bytes, got %d", len(key))
	}
	if key[0] != 1 || key[1] != 2 {
		t.Errorf("unexpected key bytes: %v", key)
	}
}
