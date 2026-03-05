package compile

import (
	"fmt"
	"os/exec"
	"strings"

	"smartvar/internal/config"
	"smartvar/internal/crypto"
	"smartvar/internal/validate"
)

// Options controls how compilation runs.
type Options struct {
	Config     *config.Config
	StdinEnv   map[string]string
	ProcessEnv map[string]string
	Strict     bool
}

// Result holds the compiled environment variables.
type Result struct {
	Vars   map[string]string // all resolved variable values
	Hidden map[string]bool   // vars excluded from output (hidden: true in decrypt config)
	Secret map[string]bool   // vars that should be masked in logs
}

// Compile resolves all variables defined in opts.Config and returns the result.
//
// Precedence: stdin env > process env > YAML definitions.
// Variables overridden by stdin/process env bypass YAML processing entirely.
func Compile(opts Options) (*Result, error) {
	cfg := opts.Config

	// Build the combined override env: process env as baseline, stdin wins.
	overrideEnv := make(map[string]string, len(opts.ProcessEnv)+len(opts.StdinEnv))
	for k, v := range opts.ProcessEnv {
		overrideEnv[k] = v
	}
	for k, v := range opts.StdinEnv {
		overrideEnv[k] = v
	}

	// Build dependency graph: only include refs to other YAML vars.
	deps := make(map[string][]string, len(cfg.Vars))
	for name, def := range cfg.Vars {
		var refs []string
		refs = append(refs, ExtractRefs(def.Value)...)
		refs = append(refs, ExtractRefs(def.Command)...)
		if def.Decrypt != nil {
			refs = append(refs, ExtractRefs(def.Decrypt.Key)...)
		}
		if def.Signature != nil {
			refs = append(refs, ExtractRefs(def.Signature.Key)...)
			refs = append(refs, ExtractRefs(def.Signature.Signature)...)
		}
		deps[name] = refs
	}

	order, err := TopologicalSort(deps)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %w", err)
	}

	resolved := make(map[string]string, len(cfg.Vars))
	hidden := make(map[string]bool)
	secret := make(map[string]bool)

	for _, name := range order {
		def := cfg.Vars[name]

		// Stdin/process env overrides bypass YAML processing.
		if val, ok := overrideEnv[name]; ok {
			resolved[name] = val
			if def.Secret {
				secret[name] = true
			}
			continue
		}

		val, resolveErr := resolveVar(name, def, resolved, overrideEnv)
		if resolveErr != nil {
			if def.IsRequired() || opts.Strict {
				return nil, resolveErr
			}
			continue
		}

		resolved[name] = val
		if def.Secret {
			secret[name] = true
		}
		if def.Decrypt != nil && def.Decrypt.Hidden {
			hidden[name] = true
		}
	}

	return &Result{
		Vars:   resolved,
		Hidden: hidden,
		Secret: secret,
	}, nil
}

func resolveVar(
	name string,
	def config.VarDef,
	resolved map[string]string,
	overrideEnv map[string]string,
) (string, error) {
	// Build lookup env: override env + previously resolved YAML vars.
	env := make(map[string]string, len(overrideEnv)+len(resolved))
	for k, v := range overrideEnv {
		env[k] = v
	}
	for k, v := range resolved {
		env[k] = v
	}

	var val string

	switch {
	case def.Command != "":
		cmdStr := Interpolate(def.Command, env)
		out, err := runCommand(cmdStr, env)
		if err != nil {
			return "", fmt.Errorf("ERROR: command failed for %s: %w", name, err)
		}
		val = out

	case def.Value != "":
		val = Interpolate(def.Value, env)

	default:
		if def.IsRequired() {
			return "", fmt.Errorf("ERROR: missing required variable %s", name)
		}
		return "", fmt.Errorf("no value defined for %s", name)
	}

	if def.Decrypt != nil {
		keyStr := Interpolate(def.Decrypt.Key, env)
		decrypted, err := crypto.Decrypt(def.Decrypt.Algorithm, keyStr, val)
		if err != nil {
			return "", fmt.Errorf("ERROR: decryption failed for %s: %w", name, err)
		}
		val = decrypted
	}

	if def.Signature != nil {
		keyStr := Interpolate(def.Signature.Key, env)
		sigStr := Interpolate(def.Signature.Signature, env)
		if err := crypto.VerifySignature(def.Signature.Algorithm, keyStr, sigStr, val); err != nil {
			return "", fmt.Errorf("ERROR: signature verification failed for %s: %w", name, err)
		}
	}

	if def.Test != "" {
		if err := validate.Validate(def.Test, val); err != nil {
			return "", fmt.Errorf("ERROR: validation failed for %s (%s): %w", name, def.Test, err)
		}
	}

	return val, nil
}

func runCommand(cmdStr string, env map[string]string) (string, error) {
	cmd := exec.Command("sh", "-c", cmdStr)
	envSlice := make([]string, 0, len(env))
	for k, v := range env {
		envSlice = append(envSlice, k+"="+v)
	}
	cmd.Env = envSlice

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(out), "\r\n"), nil
}
