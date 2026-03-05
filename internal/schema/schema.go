package schema

import "encoding/json"

// Schema returns the JSON schema describing the smartvar YAML configuration.
func Schema() ([]byte, error) {
	s := map[string]any{
		"$schema":     "http://json-schema.org/draft-07/schema#",
		"title":       "smartvar configuration",
		"description": "Declarative environment variable definition and validation.",
		"type":        "object",
		"properties": map[string]any{
			"vars": map[string]any{
				"type":        "object",
				"description": "Map of variable name to variable definition.",
				"additionalProperties": map[string]any{
					"$ref": "#/definitions/VarDef",
				},
			},
		},
		"definitions": map[string]any{
			"VarDef": map[string]any{
				"type": "object",
				"oneOf": []any{
					map[string]any{"required": []string{"value"}},
					map[string]any{"required": []string{"command"}},
				},
				"properties": map[string]any{
					"value": map[string]any{
						"type":        "string",
						"description": "Literal or interpolated value using ${VAR} syntax.",
					},
					"command": map[string]any{
						"type":        "string",
						"description": "Shell command whose stdout (trimmed) becomes the value.",
					},
					"test": map[string]any{
						"type":        "string",
						"description": "Built-in validation pattern to apply to the resolved value.",
						"enum":        []string{"email", "url", "hostname", "uuid", "int", "bool", "port"},
					},
					"required": map[string]any{
						"type":        "boolean",
						"description": "Whether the variable must resolve. Defaults to true.",
					},
					"cache": map[string]any{
						"type":        "boolean",
						"description": "Cache the command output (run the command only once).",
					},
					"secret": map[string]any{
						"type":        "boolean",
						"description": "Mark value as sensitive; masked in logs.",
					},
					"decrypt": map[string]any{
						"$ref": "#/definitions/DecryptDef",
					},
					"signature": map[string]any{
						"$ref": "#/definitions/SigDef",
					},
				},
				"additionalProperties": false,
			},
			"DecryptDef": map[string]any{
				"type":        "object",
				"description": "Decryption configuration for ENC:-prefixed values.",
				"required":    []string{"algorithm", "key"},
				"properties": map[string]any{
					"algorithm": map[string]any{
						"type":        "string",
						"description": "Decryption algorithm.",
						"enum":        []string{"aes256-gcm"},
					},
					"key": map[string]any{
						"type":        "string",
						"description": "Decryption key (hex or base64). Supports ${VAR} interpolation.",
					},
					"hidden": map[string]any{
						"type":        "boolean",
						"description": "Exclude the decrypted value from compile output.",
					},
				},
				"additionalProperties": false,
			},
			"SigDef": map[string]any{
				"type":        "object",
				"description": "Signature verification configuration.",
				"required":    []string{"algorithm", "key", "signature"},
				"properties": map[string]any{
					"algorithm": map[string]any{
						"type":        "string",
						"description": "Signature algorithm.",
						"enum":        []string{"ed25519"},
					},
					"key": map[string]any{
						"type":        "string",
						"description": "Base64-encoded public key. Supports ${VAR} interpolation.",
					},
					"signature": map[string]any{
						"type":        "string",
						"description": "Base64-encoded signature. Supports ${VAR} interpolation.",
					},
				},
				"additionalProperties": false,
			},
		},
	}
	return json.MarshalIndent(s, "", "  ")
}
