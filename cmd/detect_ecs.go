package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/drift"
	"github.com/snyk/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

func init() {
	var statePath string

	var detectECSCmd = &cobra.Command{
		Use:   "detect-ecs",
		Short: "Detect drift in ECS clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			state, err := tfstate.ParseStateFile(statePath)
			if err != nil {
				return fmt.Errorf("failed to parse state file: %w", err)
			}

			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to load AWS config: %w", err)
			}

			fetcher := aws.NewECSFetcher(ecs.NewFromConfig(cfg))
			clusters, err := fetcher.FetchAll(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch ECS clusters: %w", err)
			}

			results := drift.DetectECSDrift(state.Resources, clusters)
			if len(results) == 0 {
				fmt.Println("No ECS drift detected.")
				return nil
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(results)
		},
	}

	detectECSCmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	rootCmd.AddCommand(detectECSCmd)
}
