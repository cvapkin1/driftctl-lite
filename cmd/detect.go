package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/example/driftctl-lite/internal/aws"
	"github.com/example/driftctl-lite/internal/drift"
	"github.com/example/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect drift between Terraform state and live AWS resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath, _ := cmd.Flags().GetString("state")
		state, err := tfstate.ParseStateFile(statePath)
		if err != nil {
			return fmt.Errorf("parsing state: %w", err)
		}

		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return fmt.Errorf("loading AWS config: %w", err)
		}

		fetcher := aws.NewEC2Fetcher(cfg)
		live, err := fetcher.FetchAll(context.Background())
		if err != nil {
			return fmt.Errorf("fetching EC2 instances: %w", err)
		}

		results := drift.DetectEC2Drift(state.Resources, live)
		if len(results) == 0 {
			fmt.Println("No drift detected.")
			return nil
		}
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "DRIFT [%s] %s: %s\n", r.ResourceType, r.ResourceID, r.Reason)
		}
		return nil
	},
}

func init() {
	detectCmd.Flags().String("state", "terraform.tfstate", "path to Terraform state file")
	rootCmd.AddCommand(detectCmd)
}
