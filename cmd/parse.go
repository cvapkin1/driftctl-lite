package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftctl-lite/internal/tfstate"
)

var parseCmd = &cobra.Command{
	Use:   "parse <state-file>",
	Short: "Parse a Terraform state file and list resources",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath := args[0]

		resources, err := tfstate.ParseStateFile(statePath)
		if err != nil {
			return fmt.Errorf("failed to parse state file: %w", err)
		}

		outputJSON, _ := cmd.Flags().GetBool("json")
		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(resources)
		}

		fmt.Printf("Found %d resource(s) in state:\n\n", len(resources))
		for _, r := range resources {
			fmt.Printf("  [%s] %s.%s\n", r.Provider, r.Type, r.Name)
		}
		return nil
	},
}

func init() {
	parseCmd.Flags().Bool("json", false, "Output resources as JSON")
	rootCmd.AddCommand(parseCmd)
}
