# interpolation

Exercises `${VAR}` interpolation, including variables defined after the vars
that reference them (dependency graph resolution).

## Features

- `${VAR}` substitution in `value` fields
- Automatic topological ordering of dependencies
- Variables defined in any order

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```

## Expected output

```
API_ENDPOINT=https://api.example.com/v1
API_VERSION=v1
BASE_URL=https://api.example.com
BIND_ADDR=example.com:8080
HOST=example.com
PORT=8080
SERVICE_NAME=myservice
SERVICE_URL=https://api.example.com/services/myservice
```
