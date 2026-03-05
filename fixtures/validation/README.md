# validation

Exercises built-in validation patterns applied to resolved variable values.

## Features

- `test: email` — validates email format
- `test: url` — validates HTTP/HTTPS URL
- `test: hostname` — validates DNS hostname
- `test: uuid` — validates UUID format
- `test: int` — validates integer value
- `test: bool` — validates boolean (true/false/1/0)
- `test: port` — validates TCP port (1–65535)

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```

Validation failures produce:

```
ERROR: validation failed for VARNAME (pattern): reason
```

and exit with code 1.
