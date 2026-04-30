package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/tfstate"
	"github.com/spf13/cobra"
)

var detectInternetGatewayCmd = &cobra.Command{
	Use:   "internet-gateway",
	Short: "Detect drift in AWS Internet Gateways",
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath, _ := cmd.Flags().GetString("state")
		state, err := tfstate.ParseStateFile(statePath)
		if err != nil {
			return fmt.Errorf("failed to parse state file: %w", err)
		}

		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		fetcher := aws.NewInternetGatewayFetcher(ec2.NewFromConfig(cfg))
		live, err := fetcher.FetchAll(context.Background())
		if err != nil {
			return fmt.Errorf("failed to fetch internet gateways: %w", err)
		}

		results := drift.DetectInternetGatewayDrift(state.Resources, live)
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "[%s] %s\n", r.Status, r.Message)
		}
		return nil
	},
}

func init() {
	detectInternetGatewayCmd.Flags().String("state", "terraform.tfstate", "Path to Terraform state file")
	detectCmd.AddCommand(detectInternetGatewayCmd)
}
