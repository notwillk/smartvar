package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "smartvar",
	Short: "Define, validate, and compile environment variables from declarative YAML",
	Long: `smartvar resolves environment variables from static values, interpolated
strings, command output, and incoming environment — with built-in validation,
decryption, and signature verification.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(jsonSchemaCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(execCmd)
}
