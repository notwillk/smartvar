package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"smartvar/internal/validate"
)

var validateCmd = &cobra.Command{
	Use:   "validate <pattern>",
	Short: "Validate a value from stdin against a named pattern",
	Long: `Read a value from stdin and validate it against a built-in pattern.

Built-in patterns: email, url, hostname, uuid, int, bool, port

Exit codes:
  0  valid
  1  invalid

Examples:
  echo "foo@example.com" | smartvar validate email
  echo "8080" | smartvar validate port`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: validate.Patterns(),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: reading stdin: %v\n", err)
			os.Exit(2)
		}
		value := strings.TrimRight(string(data), "\r\n")
		if err := validate.Validate(pattern, value); err != nil {
			os.Exit(1)
		}
		return nil
	},
}
