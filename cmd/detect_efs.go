package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

var detectEFSCmd = &cobra.Command{
	Use:   "efs",
	Short: "Detect drift in AWS EFS file systems",
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath, _ := cmd.Flags().GetString("state")
		region, _ := cmd.Flags().GetString("region")

		state, err := tfstate.ParseStateFile(statePath)
		if err != nil {
			return fmt.Errorf("failed to parse state file: %w", err)
		}

		ctx := context.Background()
		fetcher, err := aws.NewEFSFetcher(ctx, region)
		if err != nil {
			return fmt.Errorf("failed to create EFS fetcher: %w", err)
		}

		fileSystems, err := fetcher.FetchAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch EFS file systems: %w", err)
		}

		results := drift.DetectEFSDrift(state.Resources, fileSystems)
		hasDrift := false
		for _, r := range results {
			if r.Status != "ok" {
				hasDrift = true
				fmt.Fprintf(os.Stdout, "[%s] %s\n", r.Status, r.Message)
			}
		}
		if !hasDrift {
			fmt.Println("No EFS drift detected.")
		}
		return nil
	},
}

func init() {
	detectEFSCmd.Flags().String("state", "terraform.tfstate", "Path to Terraform state file")
	detectEFSCmd.Flags().String("region", "us-east-1", "AWS region")
	detectCmd.AddCommand(detectEFSCmd)
}
