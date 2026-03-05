package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"smartvar/internal/schema"
)

var jsonSchemaCmd = &cobra.Command{
	Use:   "json-schema",
	Short: "Output JSON schema for the smartvar config file",
	Long: `Output the JSON schema describing the smartvar YAML configuration format.
Useful for IDE validation, CI validation, and documentation.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := schema.Schema()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: generating schema: %v\n", err)
			os.Exit(2)
		}
		fmt.Println(string(data))
		return nil
	},
}
