package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/acme/driftctl-lite/internal/aws"
	"github.com/acme/driftctl-lite/internal/drift"
	"github.com/acme/driftctl-lite/internal/tfstate"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

func init() {
	var statePath string

	cmd := &cobra.Command{
		Use:   "detect-elasticip",
		Short: "Detect drift in AWS Elastic IP allocations",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := tfstate.ParseStateFile(statePath)
			if err != nil {
				return fmt.Errorf("failed to parse state file: %w", err)
			}

			cfg, err := awsconfig.LoadDefaultConfig(context.Background())
			if err != nil {
				return fmt.Errorf("failed to load AWS config: %w", err)
			}

			client := ec2.NewFromConfig(cfg)
			fetcher := aws.NewElasticIPFetcher(client)

			live, err := fetcher.FetchAll(context.Background())
			if err != nil {
				return fmt.Errorf("failed to fetch Elastic IPs: %w", err)
			}

			results := drift.DetectElasticIPDrift(state.Resources, live)
			if len(results) == 0 {
				log.Println("No Elastic IP drift detected.")
				return nil
			}

			for _, r := range results {
				fmt.Printf("[%s] %s (%s): %s\n", r.DriftType, r.ResourceID, r.ResourceType, r.Details)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	rootCmd.AddCommand(cmd)
}
