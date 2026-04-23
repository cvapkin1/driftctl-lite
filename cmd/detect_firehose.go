package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/eliraz-refael/driftctl-lite/internal/aws"
	"github.com/eliraz-refael/driftctl-lite/internal/drift"
	"github.com/eliraz-refael/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

func init() {
	var statePath string

	cmd := &cobra.Command{
		Use:   "detect-firehose",
		Short: "Detect drift in Kinesis Firehose delivery streams",
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

			fetcher := aws.NewFirehoseFetcher(firehose.NewFromConfig(cfg))
			streams, err := fetcher.FetchAll(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch Firehose streams: %w", err)
			}

			results := drift.DetectFirehoseDrift(state.Resources, streams)
			if len(results) == 0 {
				fmt.Println("No Firehose drift detected.")
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
