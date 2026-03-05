package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"smartvar/internal/compile"
	"smartvar/internal/config"
	"smartvar/internal/env"
	"smartvar/internal/validate"
)

var testFlags struct {
	cfgPath string
	noStdin bool
	envFile string
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Validate env vars defined in config",
	Long: `Load the config, resolve all variables, and check that all validation
rules and signature verifications pass.

Exit codes:
  0  all valid
  1  validation failure
  2  configuration error

Examples:
  smartvar test
  smartvar test < .env
  smartvar test --config staging.yaml`,
	Args: cobra.NoArgs,
	RunE: runTest,
}

func init() {
	testCmd.Flags().StringVar(&testFlags.cfgPath, "config", "smartvar.yaml", "path to YAML config file")
	testCmd.Flags().BoolVar(&testFlags.noStdin, "no-stdin", false, "disable reading env vars from stdin")
	testCmd.Flags().StringVar(&testFlags.envFile, "env-file", "", "load additional variables from a .env file")
}

func runTest(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(testFlags.cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}

	var stdinEnv map[string]string
	if testFlags.noStdin {
		stdinEnv = make(map[string]string)
	} else {
		stdinEnv, err = env.ParseEnvLines(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: reading stdin: %v\n", err)
			os.Exit(2)
		}
	}

	if testFlags.envFile != "" {
		fileEnv, ferr := loadEnvFile(testFlags.envFile)
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", ferr)
			os.Exit(2)
		}
		for k, v := range stdinEnv {
			fileEnv[k] = v
		}
		stdinEnv = fileEnv
	}

	processEnv := env.ProcessEnv()
	overrideEnv := env.Merge(processEnv, stdinEnv)

	result, compileErr := compile.Compile(compile.Options{
		Config:     cfg,
		StdinEnv:   stdinEnv,
		ProcessEnv: processEnv,
		Strict:     true,
	})

	failed := false

	if compileErr != nil {
		fmt.Fprintln(os.Stderr, compileErr)
		failed = true
	}

	// Additionally validate any overridden vars against their YAML test patterns.
	if result != nil {
		for name, def := range cfg.Vars {
			if def.Test == "" {
				continue
			}
			if _, overridden := overrideEnv[name]; !overridden {
				continue // already handled by compile
			}
			if val, ok := result.Vars[name]; ok {
				if verr := validate.Validate(def.Test, val); verr != nil {
					fmt.Fprintf(os.Stderr, "ERROR: validation failed for %s (%s): %v\n", name, def.Test, verr)
					failed = true
				}
			}
		}
	}

	if failed {
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "OK: all variables valid")
	return nil
}
