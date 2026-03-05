# simple-values

Exercises basic literal variable definitions with no interpolation or commands.

## Features

- Static `value` fields
- Default `required: true` behavior

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```

## Expected output

```
APP_NAME=myapp
APP_VERSION=1.0.0
ENVIRONMENT=development
LOG_LEVEL=info
```
