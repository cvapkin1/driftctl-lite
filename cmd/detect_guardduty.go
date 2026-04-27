package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/drift"
	"github.com/snyk/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

func init() {
	var statePath string

	cmd := &cobra.Command{
		Use:   "detect-guardduty",
		Short: "Detect drift in AWS GuardDuty detectors",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := tfstate.ParseStateFile(statePath)
			if err != nil {
				return fmt.Errorf("failed to parse state file: %w", err)
			}

			ctx := context.Background()
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to load AWS config: %w", err)
			}

			client := guardduty.NewFromConfig(cfg)
			fetcher := aws.NewGuardDutyFetcher(client)
			detectors, err := fetcher.FetchAll(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch GuardDuty detectors: %w", err)
			}

			results := drift.DetectGuardDutyDrift(state.Resources, detectors)
			if len(results) == 0 {
				fmt.Println("No GuardDuty drift detected.")
				return nil
			}

			for _, r := range results {
				fmt.Fprintf(os.Stdout, "[%s] %s (%s): %s\n", r.DriftType, r.ResourceID, r.ResourceType, r.Details)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	rootCmd.AddCommand(cmd)
}
