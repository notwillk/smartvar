package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"smartvar/internal/compile"
	"smartvar/internal/config"
	"smartvar/internal/env"
)

var compileFlags struct {
	cfgPath string
	noStdin bool
	strict  bool
	jsonOut bool
	envFile string
}

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Produce environment variables from config",
	Long: `Resolve and emit all environment variables defined in the smartvar config.

Inputs are merged with the following precedence (highest to lowest):
  1. stdin env vars (VAR=value lines)
  2. process environment variables
  3. YAML definitions

Output format (default): VAR=value
Output format (--json): {"VAR": "value"}

Exit codes:
  0  success
  1  validation failure
  2  configuration error
  3  runtime execution error`,
	Args: cobra.NoArgs,
	RunE: runCompile,
}

func init() {
	compileCmd.Flags().StringVar(&compileFlags.cfgPath, "config", "smartvar.yaml", "path to YAML config file")
	compileCmd.Flags().BoolVar(&compileFlags.noStdin, "no-stdin", false, "disable reading env vars from stdin")
	compileCmd.Flags().BoolVar(&compileFlags.strict, "strict", false, "fail if any variable cannot be resolved")
	compileCmd.Flags().BoolVar(&compileFlags.jsonOut, "json", false, "output JSON instead of VAR=value lines")
	compileCmd.Flags().StringVar(&compileFlags.envFile, "env-file", "", "load additional variables from a .env file")
}

func runCompile(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(compileFlags.cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}

	stdinEnv, err := loadStdinEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: reading stdin: %v\n", err)
		os.Exit(2)
	}

	if compileFlags.envFile != "" {
		fileEnv, err := loadEnvFile(compileFlags.envFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(2)
		}
		// env-file is lower precedence than stdin
		for k, v := range stdinEnv {
			fileEnv[k] = v
		}
		stdinEnv = fileEnv
	}

	result, err := compile.Compile(compile.Options{
		Config:     cfg,
		StdinEnv:   stdinEnv,
		ProcessEnv: env.ProcessEnv(),
		Strict:     compileFlags.strict,
	})
	if err != nil {
		msg := err.Error()
		fmt.Fprintln(os.Stderr, msg)
		if isRuntimeError(msg) {
			os.Exit(3)
		}
		os.Exit(1)
	}

	if compileFlags.jsonOut {
		printJSON(result)
	} else {
		printEnv(result)
	}
	return nil
}

func loadStdinEnv() (map[string]string, error) {
	if compileFlags.noStdin {
		return make(map[string]string), nil
	}
	return env.ParseEnvLines(os.Stdin)
}

func loadEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()
	return env.ParseEnvLines(f)
}

func printEnv(result *compile.Result) {
	names := make([]string, 0, len(result.Vars))
	for k := range result.Vars {
		if !result.Hidden[k] {
			names = append(names, k)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s=%s\n", name, shellQuote(result.Vars[name]))
	}
}

func printJSON(result *compile.Result) {
	out := make(map[string]string, len(result.Vars))
	for k, v := range result.Vars {
		if !result.Hidden[k] {
			out[k] = v
		}
	}
	data, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(data))
}

// shellQuote returns the value quoted for safe use in a VAR=value line.
// Simple values are returned as-is; others are single-quoted.
func shellQuote(s string) string {
	if isSafeUnquoted(s) {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func isSafeUnquoted(s string) bool {
	for _, c := range s {
		switch c {
		case ' ', '\t', '\n', '\r', '"', '\'', '\\', '$', '`', '!', '&', ';', '|', '<', '>', '(', ')', '{', '}', '#', '~', '*', '?', '[', ']':
			return false
		}
	}
	return true
}

func isRuntimeError(msg string) bool {
	return strings.Contains(msg, "command failed")
}
