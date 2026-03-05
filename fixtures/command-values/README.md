# command-values

Exercises variables whose values are produced by shell command execution.

## Features

- `command` field: stdout (trimmed) becomes the variable value
- Commands receive the compiled environment
- `cache: true` flag

## Run

```sh
smartvar compile --config smartvar.yaml --no-stdin
smartvar test --config smartvar.yaml --no-stdin
```

## Notes

- `HOSTNAME` and `CURRENT_DATE` values vary by machine/time.
- `STATIC_ID` is deterministic.
- Commands must exit 0 or compilation fails.
