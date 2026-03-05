# env-overrides

Exercises environment variable precedence: stdin env and env-file values override
YAML definitions.

## Features

- Stdin env vars override YAML values
- `--env-file` loads a `.env` file (lower precedence than stdin)
- Interpolated vars pick up overridden dependency values

## Run

With env-file override (DB_HOST and APP_ENV come from input.env):

```sh
smartvar compile --config smartvar.yaml --env-file input.env --no-stdin
```

With stdin override (highest precedence):

```sh
echo "DB_HOST=custom.example.com" | smartvar compile --config smartvar.yaml --env-file input.env
```

With inline env:

```sh
DB_HOST=staging.example.com smartvar compile --config smartvar.yaml --no-stdin
```

## Precedence (highest to lowest)

1. stdin env vars (`VAR=value` lines piped to stdin)
2. Process environment variables
3. `--env-file` contents
4. YAML `value` / `command` definitions
