package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/spf13/cobra"
	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/drift"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

var detectElasticsearchCmd = &cobra.Command{
	Use:   "elasticsearch",
	Short: "Detect drift in AWS Elasticsearch Service domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath, _ := cmd.Flags().GetString("state")
		if statePath == "" {
			return fmt.Errorf("--state flag is required")
		}

		state, err := tfstate.ParseStateFile(statePath)
		if err != nil {
			return fmt.Errorf("failed to parse state file: %w", err)
		}

		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		client := elasticsearchservice.NewFromConfig(cfg)
		fetcher := aws.NewESFetcher(client)
		live, err := fetcher.FetchAll(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch Elasticsearch domains: %w", err)
		}

		results := drift.DetectElasticsearchDrift(state.Resources, live)
		for _, r := range results {
			fmt.Printf("[%s] %s/%s: %s\n", r.Status, r.ResourceType, r.ResourceID, r.Message)
		}

		if drift.HasDrift(results) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	detectElasticsearchCmd.Flags().String("state", "", "Path to Terraform state file")
	detectCmd.AddCommand(detectElasticsearchCmd)
}
