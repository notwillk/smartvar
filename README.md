# smartvar

A CLI tool for defining, validating, and compiling environment variables from a declarative YAML configuration.

## Overview

smartvar resolves environment variables from:

- Static literal values
- Interpolated values using `${VAR}` syntax
- Shell command output
- Incoming environment variables (stdin, process env, `.env` files)

With optional:

- Built-in validation patterns (email, url, hostname, uuid, int, bool, port)
- AES-256-GCM decryption of `ENC:`-prefixed values
- ed25519 signature verification
- Secret masking in logs

## Installation

**Via install script (Linux/macOS):**

```sh
curl -fsSL https://github.com/notwillk/smartvar/releases/latest/download/install.sh | sh
```

To pin a specific version:

```sh
SMARTVAR_VERSION=v0.1.0 curl -fsSL https://github.com/notwillk/smartvar/releases/latest/download/install.sh | sh
```

**Windows:**

Download the `.zip` for your architecture from the [releases page](https://github.com/notwillk/smartvar/releases/latest), extract `smartvar.exe`, and add it to your `PATH`.

**From source:**

```sh
just build
# binary at bin/smartvar
```

## Releases

Binaries for Linux, macOS, and Windows (amd64/arm64) are published automatically when a version tag is pushed:

```sh
git tag v0.1.0 && git push --tags
```

GoReleaser builds the binaries, creates a GitHub release, and uploads the install script as a release asset.

## Commands

```
smartvar <command> [flags]
```

| Command | Description |
|---------|-------------|
| `compile` | Produce environment variables |
| `test` | Validate env vars defined in config |
| `validate` | Validate a value from stdin against a pattern |
| `exec` | Execute a command with compiled env vars |
| `json-schema` | Output JSON schema for the config file |
| `completion` | Output shell completion scripts |

## Quick start

```sh
# Compile variables and write to .env
smartvar compile | tee .env

# Run application with compiled env
smartvar exec my_app

# Validate all variables
smartvar test

# Pass env vars via stdin
echo "DB_HOST=prod.example.com" | smartvar compile
```

## Configuration

Default config file: `smartvar.yaml`

```yaml
vars:
  API_HOST:
    value: "api.example.com"

  API_URL:
    value: "https://${API_HOST}"

  BUILD_ID:
    command: "git rev-parse HEAD"

  ADMIN_EMAIL:
    value: "${USER}@example.com"
    test: email

  JWT_SECRET:
    command: "openssl rand -hex 32"
    secret: true

  DATABASE_PASSWORD:
    value: "${ENCRYPTED_DB_PASSWORD}"
    decrypt:
      algorithm: aes256-gcm
      key: "${DECRYPTION_KEY}"
    secret: true
```

## Variable definition fields

| Field | Type | Description |
|-------|------|-------------|
| `value` | string | Literal or `${VAR}`-interpolated value |
| `command` | string | Shell command; trimmed stdout becomes value |
| `test` | string | Validation pattern (see below) |
| `required` | bool | Must resolve (default: `true`) |
| `cache` | bool | Cache command output |
| `secret` | bool | Mask value in logs |
| `decrypt` | object | AES-256-GCM decryption config |
| `signature` | object | ed25519 signature verification config |

Each var must have either `value` or `command`.

## Validation patterns

| Pattern | Description |
|---------|-------------|
| `email` | Email address |
| `url` | HTTP/HTTPS URL |
| `hostname` | DNS hostname |
| `uuid` | UUID (any version) |
| `int` | Integer |
| `bool` | Boolean (`true`/`false`/`1`/`0`) |
| `port` | TCP port (1–65535) |

## Precedence

When the same variable name appears in multiple sources:

1. **stdin env** (highest)
2. **Process environment**
3. **YAML definitions** (lowest)

## `compile` flags

| Flag | Description |
|------|-------------|
| `--config` | Path to YAML config (default: `smartvar.yaml`) |
| `--no-stdin` | Disable reading env vars from stdin |
| `--strict` | Fail if any variable cannot be resolved |
| `--json` | Output JSON instead of `VAR=value` lines |
| `--env-file` | Load variables from a `.env` file |

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Validation failure |
| 2 | Configuration error |
| 3 | Runtime execution error (command failed) |

## Shell integration

```sh
# Export into current shell
eval "$(smartvar compile --no-stdin)"

# Write .env file
smartvar compile --no-stdin > .env

# Run a process
smartvar exec -- node server.js

# Shell completions
smartvar completion bash > /etc/bash_completion.d/smartvar
```

## Development

```sh
just build          # compile binary to bin/smartvar
just test           # run tests
just format         # format code
just static         # lint + format check
just doctor         # check environment health
```

## Fixtures

See [fixtures/](fixtures/) for working examples of each feature.

| Fixture | Features |
|---------|----------|
| [simple-values](fixtures/simple-values/) | Literal values |
| [interpolation](fixtures/interpolation/) | `${VAR}` substitution, dependency ordering |
| [command-values](fixtures/command-values/) | Shell command output |
| [validation](fixtures/validation/) | Built-in patterns |
| [env-overrides](fixtures/env-overrides/) | stdin/env-file override precedence |
| [encrypted-values](fixtures/encrypted-values/) | AES-256-GCM decryption |
| [signed-values](fixtures/signed-values/) | ed25519 signature verification |
