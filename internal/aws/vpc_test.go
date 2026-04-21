package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type mockVPCClient struct {
	output *ec2.DescribeVpcsOutput
	err    error
}

func (m *mockVPCClient) DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error) {
	return m.output, m.err
}

func TestVPCFetchAll_ReturnsVPCs(t *testing.T) {
	mock := &mockVPCClient{
		output: &ec2.DescribeVpcsOutput{
			Vpcs: []types.Vpc{
				{
					VpcId:     aws.String("vpc-abc123"),
					CidrBlock: aws.String("10.0.0.0/16"),
					State:     types.VpcStateAvailable,
					IsDefault: aws.Bool(false),
					Tags: []types.Tag{
						{Key: aws.String("Name"), Value: aws.String("main-vpc")},
					},
				},
			},
		},
	}

	fetcher := NewVPCFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 VPC, got %d", len(results))
	}
	if results[0].ID != "vpc-abc123" {
		t.Errorf("expected vpc-abc123, got %s", results[0].ID)
	}
	if results[0].CIDR != "10.0.0.0/16" {
		t.Errorf("expected 10.0.0.0/16, got %s", results[0].CIDR)
	}
	if results[0].Tags["Name"] != "main-vpc" {
		t.Errorf("expected tag Name=main-vpc")
	}
}

func TestVPCFetchAll_Empty(t *testing.T) {
	mock := &mockVPCClient{
		output: &ec2.DescribeVpcsOutput{Vpcs: []types.Vpc{}},
	}
	fetcher := NewVPCFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 VPCs, got %d", len(results))
	}
}
