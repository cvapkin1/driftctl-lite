package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSubnetClient struct {
	output *ec2.DescribeSubnetsOutput
	err    error
}

func (m *mockSubnetClient) DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error) {
	return m.output, m.err
}

func TestSubnetFetchAll_ReturnsSubnets(t *testing.T) {
	mock := &mockSubnetClient{
		output: &ec2.DescribeSubnetsOutput{
			Subnets: []types.Subnet{
				{
					SubnetId:            aws.String("subnet-abc123"),
					VpcId:               aws.String("vpc-111"),
					CidrBlock:           aws.String("10.0.1.0/24"),
					AvailabilityZone:    aws.String("us-east-1a"),
					State:               types.SubnetStateAvailable,
					MapPublicIpOnLaunch: aws.Bool(true),
				},
			},
		},
	}

	fetcher := NewSubnetFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "subnet-abc123", results[0].ID)
	assert.Equal(t, "vpc-111", results[0].VPCID)
	assert.Equal(t, "10.0.1.0/24", results[0].CIDRBlock)
	assert.Equal(t, "us-east-1a", results[0].AvailabilityZone)
	assert.Equal(t, "available", results[0].State)
	assert.True(t, results[0].MapPublicIPOnLaunch)
}

func TestSubnetFetchAll_Empty(t *testing.T) {
	mock := &mockSubnetClient{
		output: &ec2.DescribeSubnetsOutput{Subnets: []types.Subnet{}},
	}

	fetcher := NewSubnetFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	require.NoError(t, err)
	assert.Empty(t, results)
}
