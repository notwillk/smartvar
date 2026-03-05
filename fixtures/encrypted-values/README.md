# encrypted-values

Exercises AES-256-GCM decryption of `ENC:`-prefixed values.

## Features

- `decrypt.algorithm: aes256-gcm`
- `decrypt.key`: resolved from environment (supports `${VAR}` interpolation)
- `secret: true`: value is masked in logs
- `decrypt.hidden: true`: value is decrypted for internal use but not emitted

## Encryption format

Encrypted values have the form:

```
ENC:<base64(12-byte-nonce || ciphertext+tag)>
```

The key is a hex-encoded 32-byte value (64 hex characters).

## Setup

Generate a key and encrypt values using the `smartvar-encrypt` helper
(or any tool implementing AES-256-GCM with the nonce prepended):

```sh
# Generate a key
DECRYPTION_KEY=$(openssl rand -hex 32)

# Encrypt a value using the encrypt integration test helper:
go run ./internal/crypto/cmd/encrypt -key "$DECRYPTION_KEY" -value "my-db-password"
# → ENC:...

# Set the required environment variables
export DECRYPTION_KEY
export ENCRYPTED_DB_PASSWORD="ENC:..."
export ENCRYPTED_API_KEY="ENC:..."
export ENCRYPTED_TOKEN="ENC:..."
```

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```

## Notes

- `INTERNAL_TOKEN` is decrypted but excluded from output (`hidden: true`).
- `DATABASE_PASSWORD` and `API_KEY` appear in output but are marked secret.
