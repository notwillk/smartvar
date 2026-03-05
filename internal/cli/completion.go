package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion <shell>",
	Short: "Output shell completion scripts",
	Long: `Output shell completion scripts for bash, zsh, or fish.

Examples:
  smartvar completion bash > /etc/bash_completion.d/smartvar
  smartvar completion zsh > ~/.zfunc/_smartvar
  smartvar completion fish > ~/.config/fish/completions/smartvar.fish`,
	Args:               cobra.ExactArgs(1),
	ValidArgs:          []string{"bash", "zsh", "fish"},
	DisableFlagParsing: false,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		default:
			fmt.Fprintf(os.Stderr, "ERROR: unsupported shell %q (supported: bash, zsh, fish)\n", args[0])
			os.Exit(2)
		}
		return nil
	},
}
