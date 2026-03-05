package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"smartvar/internal/compile"
	"smartvar/internal/config"
	"smartvar/internal/env"
)

var execFlags struct {
	cfgPath string
	noStdin bool
	envFile string
}

var execCmd = &cobra.Command{
	Use:   "exec <command> [args...]",
	Short: "Execute a command with compiled environment variables",
	Long: `Compile environment variables and run the given command with those
variables merged into the current process environment.

On success the process is replaced by the command; on compilation failure
smartvar exits with an error before launching the command.

Examples:
  smartvar exec my_app
  smartvar exec -- node server.js --port 3000`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: false,
	RunE:               runExec,
}

func init() {
	execCmd.Flags().StringVar(&execFlags.cfgPath, "config", "smartvar.yaml", "path to YAML config file")
	execCmd.Flags().BoolVar(&execFlags.noStdin, "no-stdin", false, "disable reading env vars from stdin")
	execCmd.Flags().StringVar(&execFlags.envFile, "env-file", "", "load additional variables from a .env file")
}

func runExec(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(execFlags.cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}

	var stdinEnv map[string]string
	if execFlags.noStdin {
		stdinEnv = make(map[string]string)
	} else {
		stdinEnv, err = env.ParseEnvLines(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: reading stdin: %v\n", err)
			os.Exit(2)
		}
	}

	if execFlags.envFile != "" {
		fileEnv, ferr := loadEnvFile(execFlags.envFile)
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", ferr)
			os.Exit(2)
		}
		for k, v := range stdinEnv {
			fileEnv[k] = v
		}
		stdinEnv = fileEnv
	}

	result, compileErr := compile.Compile(compile.Options{
		Config:     cfg,
		StdinEnv:   stdinEnv,
		ProcessEnv: env.ProcessEnv(),
	})
	if compileErr != nil {
		fmt.Fprintln(os.Stderr, compileErr)
		os.Exit(1)
	}

	// Build env: current process env + compiled vars (compiled wins for defined vars).
	envMap := env.ProcessEnv()
	for k, v := range result.Vars {
		if !result.Hidden[k] {
			envMap[k] = v
		}
	}
	envSlice := make([]string, 0, len(envMap))
	for k, v := range envMap {
		envSlice = append(envSlice, k+"="+v)
	}

	c := exec.Command(args[0], args[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = envSlice

	if err := c.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "ERROR: exec: %v\n", err)
		os.Exit(3)
	}
	return nil
}
