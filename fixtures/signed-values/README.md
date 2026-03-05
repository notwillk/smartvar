# signed-values

Exercises ed25519 signature verification of resolved variable values.

## Features

- `signature.algorithm: ed25519`
- `signature.key`: base64-encoded public key (supports `${VAR}` interpolation)
- `signature.signature`: base64-encoded signature (supports `${VAR}` interpolation)

## Verification

Signature verification runs after value resolution:

1. The variable value is resolved (from `value`, `command`, or env override)
2. The signature is verified: `ed25519.Verify(pubKey, []byte(value), sig)`
3. Failure → `ERROR: signature verification failed for VARNAME`

## Setup

Generate keys and sign a value using Go's `crypto/ed25519`:

```go
import (
    "crypto/ed25519"
    "crypto/rand"
    "encoding/base64"
)

pub, priv, _ := ed25519.GenerateKey(rand.Reader)
pubB64 := base64.StdEncoding.EncodeToString(pub)
privB64 := base64.StdEncoding.EncodeToString(priv)

value := "my-verified-payload"
sig := ed25519.Sign(priv, []byte(value))
sigB64 := base64.StdEncoding.EncodeToString(sig)
```

Or use the signature test helper in `internal/crypto/signature_test.go`.

```sh
export SIGNING_PUBLIC_KEY="<base64 public key>"
export SIGNED_PAYLOAD="my-verified-payload"
export PAYLOAD_SIGNATURE="<base64 signature over SIGNED_PAYLOAD>"
```

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```
